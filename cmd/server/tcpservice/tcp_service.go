package tcpservice

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

const readDeadline = 5 * time.Second

type tcpService struct {
	IP      net.IP
	Port    int
	Timeout time.Duration

	pendingReq *utilNet.PendingRequests
	nodeReg    *NodeRegistry
	ctrlPanReg *ControlPanelRegistry
}

// NewTcpService creates a new TCP service with the given parameters.
func NewTcpService(ip net.IP, port int, timeout time.Duration) *tcpService {
	return &tcpService{
		IP:      ip,
		Port:    port,
		Timeout: timeout,
	}
}

// Start starts the TCP service, listening for incoming connections and handling them.
func (t *tcpService) Start() {
	listener, err := net.ListenTCP("tcp4", &net.TCPAddr{IP: t.IP, Port: t.Port})
	if err != nil {
		log.Fatal("Error listening TCP:", err)
	}
	defer func() {
		_ = listener.Close()
	}()

	t.nodeReg = NewNodeRegistry()
	t.ctrlPanReg = &ControlPanelRegistry{}
	t.pendingReq = utilNet.NewPendingRequests()

	log.Printf("Listening on: %s\n", listener.Addr())

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Failed to accept connection:", err)
			continue
		}
		safeConn := utilNet.NewSafeConn(conn)
		go func() {
			log.Println("New connection from:", safeConn.RemoteAddr())
			if err := t.handleRegistration(safeConn); err != nil {
				if isConnClosedErr(err) {
					log.Println("Connection closed by client:", safeConn.RemoteAddr())
				} else {
					log.Println("Connection closed with unrecoverable error from", safeConn.RemoteAddr(), ":", err)
				}
			}
		}()
	}
}

// isConnClosedErr checks if the error indicates that the connection has been closed by the client.
func isConnClosedErr(err error) bool {
	return errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF)
}

// handleRegistration handles the registration process and all further communication with the client.
// Only returns when the connection is closed or an unrecoverable error occurs.
func (t *tcpService) handleRegistration(conn *utilNet.SafeConn) error {
	defer func() {
		_ = conn.Close()
	}()

	tlv, _, err := encoding.ReadNextMessage(conn, readDeadline)
	if err != nil {
		return fmt.Errorf("error reading registration message: %w", err)
	}

	// handle based on registration type
	switch tlv.Type() {
	case uint8(constants.REGISTER_NODE):
		msg, err := messages.DecodeRegisterNodeMessage(tlv)
		if err != nil {
			return fmt.Errorf("error decoding register node message %w", err)
		}

		// create and store unique node ID
		id := t.nodeReg.CreateNodeID(conn, msg.Sensors, msg.Actuators)
		// send assigned node ID back to node
		nodeID, err := encoding.NewTLV(uint8(constants.NODE_ID), []byte{byte(id)})
		if err != nil {
			return fmt.Errorf("error creating node ID TLV: %w", err)
		}
		err = writeMessage(conn, nil, messages.NewAckMessage(string(nodeID.Encode())))
		if err != nil {
			return fmt.Errorf("error sending assigned node ID to node: %w", err)
		}
		log.Printf("Node successfully registered. Assigned NodeID = %d\n", id)

		// notify all control panels about new node
		updateMsg := messages.NewNodeAddedMessage(entity.Node{
			ID:        id,
			Sensors:   msg.Sensors,
			Actuators: msg.Actuators,
		})
		t.NotifyAllControlPanels(updateMsg)

		defer func() {
			t.nodeReg.RemoveNodeID(id)
			// notify all control panels about node removal
			removeMsg := messages.NewNodeRemovedMessage(id)
			t.NotifyAllControlPanels(removeMsg)
		}()
		return t.handleConn(conn, handleNode)

	case uint8(constants.REGISTER_CONTROL):
		_, err := messages.DecodeRegisterControlMessage(tlv)
		if err != nil {
			return fmt.Errorf("error decoding register control panel message %w", err)
		}
		t.ctrlPanReg.AddControlPanel(conn)
		defer t.ctrlPanReg.RemoveControlPanel(conn)

		msg, err := messages.NewAckNodeListMessage(t.nodeReg.GetAllNodes())
		if err == nil {
			_ = writeMessage(conn, nil, msg) // ignoring error sending response
		}
		return t.handleConn(conn, handleControlPanel)
	default:
		return fmt.Errorf("invalid registration type %v", tlv.Type())
	}
}

func (t *tcpService) NotifyAllControlPanels(msg messages.TopLevelMessage) {
	for _, c := range t.ctrlPanReg.GetAll() {
		_ = writeMessage(c, nil, msg)
	}
}

// messageHandler is a function that handles a single top-level message from a connection.
type messageHandler func(t *tcpService, safeConn *utilNet.SafeConn, tlv encoding.TLV, reqID *uint16)

// Generic connection handler.
// Reads TLVs from the connection and passes them to the provided handler function.
// Read loop continues until the connection is closed.
func (t *tcpService) handleConn(conn *utilNet.SafeConn, handle messageHandler) error {
	for {
		msg, reqID, err := encoding.ReadNextMessage(conn, readDeadline)
		if err != nil {
			if !isConnClosedErr(err) {
				// send error message before closing connection
				_ = writeMessage(conn, nil, messages.AckErrorRequestIDMessage{
					Code: constants.ERR_MALFORMED_MESSAGE,
				})
			}
			return err
		}
		handle(t, conn, msg, reqID)
	}
}

func writeMessage(conn *utilNet.SafeConn, reqID *uint16, msg messages.TopLevelMessage) error {
	tlv, err := msg.Encode()
	if err != nil {
		return fmt.Errorf("error encoding response message: %w", err)
	}

	var data []byte
	if reqID != nil {
		data = tlv.EncodeWithRequestID(*reqID)
	} else {
		data = tlv.Encode()
	}

	_, err = conn.Write(data)
	if err != nil {
		return fmt.Errorf("error writing response message: %w", err)
	}
	return nil
}
