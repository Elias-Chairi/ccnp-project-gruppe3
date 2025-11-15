package messages

import (
	"fmt"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/entity"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/constants"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/encoding"
)

// ackErrorMessage represents an ACK or ERROR message (ACK_ERROR).
//
// Purpose: Reply with ACK_SUCCESS (0xA0) or an error code (0xA1+). Optional data
// may include additional info.
type AckErrorMessage struct {
	Code constants.AckErrorCode
	Data string // Optional additional data
}

// AckSuccessMessage creates a standard ACK_SUCCESS message.
func AckSuccessMessage() AckErrorMessage {
	return AckErrorMessage{
		Code: constants.ACK_SUCCESS,
		Data: "",
	}
}

// NewAckNodeListMessage creates an ACK_SUCCESS message containing a list of nodes.
func NewAckNodeListMessage(nodes []entity.Node) (*AckErrorMessage, error) {
	var tlvs []encoding.TLV
	for _, node := range nodes {
		tlv, err := encoding.EncodeNodeEntry(node)
		if err != nil {
			return nil, fmt.Errorf("failed to encode node entry: %w", err)
		}
		tlvs = append(tlvs, tlv)
	}

	value := encoding.EncodeMultipleTLVs(tlvs)
	return &AckErrorMessage{
		Code: constants.ACK_SUCCESS,
		Data: string(value),
	}, nil
}

func (m AckErrorMessage) Type() constants.MessageType {
	return constants.ACK_ERROR
}

// Encode encodes the ACK/ERROR message to bytes.
// TLV: [Type:ACK_ERROR][Length:N][Value: Code(1) + Data(N-1)]
func (m AckErrorMessage) Encode() (encoding.TLV, error) {
	value := make([]byte, 1+len(m.Data))
	value[0] = uint8(m.Code)
	copy(value[1:], m.Data)

	tlv, err := encoding.NewTLV(uint8(constants.ACK_ERROR), value)
	if err != nil {
		return nil, fmt.Errorf("failed to create ACK_ERROR TLV: %w", err)
	}
	return tlv, nil
}

// IsError reports whether this is an error message (code != ACK_SUCCESS).
func (m *AckErrorMessage) IsError() bool {
	return m.Code != constants.ACK_SUCCESS
}

// DecodeAckErrorMessage decodes an ACK/ERROR message (Type ACK_ERROR) from TLV.
// Validates type, minimum length (1), and code range.
func DecodeAckErrorMessage(tlv encoding.TLV) (*AckErrorMessage, error) {
	if tlv == nil {
		return nil, fmt.Errorf("TLV is nil")
	}
	if tlv.Type() != uint8(constants.ACK_ERROR) {
		return nil, fmt.Errorf("expected ACK_ERROR type, got %x", tlv.Type())
	}
	if tlv.Length() < 1 {
		return nil, fmt.Errorf("ACK_ERROR message should have at least 1 byte for code")
	}

	code := constants.AckErrorCode(tlv.Value()[0])
	if !code.IsValid() {
		return nil, fmt.Errorf("invalid ACK_ERROR code: %x", tlv.Value()[0])
	}

	return &AckErrorMessage{
		Code: code,
		Data: string(tlv.Value()[1:]),
	}, nil
}
