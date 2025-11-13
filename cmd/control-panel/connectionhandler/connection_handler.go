package connectionhandler

import (
	"fmt"
	"net"
	"time"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/entity"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/encoding"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/messages"
)

// ConnectionHandler manages the connection to the server.
type ConnectionHandler struct {
	conn      net.Conn
	ServerIPs []net.IP
}

// Connect establishes a TCP connection to the server.
func (c *ConnectionHandler) Connect() error {
	conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{IP: c.ServerIPs[0], Port: 6000})
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	c.conn = conn
	return nil
}

// Disconnect closes the TCP connection to the server.
//
// Panics if the connection is not established.
func (c *ConnectionHandler) Disconnect() error {
	return c.conn.Close()
}

// Register sends a registration message to the server and waits for the response containing the list of nodes.
//
// Panics if the connection is not established.
func (c *ConnectionHandler) Register() (*[]entity.Node, error) {
	msg := messages.NewRegisterControlMessage()
	encoded, err := msg.Encode()
	if err != nil {
		return nil, fmt.Errorf("failed to encode: %w", err)
	}

	_, err = c.conn.Write(encoded)
	if err != nil {
		return nil, fmt.Errorf("failed to write to conn: %w", err)
	}

	readmsg, err := encoding.ReadNextMessage(c.conn, time.Second*10)

	if err != nil {
		return nil, fmt.Errorf("failed to read next message: %w", err)
	}

	decodedmsg, err := messages.DecodeAckErrorMessage(readmsg.TLV)
	if err != nil {
		return nil, fmt.Errorf("failed to decode ack error message: %w", err)
	}

	if decodedmsg.IsError() {
		return nil, fmt.Errorf("errorcode %d: register response is error: %s", decodedmsg.Code, decodedmsg.Data)
	}

	tlvs, err := encoding.DecodeMultipleTLVs([]byte(decodedmsg.Data))
	if err != nil {
		return nil, fmt.Errorf("failed to decode message: %w", err)
	}

	var nodes []entity.Node

	for _, tlv := range tlvs {
		node, err := encoding.DecodeNodeEntry(tlv)
		if err != nil {
			return nil, fmt.Errorf("failed to decode node entry: %w", err)
		}
		nodes = append(nodes, *node)
	}

	return &nodes, nil
}
