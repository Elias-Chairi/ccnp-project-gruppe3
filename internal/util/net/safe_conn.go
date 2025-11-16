package utilnet

import (
	"net"
	"sync"
	"time"
)

type SafeConn struct {
	conn net.Conn
	mu   sync.Mutex
}

func NewSafeConn(conn net.Conn) *SafeConn {
	return &SafeConn{
		conn: conn,
	}
}

func (sc *SafeConn) Write(data []byte) (int, error) {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	return sc.conn.Write(data)
}

func (sc *SafeConn) Read(data []byte) (int, error) {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	return sc.conn.Read(data)
}

func (sc *SafeConn) Close() error {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	return sc.conn.Close()
}

func (sc *SafeConn) RemoteAddr() net.Addr {
	return sc.conn.RemoteAddr()
}

func (sc *SafeConn) SetReadDeadline(t time.Time) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	return sc.conn.SetReadDeadline(t)
}
