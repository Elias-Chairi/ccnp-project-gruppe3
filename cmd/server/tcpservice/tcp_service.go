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

var tcpServiceAddress = net.TCPAddr{
	// localhost address
	IP: net.ParseIP("127.0.0.1"),
	// port sat in protocol document
	Port: 6000,
}

const readDeadline = 5 * time.Second

type TcpService struct {
	nodeReg    *NodeRegistry
	ctrlPanReg *ControlPanelRegistry
}

func (t *TcpService) Start() {

	listener, err := net.ListenTCP("tcp4", &tcpServiceAddress)
	if err != nil {
		log.Fatal("Error listening TCP:", err)
	}
	defer func() {
		_ = listener.Close()
	}()

	t.nodeReg = NewNodeRegistry()
	t.ctrlPanReg = &ControlPanelRegistry{}

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
func (t *TcpService) handleRegistration(conn net.Conn) error {
	defer func() {
		_ = conn.Close()
	}()

	msg, err := encoding.ReadNextMessage(conn, readDeadline)
	if err != nil {
		return fmt.Errorf("error reading registration message: %w", err)
	}

	// handle based on registration type
	switch msg.TLV.Type() {
	case uint8(constants.REGISTER_NODE):
		msg, err := messages.DecodeRegisterNodeMessage(msg.TLV)
		if err != nil {
			return fmt.Errorf("error decoding register node message %w", err)
		}

		// create and store unique node ID
		id := t.nodeReg.CreateNodeID(conn, msg.Sensors, msg.Actuators)

		// when function returns, remove node ID from registry
		// todo: send new node to control panel(s)

		defer func() {
			t.nodeReg.RemoveNodeID(id)
			// todo: send delete node to control panel(s)
		}()
		return handleConn(conn, handleNode)

	case uint8(constants.REGISTER_CONTROL):
		_, err := messages.DecodeRegisterControlMessage(msg.TLV)
		if err != nil {
			return fmt.Errorf("error decoding register control panel message %w", err)
		}
		t.ctrlPanReg.AddControlPanel(conn)
		defer t.ctrlPanReg.RemoveControlPanel(conn)

		msg, err := messages.NewAckNodeListMessage(t.nodeReg.GetAllNodes())
		if err == nil {
			_ = messages.WriteMessage(conn, msg) // ignoring error sending response
		}
		return handleConn(conn, handleControlPanel)
	default:
		return fmt.Errorf("invalid registration type %v", msg.TLV.Type())
	}
}

// messageHandler is a function that handles a single top-level message from a connection.
type messageHandler func(topLevelMessage *encoding.Message) messages.AckErrorMessage

// Generic connection handler.
// Reads TLVs from the connection and passes them to the provided handler function.
// Read loop continues until the connection is closed.
func handleConn(conn net.Conn, f messageHandler) error {
	for {
		msg, err := encoding.ReadNextMessage(conn, readDeadline)
		if err != nil {
			if !isConnClosedErr(err) {
				// send error message before closing connection
				_ = messages.WriteMessage(conn, messages.AckErrorMessage{
					Code: constants.ERR_MALFORMED_MESSAGE,
				})
			}
			return err
		}
		// send response, ignoring errors
		_ = messages.WriteMessage(conn, f(msg))
	}
}
