package messages

import (
	"fmt"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/constants"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/tlv"
)

// ackErrorMessage represents an ACK or ERROR message (Type 0x40)
type ackErrorMessage struct {
	Code constants.AckErrorCode
	Data string // Optional data (e.g., TCP address for ACK, error message for ERROR)
}

// NewAckMessage creates a new ACK message
func NewAckMessage(data string) *ackErrorMessage {
	return &ackErrorMessage{
		Code: constants.ACK_SUCCESS,
		Data: data,
	}
}

// NewErrorMessage creates a new ERROR message
func NewErrorMessage(errorCode constants.AckErrorCode, errorMessage string) *ackErrorMessage {
	return &ackErrorMessage{
		Code: errorCode,
		Data: errorMessage,
	}
}

// Encode encodes the ACK/ERROR message to bytes
func (m *ackErrorMessage) Encode() ([]byte, error) {
	value := make([]byte, 1+len(m.Data))
	value[0] = uint8(m.Code)
	copy(value[1:], []byte(m.Data))

	tlv, err := tlv.NewTLV(uint8(constants.ACK_ERROR), value)
	if err != nil {
		return nil, fmt.Errorf("failed to create ACK_ERROR TLV: %w", err)
	}
	return tlv.Encode(), nil
}

// Type returns the message type
func (m *ackErrorMessage) Type() constants.MessageType {
	return constants.ACK_ERROR
}

// IsError returns true if this is an error message
func (m *ackErrorMessage) IsError() bool {
	return m.Code != constants.ACK_SUCCESS
}

// DecodeAckErrorMessage decodes an ACK/ERROR message from TLV
func DecodeAckErrorMessage(tlv tlv.TLV) (ackErrorMessage, error) {
	if tlv == nil {
		return ackErrorMessage{}, fmt.Errorf("TLV is nil")
	}
	if tlv.Type() != uint8(constants.ACK_ERROR) {
		return ackErrorMessage{}, fmt.Errorf("expected ACK_ERROR type, got %x", tlv.Type())
	}
	if tlv.Length() < 1 {
		return ackErrorMessage{}, fmt.Errorf("ACK_ERROR message should have at least 1 byte for code")
	}

	code := constants.AckErrorCode(tlv.Value()[0])
	if !code.IsValid() {
		return ackErrorMessage{}, fmt.Errorf("invalid ACK_ERROR code: %x", tlv.Value()[0])
	}
	data := ""
	if len(tlv.Value()) > 1 {
		data = string(tlv.Value()[1:])
	}

	return ackErrorMessage{
		Code: code,
		Data: data,
	}, nil
}
