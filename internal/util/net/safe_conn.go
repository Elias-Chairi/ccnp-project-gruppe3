package utilnet

import (
	"net"
	"sync"
	"time"
)

type SafeConn struct {
	conn    net.Conn
	writeMu sync.Mutex
	readMu  sync.Mutex
}

func NewSafeConn(conn net.Conn) *SafeConn {
	return &SafeConn{
		conn: conn,
	}
}

func (s *SafeConn) Write(data []byte) (int, error) {
	s.writeMu.Lock()
	defer s.writeMu.Unlock()
	return s.conn.Write(data)
}

func (s *SafeConn) Read(data []byte) (int, error) {
	s.readMu.Lock()
	defer s.readMu.Unlock()
	return s.conn.Read(data)
}

func (s *SafeConn) Close() error {
	return s.conn.Close()
}

func (s *SafeConn) RemoteAddr() net.Addr {
	return s.conn.RemoteAddr()
}

func (s *SafeConn) SetReadDeadline(t time.Time) error {
	s.readMu.Lock()
	defer s.readMu.Unlock()
	return s.conn.SetReadDeadline(t)
}
