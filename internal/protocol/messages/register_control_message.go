package messages

import (
	"fmt"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/constants"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/tlv"
)

// registerControlMessage represents a REGISTER_CONTROL message (Type 0x11)
type registerControlMessage struct{}

// NewRegisterControlMessage creates a new REGISTER_CONTROL message
func NewRegisterControlMessage() *registerControlMessage {
	return &registerControlMessage{}
}

// Encode encodes the REGISTER_CONTROL message to bytes
func (m *registerControlMessage) Encode() ([]byte, error) {
	tlv, err := tlv.NewTLV(uint8(constants.REGISTER_CONTROL), []byte{})
	if err != nil {
		return nil, fmt.Errorf("failed to create REGISTER_CONTROL TLV: %w", err)
	}
	return tlv.Encode(), nil
}

// GetType returns the message type
func (m *registerControlMessage) Type() constants.MessageType {
	return constants.REGISTER_CONTROL
}

// DecodeRegisterControlMessage decodes a REGISTER_CONTROL message from TLV
func DecodeRegisterControlMessage(tlv tlv.TLV) (registerControlMessage, error) {
	if tlv.Type() != uint8(constants.REGISTER_CONTROL) {
		return registerControlMessage{}, fmt.Errorf("expected REGISTER_CONTROL type, got %x", tlv.Type())
	}
	if tlv.Length() != 0 {
		return registerControlMessage{}, fmt.Errorf("REGISTER_CONTROL message should have empty value, got length %d", tlv.Length())
	}
	return registerControlMessage{}, nil
}
