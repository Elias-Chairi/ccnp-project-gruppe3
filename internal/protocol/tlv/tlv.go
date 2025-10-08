package tlv

import (
	"encoding/binary"
	"fmt"
	"math"
)

type TLV interface {
	Type() uint8
	Length() uint16
	Value() []byte
	Encode() []byte
}

// tlv represents a Type-Length-Value structure
type tlv struct {
	tlvType uint8  // 1 byte
	length  uint16 // 2 bytes
	value   []byte // variable length
}

// NewTLV creates a new TLV instance
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

// Type returns the TLV type
func (t *tlv) Type() uint8 {
	return t.tlvType
}

// Length returns the TLV length
func (t *tlv) Length() uint16 {
	return t.length
}

// Value returns the TLV value
func (t *tlv) Value() []byte {
	return t.value
}

// Encode serializes the TLV to bytes
func (t *tlv) Encode() []byte {
	result := make([]byte, 3+len(t.value))
	result[0] = t.tlvType
	binary.BigEndian.PutUint16(result[1:3], t.length)
	copy(result[3:], t.value)
	return result
}

// DecodeTLV parses bytes into a TLV structure
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

// DecodeMultipleTLVs parses multiple TLVs from a byte slice
func DecodeMultipleTLVs(data []byte) ([]TLV, error) {
	var tlvs []TLV
	offset := 0

	for offset < len(data) && len(data[offset:]) >= 3 {
		tlv, err := DecodeTLV(data[offset:])
		if err != nil {
			return nil, err
		}

		tlvs = append(tlvs, tlv)
		offset += int(3 + tlv.Length())
	}

	return tlvs, nil
}

// EncodeMultipleTLVs encodes multiple TLVs into a single byte slice
func EncodeMultipleTLVs(tlvs []TLV) []byte {
	var result []byte
	for _, tlv := range tlvs {
		result = append(result, tlv.Encode()...)
	}
	return result
}
