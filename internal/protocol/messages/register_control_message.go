package messages

import (
	"fmt"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/constants"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/encoding"
)

// registerControlMessage represents a REGISTER_CONTROL message (Type REGISTER_CONTROL).
//
// Purpose: Control panel registers with server over TCP. Payload is empty.
// TLV: [Type:REGISTER_CONTROL][Length:0][Value:(empty)]
type registerControlMessage struct{}

// NewRegisterControlMessage creates a new REGISTER_CONTROL message.
func NewRegisterControlMessage() *registerControlMessage {
	return &registerControlMessage{}
}

// Encode encodes the REGISTER_CONTROL message to bytes.
func (m *registerControlMessage) Encode() ([]byte, error) {
	tlv, err := encoding.NewTLV(uint8(constants.REGISTER_CONTROL), []byte{})
	if err != nil {
		return nil, fmt.Errorf("failed to create REGISTER_CONTROL TLV: %w", err)
	}
	return tlv.Encode(), nil
}

// Type returns the message type.
func (m *registerControlMessage) Type() constants.MessageType {
	return constants.REGISTER_CONTROL
}

// DecodeRegisterControlMessage decodes a REGISTER_CONTROL message from TLV.
// Validates type (REGISTER_CONTROL) and that length is zero.
func DecodeRegisterControlMessage(tlv encoding.TLV) (registerControlMessage, error) {
	if tlv == nil {
		return registerControlMessage{}, fmt.Errorf("nil TLV provided")
	}
	if tlv.Type() != uint8(constants.REGISTER_CONTROL) {
		return registerControlMessage{}, fmt.Errorf("expected REGISTER_CONTROL type, got %x", tlv.Type())
	}
	if tlv.Length() != 0 {
		return registerControlMessage{}, fmt.Errorf("REGISTER_CONTROL message should have empty value, got length %d", tlv.Length())
	}
	return registerControlMessage{}, nil
}
