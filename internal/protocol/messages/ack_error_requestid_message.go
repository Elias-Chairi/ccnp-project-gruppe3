package messages

import (
	"fmt"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/constants"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/encoding"
)

// ackErrorMessage represents an ACK or ERROR message (ACK_ERROR).
//
// Purpose: Reply with ACK_SUCCESS (0xA0) or an error code (0xA1+). Optional data
// may include additional info.
type AckErrorRequestIDMessage struct {
	Code constants.AckErrorCode
	Data string // Optional additional data
}

// AckRequestIDSuccessMessage creates a standard ACK_SUCCESS message.
func AckRequestIDSuccessMessage() AckErrorRequestIDMessage {
	return AckErrorRequestIDMessage{
		Code: constants.ACK_SUCCESS,
		Data: "",
	}
}

// NewAckRequestIDMessage creates an ACK_SUCCESS message with the given data.
func NewAckRequestIDMessage(data string) AckErrorRequestIDMessage {
	return AckErrorRequestIDMessage{
		Code: constants.ACK_SUCCESS,
		Data: data,
	}
}

func (m AckErrorRequestIDMessage) Type() constants.MessageType {
	return constants.ACK_ERROR_REQUESTID
}

// Encode encodes the ACK/ERROR message to bytes.
// TLV: [Type:ACK_ERROR][Length:N][Value: Code(1) + Data(N-1)]
func (m AckErrorRequestIDMessage) Encode() (encoding.TLV, error) {
	value := make([]byte, 1+len(m.Data))
	value[0] = uint8(m.Code)
	copy(value[1:], m.Data)

	tlv, err := encoding.NewTLV(uint8(constants.ACK_ERROR_REQUESTID), value)
	if err != nil {
		return nil, fmt.Errorf("failed to create ACK_ERROR_REQUESTID TLV: %w", err)
	}
	return tlv, nil
}

// IsError reports whether this is an error message (code != ACK_SUCCESS).
func (m *AckErrorRequestIDMessage) IsError() bool {
	return m.Code != constants.ACK_SUCCESS
}

// DecodeAckErrorRequestIDMessage decodes an ACK/ERROR message (Type ACK_ERROR_REQUESTID) from TLV.
// Validates type, minimum length (1), and code range.
func DecodeAckErrorRequestIDMessage(tlv encoding.TLV) (*AckErrorRequestIDMessage, error) {
	if tlv == nil {
		return nil, fmt.Errorf("TLV is nil")
	}
	if tlv.Type() != uint8(constants.ACK_ERROR_REQUESTID) {
		return nil, fmt.Errorf("expected ACK_ERROR_REQUESTID type, got %x", tlv.Type())
	}
	if tlv.Length() < 1 {
		return nil, fmt.Errorf("ACK_ERROR_REQUESTID message should have at least 1 byte for code")
	}

	code := constants.AckErrorCode(tlv.Value()[0])
	if !code.IsValid() {
		return nil, fmt.Errorf("invalid ACK_ERROR_REQUESTID code: %x", tlv.Value()[0])
	}

	return &AckErrorRequestIDMessage{
		Code: code,
		Data: string(tlv.Value()[1:]),
	}, nil
}
