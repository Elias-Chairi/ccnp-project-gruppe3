package connectionhandler

import (
	"fmt"
	"log"
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
func (c *ConnectionHandler) Register(node entity.Node) (uint8, error) {
	msg, err := messages.NewRegisterNodeMessage(node.Sensors, node.Actuators)
	if err != nil {
		return 0, err
	}
	encoded, err := msg.Encode()
	if err != nil {
		return 0, fmt.Errorf("failed to encode: %w", err)
	}

	if _, err := c.conn.Write(encoded); err != nil {
		return 0, fmt.Errorf("failed to send: %w", err)
	}

	tlv, _, err := encoding.ReadNextMessage(c.conn, time.Second*2)
	if err != nil {
		return 0, fmt.Errorf("failed to read next message: %w", err)
	}

	ack, err := messages.DecodeAckErrorMessage(tlv)
	if err != nil {
		return 0, fmt.Errorf("failed to decode ack error message: %w", err)
	}

	if ack.IsError() {
		return 0, fmt.Errorf("errorcode %d: register response is error: %s", ack.Code, ack.Data)
	}

	inner, err := encoding.DecodeTLV([]byte(ack.Data))
	if err != nil {
		return 0, fmt.Errorf("failed to decode inner TLVs: %w", err)
	}

	if inner.Length() != 1 {
		return 0, fmt.Errorf("invalid NodeID TLV length")
	}

	assignedID := inner.Value()[0]

	log.Printf("Node successfully registered. Assigned NodeID = %d\n", assignedID)

	return assignedID, nil
}

// SendSensorUpdate sends a sensor update message from a node to a server
func (c *ConnectionHandler) SendSensorUpdate(sensor entity.Sensor[any]) error {
	// Build message (node → server, so no nodeID is needed)
	msg := messages.NewSensorUpdateMessage(sensor)

	// Encode into raw TLV bytes
	tlv, err := msg.Encode()
	if err != nil {
		return fmt.Errorf("failed to encode sensor update: %w", err)
	}

	// Write to TCP connection
	_, err = c.conn.Write(tlv.Encode())
	if err != nil {
		return fmt.Errorf("failed to send sensor update: %w", err)
	}

	return nil
}
