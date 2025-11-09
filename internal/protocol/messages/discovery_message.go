package messages

import (
	"fmt"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/constants"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/encoding"
)

// discoveryMessage represents a DISCOVERY message (Type DISCOVERY).
//
// Purpose: UDP multicast discovery. Payload is empty.
// TLV: [Type:DISCOVERY][Length:0][Value:(empty)]
type discoveryMessage struct{}

// NewDiscoveryMessage creates a new DISCOVERY message.
func NewDiscoveryMessage() *discoveryMessage {
	return &discoveryMessage{}
}

// Encode encodes the DISCOVERY message to bytes.
func (m *discoveryMessage) Encode() ([]byte, error) {
	tlv, err := encoding.NewTLV(uint8(constants.DISCOVERY), []byte{})
	if err != nil {
		return nil, err
	}
	return tlv.Encode(), nil
}

// Type returns the message type.
func (m *discoveryMessage) Type() constants.MessageType {
	return constants.DISCOVERY
}

// DecodeDiscoveryMessage decodes a DISCOVERY message from TLV.
// Validates type (DISCOVERY) and that length is zero.
func DecodeDiscoveryMessage(tlv encoding.TLV) (discoveryMessage, error) {
	if tlv == nil {
		return discoveryMessage{}, fmt.Errorf("TLV is nil")
	}
	if tlv.Type() != uint8(constants.DISCOVERY) {
		return discoveryMessage{}, fmt.Errorf("expected DISCOVERY type, got %x", tlv.Type())
	}
	if tlv.Length() != 0 {
		return discoveryMessage{}, fmt.Errorf("DISCOVERY message should have empty value, got length %d", tlv.Length())
	}
	return discoveryMessage{}, nil
}
