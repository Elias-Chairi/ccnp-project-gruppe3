package connectionhandler

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/entity"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/constants"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/encoding"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/messages"
	utilNet "github.com/Elias-Chairi/ccnp-project-gruppe3/internal/util/net"
)

// ConnectionHandler manages the connection to the server.
type ConnectionHandler struct {
	conn      *utilNet.SafeConn
	ServerIPs []net.IP
	Node      *entity.Node
}

// Connect establishes a TCP connection to the server.
func (c *ConnectionHandler) Connect() error {
	conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{IP: c.ServerIPs[0], Port: 6000})
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	c.conn = utilNet.NewSafeConn(conn)
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

// StartListeningForCommands starts a goroutine that listens for incoming command messages from the server
// and applies them to the given node's actuators.
func (c *ConnectionHandler) StartListeningForCommands() {
	for {
		// read forever incoming messages from server
		tlv, reqID, err := encoding.ReadNextMessage(c.conn, time.Second*10)
		if err != nil {
			if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
				log.Println("Connection closed by server")
				return
			}
			log.Println("Error reading message from server:", err)
			continue
		}
		switch tlv.Type() {
		case uint8(constants.COMMAND):
			cmdMsg, err := messages.DecodeCommandMessage(tlv, false)
			if err != nil {
				log.Println("Error decoding command message:", err)
				continue
			}

			switch cmdMsg.ActuatorSelector.Type {
			case constants.SINGLE_ACTUATOR:
				// find actuator in node
				var actuator *entity.Actuator[any]
				for i := range c.Node.Actuators {
					if c.Node.Actuators[i].ID == cmdMsg.ActuatorSelector.ActuatorIDs[0] {
						actuator = &c.Node.Actuators[i]
						break
					}
				}

				var msg messages.AckErrorRequestIDMessage
				if actuator == nil {
					msg = messages.AckErrorRequestIDMessage{
						Code: constants.ERR_UNKNOWN_ACTUATOR_ID,
						Data: fmt.Sprintf("Actuator ID %d not found", cmdMsg.ActuatorSelector.ActuatorIDs[0]),
					}
				} else {
					msg = messages.AckRequestIDSuccessMessage()
					// apply state to actuator
					actuator.State = cmdMsg.ActuatorState
					log.Printf("Applied new state to actuator ID %d: %+v\n", actuator.ID, actuator.State)
				}
				tlv, err := msg.Encode()
				if err != nil {
					log.Println("Error encoding ACK message for actuator command:", err)
					continue
				}
				_, err = c.conn.Write(tlv.EncodeWithRequestID(*reqID))
				if err != nil {
					log.Println("Error sending ACK message for actuator command:", err)
				}

			case constants.ACTUATOR_LIST:
				// maybe future implementation
			case constants.ACTUATOR_TYPE:
				// maybe future implementation
			case constants.ALL_ACTUATORS:
				// maybe future implementation
			}
		default:
			log.Printf("Received unknown message type: %d\n", tlv.Type())
			continue
		}
	}
}
