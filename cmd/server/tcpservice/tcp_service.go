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

// Make map to assign node IDs to the net.Conn
var nodeIDs = make(map[int16]net.Conn)

func createNodeID(conn net.Conn) int16 {
	nodeID := int16(len(nodeIDs) + 1)
	nodeIDs[nodeID] = conn
	return nodeID
}

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

// readNextTRLV reads the next TRLV from the connection.
func readNextTRLV(conn net.Conn) (*encoding.TRLV, error) {
	for {
		// doesn't make sense to read forever since if the received data is too far apart in time
		// it is not likely that they belong to the same message or that the client is dead.
		err := conn.SetReadDeadline(time.Now().Add(readDeadline))
		if err != nil {
			return nil, fmt.Errorf("error setting read deadline: %w", err)
		}
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

		// Create and assign new node ID
		nodeID := createNodeID(conn)

		// Send ACK_SUCCESS response with the assigned node ID
		response := messages.AckErrorMessage{
			Code: constants.ACK_SUCCESS,
			Data: fmt.Sprintf("NodeID:%d", nodeID),
		}

		if err := sendResponse(conn, response); err != nil {
			return fmt.Errorf("failed to send node ID response: %w", err)
		}

		// todo: send new node to control panel(s)
		return handleConn(conn, handleNode)
	case uint8(constants.REGISTER_CONTROL):
		_, err := messages.DecodeRegisterControlMessage(trlv.TLV)
		if err != nil {
			return fmt.Errorf("error decoding register control panel message %w", err)
		}
		// todo: reply with current node list
		return handleConn(conn, handleControlPanel)
	default:
		return fmt.Errorf("invalid registration type %v", trlv.TLV.Type())
	}
}

// messageHandler is a function that handles a single top-level TRLV message from a connection.
type messageHandler func(topLevelMessage *encoding.TRLV) messages.AckErrorMessage

// Generic connection handler.
// Reads TLVs from the connection and passes them to the provided handler function.
// Read loop continues until the connection is closed.
func handleConn(conn net.Conn, f messageHandler) error {
	for {
		t, err := readNextTRLV(conn)
		if err != nil {
			if !isConnClosedErr(err) {
				_ = sendResponse(conn, messages.AckErrorMessage{
					Code: constants.ERR_MALFORMED_MESSAGE,
				})
			}
			return err
		}
		_ = sendResponse(conn, f(t)) // ignoring error sending response
	}
}

// sendResponse encodes and sends an AckErrorMessage response over the connection.
func sendResponse(conn net.Conn, message messages.AckErrorMessage) error {
	encodedMsg, err := message.Encode()
	if err != nil {
		return fmt.Errorf("failed to encode response message: %w", err)
	}
	_, err = conn.Write(encodedMsg)
	if err != nil {
		return fmt.Errorf("failed to send response message: %w", err)
	}
	return nil
}
