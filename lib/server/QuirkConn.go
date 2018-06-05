package server

import (
	"net"
	"time"
	"sync/atomic"
)

type QuirkConn struct {
	net.Conn
	hadReadDeadlineInPast int32 // atomic
}

// SetReadDeadline - implements a workaround of SetReadDeadline go bug
func (q *QuirkConn) SetReadDeadline(t time.Time) error {
	inPast := int32(0)
	if t.Before(time.Now()) {
		inPast = 1
	}
	atomic.StoreInt32(&q.hadReadDeadlineInPast, inPast)
	return q.Conn.SetReadDeadline(t)
}

// canSetReadDeadline - returns if it is safe to set a new
// read deadline without triggering golang/go#21133 issue.
func (q *QuirkConn) canSetReadDeadline() bool {
	return atomic.LoadInt32(&q.hadReadDeadlineInPast) != 1
}
