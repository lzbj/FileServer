package server

import (
	"crypto/tls"
	"errors"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/dustin/go-humanize"
)

const (
	serverShutdownPoll = 500 * time.Millisecond

	// DefaultShutdownTimeout - default shutdown timeout used for graceful http server shutdown.
	DefaultShutdownTimeout = 5 * time.Second

	// DefaultTCPKeepAliveTimeout - default TCP keep alive timeout for accepted connection.
	DefaultTCPKeepAliveTimeout = 10 * time.Second

	// DefaultReadTimeout - default timout to read data from accepted connection.
	DefaultReadTimeout = 5 * time.Minute

	// DefaultWriteTimeout - default timout to write data to accepted connection.
	DefaultWriteTimeout = 5 * time.Minute

	// DefaultMaxHeaderBytes - default maximum HTTP header size in bytes.
	DefaultMaxHeaderBytes = 1 * humanize.MiByte
)

// Server - extended http.Server supports multiple addresses to serve and enhanced connection handling.
type Server struct {
	http.Server
	Addrs                  string        // addresses on which the server listens for new connection.
	ShutdownTimeout        time.Duration // timeout used for graceful server shutdown.
	TCPKeepAliveTimeout    time.Duration // timeout used for underneath TCP connection.
	UpdateBytesReadFunc    func(int)     // function to be called to update bytes read in bufConn.
	UpdateBytesWrittenFunc func(int)     // function to be called to update bytes written in bufConn.
	listenerMutex          *sync.Mutex   // to guard 'listener' field.
	listener               *httpListener // HTTP listener for all 'Addrs' field.
	inShutdown             uint32        // indicates whether the server is in shutdown or not
	requestCount           int32         // counter holds no. of request in process.
}

// Start - start HTTP server
func (srv *Server) Start() (err error) {
	// Take a copy of server fields.
	var tlsConfig *tls.Config
	if srv.TLSConfig != nil {
		tlsConfig = srv.TLSConfig.Clone()
	}
	readTimeout := srv.ReadTimeout
	writeTimeout := srv.WriteTimeout
	handler := srv.Handler // if srv.Handler holds non-synced state -> possible data race

	addrs := srv.Addrs // copy and remove duplicates
	tcpKeepAliveTimeout := srv.TCPKeepAliveTimeout
	updateBytesReadFunc := srv.UpdateBytesReadFunc
	updateBytesWrittenFunc := srv.UpdateBytesWrittenFunc

	// Create new HTTP listener.
	var listener *httpListener
	listener, err = NewHTTPListener(
		addrs,
		tlsConfig,
		tcpKeepAliveTimeout,
		readTimeout,
		writeTimeout,
		updateBytesReadFunc,
		updateBytesWrittenFunc,
	)
	if err != nil {
		return err
	}

	// Wrap given handler to do additional
	// * return 503 (service unavailable) if the server in shutdown.
	wrappedHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&srv.requestCount, 1)
		defer atomic.AddInt32(&srv.requestCount, -1)

		// If server is in shutdown, return 503 (service unavailable)
		if atomic.LoadUint32(&srv.inShutdown) != 0 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		// Handle request using passed handler.
		handler.ServeHTTP(w, r)
	})

	srv.listenerMutex.Lock()
	srv.Handler = wrappedHandler
	srv.listener = listener
	srv.listenerMutex.Unlock()

	// Start servicing with listener.
	return srv.Server.Serve(listener)
}

// Shutdown - shuts down HTTP server.
func (srv *Server) Shutdown() error {
	srv.listenerMutex.Lock()
	if srv.listener == nil {
		srv.listenerMutex.Unlock()
		return errors.New("server not initialized")
	}
	srv.listenerMutex.Unlock()

	if atomic.AddUint32(&srv.inShutdown, 1) > 1 {
		// shutdown in progress
		return errors.New("http server already in shutdown")
	}

	// Close underneath HTTP listener.
	srv.listenerMutex.Lock()
	err := srv.listener.Close()
	srv.listenerMutex.Unlock()

	// Wait for opened connection to be closed up to Shutdown timeout.
	shutdownTimeout := srv.ShutdownTimeout
	shutdownTimer := time.NewTimer(shutdownTimeout)
	ticker := time.NewTicker(serverShutdownPoll)
	defer ticker.Stop()
	for {
		select {
		case <-shutdownTimer.C:
			return errors.New("timed out. some connections are still active. doing abnormal shutdown")
		case <-ticker.C:
			if atomic.LoadInt32(&srv.requestCount) <= 0 {
				return err
			}
		}
	}
}

// NewServer - creates new HTTP server using given arguments.
func NewServer(addrs string, handler http.Handler) *Server {

	tlsConfig := &tls.Config{
		// TLS hardening
		PreferServerCipherSuites: false,
		MinVersion:               tls.VersionTLS12,
		NextProtos:               []string{"http/1.1", "h2"},
	}

	httpServer := &Server{
		Addrs:               addrs,
		ShutdownTimeout:     DefaultShutdownTimeout,
		TCPKeepAliveTimeout: DefaultTCPKeepAliveTimeout,
		listenerMutex:       &sync.Mutex{},
	}
	httpServer.Handler = handler
	httpServer.TLSConfig = tlsConfig
	httpServer.ReadTimeout = DefaultReadTimeout
	httpServer.WriteTimeout = DefaultWriteTimeout
	httpServer.MaxHeaderBytes = DefaultMaxHeaderBytes

	return httpServer
}
