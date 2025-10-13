package messages

import (
	"fmt"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/constants"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/tlv"
)

// ackErrorMessage represents an ACK or ERROR message (ACK_ERROR).
//
// Purpose: Reply with ACK_SUCCESS (0xA0) or an error code (0xA1+). Optional data
// may include additional info.
type ackErrorMessage struct {
	Code constants.AckErrorCode
	Data *string // Optional data
}

// NewAckMessage creates a new ACK message with code ACK_SUCCESS.
// Optional data may be nil or empty.
func NewAckMessage(data *string) *ackErrorMessage {
	return &ackErrorMessage{
		Code: constants.ACK_SUCCESS,
		Data: data,
	}
}

// NewErrorMessage creates a new ERROR message with the specified error code.
// Optional errorMessage may be nil or empty.
func NewErrorMessage(errorCode constants.AckErrorCode, errorMessage *string) (*ackErrorMessage, error) {
	if !errorCode.IsValid() || errorCode == constants.ACK_SUCCESS {
		return nil, fmt.Errorf("invalid error code: %v", errorCode)
	}
	return &ackErrorMessage{
		Code: errorCode,
		Data: errorMessage,
	}, nil
}

// Encode encodes the ACK/ERROR message to bytes.
// TLV: [Type:ACK_ERROR][Length:N][Value: Code(1) + Data(N-1)]
func (m *ackErrorMessage) Encode() ([]byte, error) {
	value := make([]byte, 1+len(*m.Data))
	copy(value[1:], *m.Data)
	value[0] = uint8(m.Code)

	tlv, err := tlv.NewTLV(uint8(constants.ACK_ERROR), value)
	if err != nil {
		return nil, fmt.Errorf("failed to create ACK_ERROR TLV: %w", err)
	}
	return tlv.Encode(), nil
}

// Type returns the message type.
func (m *ackErrorMessage) Type() constants.MessageType {
	return constants.ACK_ERROR
}

// IsError reports whether this is an error message (code != ACK_SUCCESS).
func (m *ackErrorMessage) IsError() bool {
	return m.Code != constants.ACK_SUCCESS
}

// DecodeAckErrorMessage decodes an ACK/ERROR message (Type ACK_ERROR) from TLV.
// Validates type, minimum length (1), and code range.
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

	var data *string
	if len(tlv.Value()) > 1 {
		temp := string(tlv.Value()[1:])
		data = &temp
	}

	return ackErrorMessage{
		Code: code,
		Data: data,
	}, nil
}
