package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/minio/minio/cmd/logger"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"syscall"
	"time"
)

var sslRequiredErrMsg = []byte("HTTP/1.0 403 Forbidden\r\n\r\nSSL required")

// HTTP methods.
var methods = []string{
	http.MethodGet,
	http.MethodHead,
	http.MethodPost,
	http.MethodPut,
	http.MethodPatch,
	http.MethodDelete,
	http.MethodConnect,
	http.MethodOptions,
	http.MethodTrace,
	"PRI", // HTTP 2 method
}

// maximum length of above methods + one space.
var methodMaxLen = getMethodMaxLen() + 1

func getMethodMaxLen() int {
	maxLen := 0
	for _, method := range methods {
		if len(method) > maxLen {
			maxLen = len(method)
		}
	}

	return maxLen
}

func isHTTPMethod(s string) bool {
	for _, method := range methods {
		if s == method {
			return true
		}
	}

	return false
}

type acceptResult struct {
	conn net.Conn
	err  error
}

// httpListener
type httpListener struct {
	mutex                  sync.Mutex        // to guard Close() method.
	tcpListener            *net.TCPListener  // underlaying TCP listener.
	acceptCh               chan acceptResult // channel where all TCP listeners write accepted connection.
	doneCh                 chan struct{}     // done channel for TCP listener goroutines.
	tlsConfig              *tls.Config       // TLS configuration
	tcpKeepAliveTimeout    time.Duration
	readTimeout            time.Duration
	writeTimeout           time.Duration
	updateBytesReadFunc    func(int) // function to be called to update bytes read in BufConn.
	updateBytesWrittenFunc func(int) // function to be called to update bytes written in BufConn.
}

// isRoutineNetErr returns true if error is due to a network timeout,
// connect reset or io.EOF and false otherwise
func isRoutineNetErr(err error) bool {
	if nErr, ok := err.(*net.OpError); ok {
		// Check if the error is a tcp connection reset
		if syscallErr, ok := nErr.Err.(*os.SyscallError); ok {
			if errno, ok := syscallErr.Err.(syscall.Errno); ok {
				return errno == syscall.ECONNRESET
			}
		}
		// Check if the error is a timeout
		return nErr.Timeout()
	}
	return err == io.EOF
}

// start - starts separate goroutine for each TCP listener.  A valid insecure/TLS HTTP new connection is passed to httpListener.acceptCh.
func (listener *httpListener) start() {
	listener.acceptCh = make(chan acceptResult)
	listener.doneCh = make(chan struct{})

	// Closure to send acceptResult to acceptCh.
	// It returns true if the result is sent else false if returns when doneCh is closed.
	send := func(result acceptResult, doneCh <-chan struct{}) bool {
		select {
		case listener.acceptCh <- result:
			// Successfully written to acceptCh
			return true
		case <-doneCh:
			// As stop signal is received, close accepted connection.
			if result.conn != nil {
				result.conn.Close()
			}
			return false
		}
	}

	// Closure to handle single connection.
	handleConn := func(tcpConn *net.TCPConn, doneCh <-chan struct{}) {
		// Tune accepted TCP connection.
		tcpConn.SetKeepAlive(true)
		tcpConn.SetKeepAlivePeriod(listener.tcpKeepAliveTimeout)

		bufconn := newBufConn(tcpConn, listener.readTimeout, listener.writeTimeout,
			listener.updateBytesReadFunc, listener.updateBytesWrittenFunc)

		// Peek bytes of maximum length of all HTTP methods.
		data, err := bufconn.Peek(methodMaxLen)
		if err != nil {
			// Peek could fail legitimately when clients abruptly close
			// connection. E.g. Chrome browser opens connections speculatively to
			// speed up loading of a web page. Peek may also fail due to network
			// saturation on a transport with read timeout set. All other kind of
			// errors should be logged for further investigation. Thanks @brendanashworth.
			if !isRoutineNetErr(err) {
				reqInfo := (&logger.ReqInfo{}).AppendTags("remoteAddr", bufconn.RemoteAddr().String())
				reqInfo.AppendTags("localAddr", bufconn.LocalAddr().String())
				ctx := logger.SetReqInfo(context.Background(), reqInfo)
				logger.LogIf(ctx, err)
			}
			bufconn.Close()
			return
		}

		// Return bufconn if read data is a valid HTTP method.
		tokens := strings.SplitN(string(data), " ", 2)
		if isHTTPMethod(tokens[0]) {
			if listener.tlsConfig == nil {
				send(acceptResult{bufconn, nil}, doneCh)
			} else {
				// As TLS is configured and we got plain text HTTP request,
				// return 403 (forbidden) error.
				bufconn.Write(sslRequiredErrMsg)
				bufconn.Close()
			}
			return
		}

		if listener.tlsConfig != nil {
			// As the listener is configured with TLS, try to do TLS handshake, drop the connection if it fails.
			tlsConn := tls.Server(bufconn, listener.tlsConfig)
			if err = tlsConn.Handshake(); err != nil {
				reqInfo := (&logger.ReqInfo{}).AppendTags("remoteAddr", bufconn.RemoteAddr().String())
				reqInfo.AppendTags("localAddr", bufconn.LocalAddr().String())
				ctx := logger.SetReqInfo(context.Background(), reqInfo)
				logger.LogIf(ctx, err)
				bufconn.Close()
				return
			}

			// Check whether the connection contains HTTP request or not.
			bufconn = newBufConn(tlsConn, listener.readTimeout, listener.writeTimeout,
				listener.updateBytesReadFunc, listener.updateBytesWrittenFunc)

			// Peek bytes of maximum length of all HTTP methods.
			data, err = bufconn.Peek(methodMaxLen)
			if err != nil {
				if !isRoutineNetErr(err) {
					reqInfo := (&logger.ReqInfo{}).AppendTags("remoteAddr", bufconn.RemoteAddr().String())
					reqInfo.AppendTags("localAddr", bufconn.LocalAddr().String())
					ctx := logger.SetReqInfo(context.Background(), reqInfo)
					logger.LogIf(ctx, err)
				}
				bufconn.Close()
				return
			}

			// Return bufconn if read data is a valid HTTP method.
			tokens := strings.SplitN(string(data), " ", 2)
			if isHTTPMethod(tokens[0]) {
				send(acceptResult{bufconn, nil}, doneCh)
				return
			}
		}
		reqInfo := (&logger.ReqInfo{}).AppendTags("remoteAddr", bufconn.RemoteAddr().String())
		reqInfo.AppendTags("localAddr", bufconn.LocalAddr().String())
		ctx := logger.SetReqInfo(context.Background(), reqInfo)
		logger.LogIf(ctx, err)

		bufconn.Close()
		return
	}

	// Closure to handle TCPListener until done channel is closed.
	handleListener := func(tcpListener *net.TCPListener, doneCh <-chan struct{}) {
		for {
			tcpConn, err := tcpListener.AcceptTCP()
			if err != nil {
				// Returns when send fails.
				if !send(acceptResult{nil, err}, doneCh) {
					return
				}
			} else {
				go handleConn(tcpConn, doneCh)
			}
		}
	}

	// Start separate goroutine for each TCP listener to handle connection.

	go handleListener(listener.tcpListener, listener.doneCh)

}

// Accept - reads from httpListener.acceptCh for one of previously accepted TCP connection and returns the same.
func (listener *httpListener) Accept() (conn net.Conn, err error) {
	result, ok := <-listener.acceptCh
	if ok {
		return result.conn, result.err
	}

	return nil, syscall.EINVAL
}

// Close - closes underneath all TCP listeners.
func (listener *httpListener) Close() (err error) {
	listener.mutex.Lock()
	defer listener.mutex.Unlock()
	if listener.doneCh == nil {
		return syscall.EINVAL
	}

	listener.tcpListener.Close()
	close(listener.doneCh)

	listener.doneCh = nil
	return nil
}

// Addr - net.Listener interface compatible method returns net.Addr.  In case of multiple TCP listeners, it returns '0.0.0.0' as IP address.
func (listener *httpListener) Addr() (addr net.Addr) {
	addr = listener.tcpListener.Addr()

	tcpAddr := addr.(*net.TCPAddr)
	if ip := net.ParseIP("0.0.0.0"); ip != nil {
		tcpAddr.IP = ip
	}

	addr = tcpAddr
	return addr
}

// newHTTPListener - creates new httpListener object which is interface compatible to net.Listener.
// httpListener is capable to
// * listen to multiple addresses
// * controls incoming connections only doing HTTP protocol
func NewHTTPListener(serverAddr string,
	tlsConfig *tls.Config,
	tcpKeepAliveTimeout time.Duration,
	readTimeout time.Duration,
	writeTimeout time.Duration,
	updateBytesReadFunc func(int),
	updateBytesWrittenFunc func(int)) (listener *httpListener, err error) {

	var tcpListener *net.TCPListener
	// Close all opened listeners on error
	defer func() {
		if err == nil {
			return
		}

		tcpListener.Close()

	}()

	var l net.Listener
	if l, err = net.Listen("tcp", serverAddr); err != nil {
		return nil, err
	}

	tcpListener, ok := l.(*net.TCPListener)
	if !ok {
		return nil, fmt.Errorf("unexpected listener type found %v, expected net.TCPListener", l)
	}

	listener = &httpListener{
		tcpListener:            tcpListener,
		tlsConfig:              nil,
		tcpKeepAliveTimeout:    tcpKeepAliveTimeout,
		readTimeout:            readTimeout,
		writeTimeout:           writeTimeout,
		updateBytesReadFunc:    updateBytesReadFunc,
		updateBytesWrittenFunc: updateBytesWrittenFunc,
	}
	listener.start()

	return listener, nil
}
