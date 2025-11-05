package tcpservice

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/constants"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/messages"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/tlv"
)

var tcpServiceAddress = net.TCPAddr{
	// localhost address
	IP: net.ParseIP("127.0.0.1"),
	// port sat in protocol document
	Port: 6000,
}

func StartTCPService() {

	listner, err := net.ListenTCP("tcp4", &tcpServiceAddress)
	if err != nil {
		log.Fatal("Error listening TCP:", err)
	}
	defer func() {
		_ = listner.Close()
	}()

	log.Printf("Listening on: %s\n", listner.Addr())

	for {
		conn, err := listner.Accept()
		if err != nil {
			log.Println("Failed to accept connection:", err)
			continue
		}
		go func() {
			log.Println("New connection from:", conn.RemoteAddr())
			if err := handleRegistration(conn); err != nil {
				log.Println("Error handling registration:", err)
			}
			log.Println("Connection closed:", conn.RemoteAddr())
		}()
	}
}

// handleRegistration handles the registration process and all further communication with the client.
// Only returns when the connection is closed or an unrecoverable error occurs.
func handleRegistration(conn net.Conn) error {
	defer func() {
		_ = conn.Close()
	}()

	// read first TLV to determine type of registration
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		_ = conn.Close()
		return fmt.Errorf("error reading conn bytes %w", err)
	}

	// decode TLV to get registration type
	t, err := tlv.DecodeTLV(buf[:n])
	if err != nil {
		return fmt.Errorf("error decoding TLV %w", err)
	}

	// handle based on registration type
	switch t.Type() {
	case uint8(constants.REGISTER_NODE):
		_, err := messages.DecodeRegisterNodeMessage(t)
		if err != nil {
			return fmt.Errorf("error decoding register node message %w", err)
		}
		// todo: reply with assigned node ID
		// todo: send new node to control panel(s)
		handleConn(conn, handleNode)
	case uint8(constants.REGISTER_CONTROL):
		_, err := messages.DecodeRegisterControlMessage(t)
		if err != nil {
			return fmt.Errorf("error decoding register control panel message %w", err)
		}
		// todo: reply with current node list
		handleConn(conn, handleControlPanel)
	default:
		return fmt.Errorf("invalid registration type %v", t.Type())
	}

	return nil
}

// ConnHandler is a function that handles a TLV received from a connection.
type ConnHandler func(t tlv.TLV) error

// Generic connection handler.
// Reads TLVs from the connection and passes them to the provided handler function.
// Read loop continues until the connection is closed.
func handleConn(conn net.Conn, f ConnHandler) {
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				return // connection closed by client
			} else {
				continue // ignore other read errors
			}
		}
		t, err := tlv.DecodeTLV(buf[:n])
		if err != nil {
			// todo: reply to client with error message
		}
		err = f(t)
		if err != nil {
			// todo: reply to client with error message
		} // todo(maybe): else ack success
	}
}
