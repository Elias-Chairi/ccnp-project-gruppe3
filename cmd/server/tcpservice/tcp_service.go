package tcpservice

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/constants"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/encoding"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/messages"
)

const readDeadline = 5 * time.Second

type tcpService struct {
	IP      net.IP
	Port    int
	Timeout time.Duration

	pendingReq *PendingRequests
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
	t.pendingReq = NewPendingRequests()

	log.Printf("Listening on: %s\n", listener.Addr())

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Failed to accept connection:", err)
			continue
		}
		go func() {
			log.Println("New connection from:", conn.RemoteAddr())
			if err := t.handleRegistration(conn); err != nil {
				if isConnClosedErr(err) {
					log.Println("Connection closed by client:", conn.RemoteAddr())
				} else {
					log.Println("Connection closed with unrecoverable error from", conn.RemoteAddr(), ":", err)
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
func (t *tcpService) handleRegistration(conn net.Conn) error {
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
		// todo: send new node to control panel(s)

		defer func() {
			t.nodeReg.RemoveNodeID(id)
			// todo: send delete node to control panel(s)
		}()
		return t.handleConn(t, conn, handleNode)

	case uint8(constants.REGISTER_CONTROL):
		_, err := messages.DecodeRegisterControlMessage(tlv)
		if err != nil {
			return fmt.Errorf("error decoding register control panel message %w", err)
		}
		t.ctrlPanReg.AddControlPanel(conn)
		defer t.ctrlPanReg.RemoveControlPanel(conn)

		msg, err := messages.NewAckNodeListMessage(t.nodeReg.GetAllNodes())
		if err == nil {
			_ = t.writeMessage(conn, nil, msg) // ignoring error sending response
		}
		return t.handleConn(t, conn, handleControlPanel)
	default:
		return fmt.Errorf("invalid registration type %v", tlv.Type())
	}
}

// messageHandler is a function that handles a single top-level message from a connection.
type messageHandler func(tcp *tcpService, tlv encoding.TLV) messages.TopLevelMessage

// Generic connection handler.
// Reads TLVs from the connection and passes them to the provided handler function.
// Read loop continues until the connection is closed.
func (t *tcpService) handleConn(tcp *tcpService, conn net.Conn, f messageHandler) error {
	for {
		msg, reqID, err := encoding.ReadNextMessage(conn, readDeadline)
		// todo: Request id ?????
		if err != nil {
			if !isConnClosedErr(err) {
				// send error message before closing connection
				_ = t.writeMessage(conn, nil, messages.AckErrorMessage{
					Code: constants.ERR_MALFORMED_MESSAGE,
				})
			}
			return err
		}
		// send response, ignoring errors
		_ = t.writeMessage(conn, reqID, f(tcp, msg))
	}
}

func (t *tcpService) writeMessage(conn net.Conn, reqID *uint16, msg messages.TopLevelMessage) error {
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

	// switch tlv.Type() {
	// case uint8(constants.COMMAND): // writing a command to node, store in pending requests
	// 	id := t.pendingReq.Add(msg)
	// 	data = tlv.EncodeWithRequestID(id)
	// case uint8(constants.ACK_ERROR_REQUESTID): // writing an ack/error with request ID 0 (back to control panel)
	// 	if reqID == nil {
	// 		return fmt.Errorf("missing request ID for ACK/ERROR message")
	// 	}
	// 	data = tlv.EncodeWithRequestID(*reqID)
	// default:
	// 	data = tlv.Encode()
	// }

	_, err = conn.Write(data)
	if err != nil {
		return fmt.Errorf("error writing response message: %w", err)
	}
	return nil
}
