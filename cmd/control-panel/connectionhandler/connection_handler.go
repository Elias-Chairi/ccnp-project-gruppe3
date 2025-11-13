package connectionhandler

import (
	"fmt"
	"net"
)


type ConnectionHandler struct {
	conn net.Conn
	ServerIPs []net.IP 
}

func (c *ConnectionHandler) Connect() error {
	conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{IP: c.ServerIPs[0], Port: 6000})
	if (err != nil) {
		return fmt.Errorf("Failed to connect: %w", err)
	}
	c.conn = conn
	return nil
}
















