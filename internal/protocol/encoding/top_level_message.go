package encoding

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/constants"
)

// readMessage reads a Message from the given reader.
//
// It first reads until it has received the message type (1 byte), if message type requires a RequestID it reads that (2 bytes),
// then reads the length (2 bytes), then reads until the full value is received based on the length.
//
// For messages not requiring a RequestID, the RequestID field in the returned Message will be nil.
//
// Returns the Message, number of bytes read, and an error if any.
func readMessage(r io.Reader) (*tlv, *uint16, int, error) {
	headerBuf := make([]byte, 1) // read Type(1)
	n, err := io.ReadFull(r, headerBuf)
	if err != nil {
		return nil, nil, n, fmt.Errorf("error reading message header: %w", err)
	}

	messageType := constants.MessageType(headerBuf[0])
	if !messageType.IsValid() {
		return nil, nil, n, fmt.Errorf("invalid message type: %d", messageType)
	}

	var requestID *uint16
	switch messageType { // these message types require a RequestID
	case constants.COMMAND, constants.ACK_ERROR_REQUESTID:
		requestIDBuf := make([]byte, 2) // read RequestID(2)
		m, err := io.ReadFull(r, requestIDBuf)
		n += m
		if err != nil {
			return nil, nil, n, fmt.Errorf("error reading message request ID: %w", err)
		}
		reqID := binary.BigEndian.Uint16(requestIDBuf)
		requestID = &reqID
	}

	lengthBuf := make([]byte, 2)
	m, err := io.ReadFull(r, lengthBuf) // read Length(2)
	n += m
	if err != nil {
		return nil, nil, n, fmt.Errorf("error reading message header: %w", err)
	}
	length := binary.BigEndian.Uint16(lengthBuf)

	value := make([]byte, length)
	m, err = io.ReadFull(r, value) // read Value(N)
	n += m
	if err != nil {
		return nil, nil, n, fmt.Errorf("error reading message value: %w", err)
	}

	return &tlv{
		tlvType: uint8(messageType),
		length:  length,
		value:   value,
	}, requestID, n, nil
}

type MessageReader interface {
	io.Reader
	SetReadDeadline(time.Time) error
}

// ReadNextMessage reads the next Message from the given net.Conn with a specified timeout for each read operation.
func ReadNextMessage(conn MessageReader, timeout time.Duration) (*tlv, *uint16, error) {
	for {
		// doesn't make sense to read forever since if the received data is too far apart in time
		// it is not likely that they belong to the same message or that the client is dead.
		err := conn.SetReadDeadline(time.Now().Add(timeout))
		if err != nil {
			return nil, nil, fmt.Errorf("error setting read deadline: %w", err)
		}
		tlv, rID, n, err := readMessage(conn)
		if err != nil {
			if errors.Is(err, os.ErrDeadlineExceeded) && n == 0 {
				// timeout occurred without reading any data, continue reading holding the connection open indefinitely
				continue
			}

			// if an error occurs during the top-level Message read, it cannot continue processing
			// because it cannot determine if the next bytes belong to the current message or the next one.
			return nil, nil, fmt.Errorf("error reading Message from connection: %w", err)
		}

		return tlv, rID, nil
	}
}
