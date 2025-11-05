package tcpservice

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
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

func StartTCPService() {

	listener, err := net.ListenTCP("tcp4", &tcpServiceAddress)
	if err != nil {
		log.Fatal("Error listening TCP:", err)
	}
	defer func() {
		_ = listener.Close()
	}()

	log.Printf("Listening on: %s\n", listener.Addr())

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Failed to accept connection:", err)
			continue
		}
		go func() {
			log.Println("New connection from:", conn.RemoteAddr())
			if err := handleRegistration(conn); err != nil {
				if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
					log.Println("Connection closed by client:", conn.RemoteAddr())
				} else {
					log.Println("Connection closed with unrecoverable error from", conn.RemoteAddr(), ":", err)
				}
			}
		}()
	}
}

// readNextTRLV reads the next TRLV from the connection.
func readNextTRLV(conn net.Conn) (*encoding.TRLV, error) {
	for {
		// dosent make sense to read forever since if the recived data is too far apart in time
		// it is not likely that they belong to the same message or that the client is dead.
		conn.SetReadDeadline(time.Now().Add(readDeadline))
		trlv, n, err := encoding.ReadTRLV(conn)
		if err != nil {
			if errors.Is(err, os.ErrDeadlineExceeded) && n == 0 {
				// timeout occurred without reading any data, continue reading holding the connection open indefinitely
				continue
			}

			// if an error occurs during the top-level TRLV read, it cannot continue processing
			// because it cannot determine if the next bytes belong to the current message or the next one.
			return nil, fmt.Errorf("error reading TRLV from connection: %w", err)
		}

		return trlv, nil
	}
}

// handleRegistration handles the registration process and all further communication with the client.
// Only returns when the connection is closed or an unrecoverable error occurs.
func handleRegistration(conn net.Conn) error {
	defer func() {
		_ = conn.Close()
	}()

	trlv, err := readNextTRLV(conn)
	if err != nil {
		return fmt.Errorf("error reading registration TRLV: %w", err)
	}

	// handle based on registration type
	switch trlv.TLV.Type() {
	case uint8(constants.REGISTER_NODE):
		_, err := messages.DecodeRegisterNodeMessage(trlv.TLV)
		if err != nil {
			return fmt.Errorf("error decoding register node message %w", err)
		}
		// todo: reply with assigned node ID
		// todo: send new node to control panel(s)
		handleConn(conn, handleNode)
	case uint8(constants.REGISTER_CONTROL):
		_, err := messages.DecodeRegisterControlMessage(trlv.TLV)
		if err != nil {
			return fmt.Errorf("error decoding register control panel message %w", err)
		}
		// todo: reply with current node list
		handleConn(conn, handleControlPanel)
	default:
		return fmt.Errorf("invalid registration type %v", trlv.TLV.Type())
	}

	return nil
}

// ConnHandler is a function that handles a TLV received from a connection.
type ConnHandler func(t encoding.TLV) (constants.AckErrorCode, error)

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
		t, err := encoding.DecodeTLV(buf[:n])
		if err != nil {
			sendError(conn, constants.ERR_MALFORMED_MESSAGE)
			continue // ignore malformed TLVs
		}
		if code, err := f(t); err != nil || code != constants.ACK_SUCCESS {
			sendError(conn, code)
		} else {
			sendAck(conn)
		}
	}
}

func sendError(conn net.Conn, code constants.AckErrorCode) {
	ackErrMsg, err := messages.NewErrorMessage(code, nil)
	if err != nil {
		return
	}
	encodedMsg, err := ackErrMsg.Encode()
	if err != nil {
		return
	}
	_, _ = conn.Write(encodedMsg)
}

func sendAck(conn net.Conn) {
	ackMsg := messages.NewAckMessage(nil)
	encodedMsg, err := ackMsg.Encode()
	if err != nil {
		return
	}
	_, _ = conn.Write(encodedMsg)
}
