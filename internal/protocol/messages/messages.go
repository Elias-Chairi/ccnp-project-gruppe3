// Package messages implements high-level TLV messages used by the protocol.
//
// Message types (MessageType 0x40-0x4F):
//   - DISCOVERY: UDP multicast discovery. Empty payload.
//   - REGISTER_NODE: Node → Server. Contains sensor and actuator entries.
//   - REGISTER_CONTROL: Control → Server. Empty payload.
//   - SENSOR_UPDATE: Sensor reading. Node→Server (no node ID) and Server→Control (with node ID).
//   - COMMAND: Actuator command. Control→Server (with node selector) and Server→Node (no node selector).
//   - ACK_ERROR: Acknowledgment (ACK_SUCCESS) or error response (ERR_*).
//
// All messages encode as a top-level TLV with the message type and a Value that
// contains zero or more nested TLVs. Decoders validate the type and inner TLVs.
package messages

import (
	"fmt"
	"io"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/constants"
)

type Message interface {
	Encode() ([]byte, error)
	Type() constants.MessageType
}

// WriteMessage encodes and writes a message to the provided writer.
func WriteMessage(w io.Writer, msg Message) error {
	encoded, err := msg.Encode()
	if err != nil {
		return fmt.Errorf("failed to encode message: %w", err)
	}

	_, err = w.Write(encoded)
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	return nil
}
