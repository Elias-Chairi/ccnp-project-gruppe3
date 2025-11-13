package encoding

import (
	"encoding/binary"
	"fmt"
	"math"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/constants"
)

// EncodeInteger32 encodes a 32-bit integer into a TLV of type DATA_TYPE_INTEGER
func EncodeInteger32(value int32) (TLV, error) {
	data := make([]byte, 4)
	binary.BigEndian.PutUint32(data, uint32(value))
	return NewTLV(uint8(constants.DATA_TYPE_INTEGER), data)
}

// EncodeInteger64 encodes a 64-bit integer into a TLV of type DATA_TYPE_INTEGER
func EncodeInteger64(value int64) (TLV, error) {
	data := make([]byte, 8)
	binary.BigEndian.PutUint64(data, uint64(value))
	return NewTLV(uint8(constants.DATA_TYPE_INTEGER), data)
}

// DecodeInteger32 decodes a TLV of type DATA_TYPE_INTEGER into a 32-bit integer.
func DecodeInteger32(tlv TLV) (int32, error) {
	if tlv == nil {
		return 0, fmt.Errorf("TLV is nil")
	}
	if tlv.Type() != uint8(constants.DATA_TYPE_INTEGER) {
		return 0, fmt.Errorf("expected integer type, got %x", tlv.Type())
	}
	if tlv.Length() != 4 {
		return 0, fmt.Errorf("invalid integer length: %d", tlv.Length())
	}
	return int32(binary.BigEndian.Uint32(tlv.Value())), nil
}

// DecodeInteger64 decodes a TLV of type DATA_TYPE_INTEGER into a 64-bit integer.
func DecodeInteger64(tlv TLV) (int64, error) {
	if tlv == nil {
		return 0, fmt.Errorf("TLV is nil")
	}
	if tlv.Type() != uint8(constants.DATA_TYPE_INTEGER) {
		return 0, fmt.Errorf("expected integer type, got %x", tlv.Type())
	}
	if tlv.Length() != 8 {
		return 0, fmt.Errorf("invalid integer length: %d", tlv.Length())
	}
	return int64(binary.BigEndian.Uint64(tlv.Value())), nil
}

// EncodeFloat32 encodes a 32-bit float into a TLV of type DATA_TYPE_FLOAT
func EncodeFloat32(value float32) (TLV, error) {
	data := make([]byte, 4)
	binary.BigEndian.PutUint32(data, math.Float32bits(value))
	return NewTLV(uint8(constants.DATA_TYPE_FLOAT), data)
}

// EncodeFloat64 encodes a 64-bit float into a TLV of type DATA_TYPE_FLOAT
func EncodeFloat64(value float64) (TLV, error) {
	data := make([]byte, 8)
	binary.BigEndian.PutUint64(data, math.Float64bits(value))
	return NewTLV(uint8(constants.DATA_TYPE_FLOAT), data)
}

// DecodeFloat32 decodes a TLV of type DATA_TYPE_FLOAT into a 32-bit float.
func DecodeFloat32(tlv TLV) (float32, error) {
	if tlv == nil {
		return 0, fmt.Errorf("TLV is nil")
	}
	if tlv.Type() != uint8(constants.DATA_TYPE_FLOAT) {
		return 0, fmt.Errorf("expected float type, got %x", tlv.Type())
	}
	if tlv.Length() != 4 {
		return 0, fmt.Errorf("invalid float length: %d", tlv.Length())
	}
	bits := binary.BigEndian.Uint32(tlv.Value())
	return math.Float32frombits(bits), nil
}

// DecodeFloat64 decodes a TLV of type DATA_TYPE_FLOAT into a 64-bit float.
func DecodeFloat64(tlv TLV) (float64, error) {
	if tlv == nil {
		return 0, fmt.Errorf("TLV is nil")
	}
	if tlv.Type() != uint8(constants.DATA_TYPE_FLOAT) {
		return 0, fmt.Errorf("expected float type, got %x", tlv.Type())
	}
	if tlv.Length() != 8 {
		return 0, fmt.Errorf("invalid float length: %d", tlv.Length())
	}
	bits := binary.BigEndian.Uint64(tlv.Value())
	return math.Float64frombits(bits), nil
}

// EncodeString encodes a string into a TLV of type DATA_TYPE_STRING
func EncodeString(value string) (TLV, error) {
	return NewTLV(uint8(constants.DATA_TYPE_STRING), []byte(value))
}

// DecodeString decodes a TLV of type DATA_TYPE_STRING into a string.
func DecodeString(tlv TLV) (string, error) {
	if tlv == nil {
		return "", fmt.Errorf("TLV is nil")
	}
	if tlv.Type() != uint8(constants.DATA_TYPE_STRING) {
		return "", fmt.Errorf("expected string type, got %x", tlv.Type())
	}
	return string(tlv.Value()), nil
}

// EncodeBoolean encodes a boolean into a TLV of type DATA_TYPE_BOOLEAN
func EncodeBoolean(value bool) (TLV, error) {
	var byteValue byte
	if value {
		byteValue = 0x01
	} else {
		byteValue = 0x00
	}
	return NewTLV(uint8(constants.DATA_TYPE_BOOLEAN), []byte{byteValue})
}

// DecodeBoolean decodes a TLV of type DATA_TYPE_BOOLEAN into a boolean (0x00=false, 0x01=true).
func DecodeBoolean(tlv TLV) (bool, error) {
	if tlv == nil {
		return false, fmt.Errorf("TLV is nil")
	}
	if tlv.Type() != uint8(constants.DATA_TYPE_BOOLEAN) {
		return false, fmt.Errorf("expected boolean type, got %x", tlv.Type())
	}
	if tlv.Length() != 1 {
		return false, fmt.Errorf("invalid boolean length: %d", tlv.Length())
	}
	switch tlv.Value()[0] {
	case 0x00:
		return false, nil
	case 0x01:
		return true, nil
	default:
		return false, fmt.Errorf("invalid boolean value: %d", tlv.Value()[0])
	}
}

// EncodeAny encodes a supported Go value as the appropriate DataType TLV.
// Supported: int32, float32, string, bool.
// Returns an error for unsupported types.
func EncodeAny(value any) (TLV, error) {
	switch v := value.(type) {
	case int:
		return EncodeInteger32(int32(v)) // hmm
	case int32:
		return EncodeInteger32(v)
	case int64:
		return EncodeInteger64(v)
	case float32:
		return EncodeFloat32(v)
	case float64:
		return EncodeFloat64(v)
	case string:
		return EncodeString(v)
	case bool:
		return EncodeBoolean(v)
	default:
		return nil, fmt.Errorf("unsupported value type: %T", v)
	}
}

// DecodeAny decodes a DataType TLV into its Go value.
// Returns an error for unsupported types.
func DecodeAny(tlv TLV) (any, error) {
	if tlv == nil {
		return nil, fmt.Errorf("TLV is nil")
	}
	switch constants.DataType(tlv.Type()) {
	case constants.DATA_TYPE_INTEGER:
		switch tlv.Length() {
		case 4:
			return DecodeInteger32(tlv)
		case 8:
			return DecodeInteger64(tlv)
		default:
			return nil, fmt.Errorf("invalid integer length: %d", tlv.Length())
		}
	case constants.DATA_TYPE_FLOAT:
		switch tlv.Length() {
		case 4:
			return DecodeFloat32(tlv)
		case 8:
			return DecodeFloat64(tlv)
		default:
			return nil, fmt.Errorf("invalid float length: %d", tlv.Length())
		}
	case constants.DATA_TYPE_STRING:
		return DecodeString(tlv)
	case constants.DATA_TYPE_BOOLEAN:
		return DecodeBoolean(tlv)
	default:
		return nil, fmt.Errorf("unsupported TLV type: %x", tlv.Type())
	}
}

// EncodeByte encodes a single byte value.
func EncodeByte(value uint8) []byte {
	return []byte{value}
}

// EncodeByteList encodes a list of bytes.
func EncodeByteList(values []uint8) []byte {
	return values
}

// DecodeByteList decodes a list of bytes.
func DecodeByteList(data []byte) []uint8 {
	return data
}

// DecodeByte decodes a single byte value.
func DecodeByte(data []byte) (uint8, error) {
	if len(data) != 1 {
		return 0, fmt.Errorf("invalid byte length: %d", len(data))
	}
	return data[0], nil
}
