package encoding

import (
	"encoding/binary"
	"fmt"
	"math"
)

// TLV is the interface representing a Type-Length-Value element.
type TLV interface {
	// Type returns the 1-byte type code.
	Type() uint8
	// Length returns the 2-byte length (number of bytes in Value).
	Length() uint16
	// Value returns the raw value bytes.
	Value() []byte
	// Encode returns the wire representation [Type][Length][Value].
	Encode() []byte
}

// tlv is the internal implementation of the TLV interface.
type tlv struct {
	tlvType uint8
	length  uint16
	value   []byte
}

// NewTLV creates a new TLV instance.
//
// It returns an error if the value length exceeds math.MaxUint16 (65535).
// The provided value slice is not copied; it is referenced by the TLV.
func NewTLV(msgType uint8, value []byte) (TLV, error) {
	if len(value) > math.MaxUint16 {
		return nil, fmt.Errorf("value length exceeds maximum of %d bytes", math.MaxUint16)
	}
	return &tlv{
		tlvType: msgType,
		length:  uint16(len(value)),
		value:   value,
	}, nil
}

func (t *tlv) Type() uint8 {
	return t.tlvType
}

func (t *tlv) Length() uint16 {
	return t.length
}

func (t *tlv) Value() []byte {
	return t.value
}

func (t *tlv) Encode() []byte {
	result := make([]byte, 3+len(t.value))
	result[0] = t.tlvType
	binary.BigEndian.PutUint16(result[1:3], t.length)
	copy(result[3:], t.value)
	return result
}

// DecodeTLV parses a single TLV from the provided byte slice.
//
// Returns an error if:
//   - fewer than 3 bytes are provided (incomplete header), or
//   - the buffer does not contain the full value bytes as indicated by Length.
//
// The returned TLV contains a defensive copy of the value bytes.
func DecodeTLV(data []byte) (TLV, error) {
	if len(data) < 3 {
		return nil, fmt.Errorf("insufficient data for TLV header")
	}

	msgType := data[0]
	length := binary.BigEndian.Uint16(data[1:3])

	if len(data) < int(3+length) {
		return nil, fmt.Errorf("insufficient data for TLV value")
	}

	value := make([]byte, length)
	copy(value, data[3:3+length])

	return &tlv{
		tlvType: msgType,
		length:  length,
		value:   value,
	}, nil
}

// DecodeMultipleTLVs parses a concatenated stream of TLVs from data.
// It iteratively decodes TLVs until the buffer is exhausted or fewer
// than 3 bytes remain (insufficient for a header). If any TLV fails
// to decode, the function returns the error and no partial result.
func DecodeMultipleTLVs(data []byte) ([]TLV, error) {
	var tlvs []TLV
	offset := 0

	for offset < len(data) && len(data[offset:]) >= 3 {
		tlv, err := DecodeTLV(data[offset:])
		if err != nil {
			return nil, fmt.Errorf("failed to decode TLV %d at offset %d: %w", len(tlvs)+1, offset, err)
		}

		tlvs = append(tlvs, tlv)
		offset += int(3 + tlv.Length())
	}

	return tlvs, nil
}

// EncodeMultipleTLVs encodes multiple TLVs and concatenates their wire
// representations in order.
func EncodeMultipleTLVs(tlvs []TLV) []byte {
	var result []byte
	for _, tlv := range tlvs {
		result = append(result, tlv.Encode()...)
	}
	return result
}
