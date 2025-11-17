package connectionhandler

import (
	"fmt"
	"net"
	"time"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/entity"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/constants"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/encoding"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/messages"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/selectors"
	utilNet "github.com/Elias-Chairi/ccnp-project-gruppe3/internal/util/net"
)

type UI interface {
	SensorUpdate(nodeID uint8, sensorID uint8, value any)
	ActuatorUpdate(nodeID uint8, actuatorID uint8, state any)
	ActuatorCommandResponse(nodeID uint8, actuatorID uint8, state any, err error)
	NodeAdded(node entity.Node)
	NodeRemoved(nodeID uint8)
}

// connectionHandler manages the connection to the server.
type connectionHandler struct {
	UI         UI
	conn       net.Conn
	ServerIPs  []net.IP
	pendingReq *utilNet.PendingRequests
}

func NewConnectionHandler(serverIPs []net.IP, ui UI) *connectionHandler {
	return &connectionHandler{
		UI:         ui,
		ServerIPs:  serverIPs,
		pendingReq: utilNet.NewPendingRequests(),
	}
}

// Connect establishes a TCP connection to the server.
func (c *connectionHandler) Connect() error {
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
func (c *connectionHandler) Disconnect() error {
	return c.conn.Close()
}

// Register sends a registration message to the server and waits for the response containing the list of nodes.
//
// Panics if the connection is not established.
func (c *connectionHandler) Register() (*[]entity.Node, error) {
	msg := messages.NewRegisterControlMessage()
	encoded, err := msg.Encode()
	if err != nil {
		return nil, fmt.Errorf("failed to encode: %w", err)
	}

	_, err = c.conn.Write(encoded)
	if err != nil {
		return nil, fmt.Errorf("failed to write to conn: %w", err)
	}

	readmsg, _, err := encoding.ReadNextMessage(c.conn, time.Second*10)

	if err != nil {
		return nil, fmt.Errorf("failed to read next message: %w", err)
	}

	decodedmsg, err := messages.DecodeAckErrorMessage(readmsg)
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

	go c.startListening()
	return &nodes, nil
}

func (c *connectionHandler) SendCommand(nodeID uint8, actuatorID uint8, state any) error {
	nodeSel := selectors.NewSingleNodeSelector(nodeID)
	actuatorSel := selectors.NewSingleActuatorSelector(actuatorID)
	msg := messages.NewCommandMessage(nodeSel, actuatorSel, state)

	tlv, err := msg.Encode()
	if err != nil {
		return fmt.Errorf("failed to encode command message: %w", err)
	}

	reqID := c.pendingReq.Add(utilNet.Request{Msg: msg})

	_, err = c.conn.Write(tlv.EncodeWithRequestID(reqID))
	if err != nil {
		c.pendingReq.Remove(reqID)
		return fmt.Errorf("failed to write to conn: %w", err)
	}
	return nil
}

func (c *connectionHandler) startListening() {
	for {
		tlv, reqID, err := encoding.ReadNextMessage(c.conn, time.Second*10)
		if err != nil {
			continue
		}

		switch tlv.Type() {
		case uint8(constants.NODE_ADDED):
			msg, err := messages.DecodeNodeAddedMessage(tlv)
			if err != nil {
				// malformed node added message, ignore
				continue
			}

			c.UI.NodeAdded(msg.Node)
		case uint8(constants.NODE_REMOVED):
			msg, err := messages.DecodeNodeRemovedMessage(tlv)
			if err != nil {
				// malformed node removed message, ignore
				continue
			}

			c.UI.NodeRemoved(msg.NodeID)
		case uint8(constants.SENSOR_UPDATE):
			msg, err := messages.DecodeSensorUpdateMessage(tlv, true)
			if err != nil {
				// malformed sensor update message, ignore
				continue
			}

			c.UI.SensorUpdate(*msg.NodeID, msg.Sensor.ID, msg.Sensor.Value)
		case uint8(constants.ACTUATOR_UPDATE):
			msg, err := messages.DecodeActuatorUpdateMessage(tlv)
			if err != nil {
				// malformed actuator update message, ignore
				continue
			}

			switch msg.ActuatorSelector.Type {
			case constants.SINGLE_ACTUATOR:
				c.UI.ActuatorUpdate(msg.NodeSelector.NodeIDs[0], msg.ActuatorSelector.ActuatorIDs[0], msg.ActuatorState)
			case constants.ACTUATOR_LIST:
				// maybe future functionality
			case constants.ACTUATOR_TYPE:
				// maybe future functionality
			case constants.ALL_ACTUATORS:
				// maybe future functionality
			}

		case uint8(constants.ACK_ERROR_REQUESTID):
			msg, err := messages.DecodeAckErrorRequestIDMessage(tlv)
			if err != nil {
				// malformed ack/error message, ignore
				continue
			}

			req, ok := c.pendingReq.Get(*reqID)
			if !ok {
				// unknown request ID, ignore
				continue
			}

			switch reqMsg := req.Msg.(type) {
			case *messages.CommandMessage:
				if msg.IsError() {
					c.UI.ActuatorCommandResponse(
						reqMsg.NodeSelector.NodeIDs[0],
						reqMsg.ActuatorSelector.ActuatorIDs[0],
						reqMsg.ActuatorState,
						fmt.Errorf("errorcode %d: command response is error: %s", msg.Code, msg.Data))
					break
				}
				c.UI.ActuatorCommandResponse(
					reqMsg.NodeSelector.NodeIDs[0],
					reqMsg.ActuatorSelector.ActuatorIDs[0],
					reqMsg.ActuatorState,
					nil)
			default:
				// unknown original message type, ignore
				continue
			}

			c.pendingReq.Remove(*reqID)
		default:
			// unknown message type, ignore
			continue
		}
	}
}
