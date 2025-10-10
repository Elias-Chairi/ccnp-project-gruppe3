package encoding

import (
	"encoding/binary"
	"fmt"
	"math"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/constants"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/tlv"
)

// EncodeInteger encodes a 32-bit integer into a TLV of type DATA_TYPE_INTEGER
func EncodeInteger(value int32) (tlv.TLV, error) {
	data := make([]byte, 4)
	binary.BigEndian.PutUint32(data, uint32(value))
	return tlv.NewTLV(uint8(constants.DATA_TYPE_INTEGER), data)
}

// DecodeInteger decodes a TLV of type DATA_TYPE_INTEGER into a 32-bit integer.
func DecodeInteger(tlv tlv.TLV) (int32, error) {
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

// EncodeFloat encodes a 32-bit float into a TLV of type DATA_TYPE_FLOAT
func EncodeFloat(value float32) (tlv.TLV, error) {
	data := make([]byte, 4)
	binary.BigEndian.PutUint32(data, math.Float32bits(value))
	return tlv.NewTLV(uint8(constants.DATA_TYPE_FLOAT), data)
}

// DecodeFloat decodes a TLV of type DATA_TYPE_FLOAT into a 32-bit float.
func DecodeFloat(tlv tlv.TLV) (float32, error) {
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

// EncodeString encodes a string into a TLV of type DATA_TYPE_STRING
func EncodeString(value string) (tlv.TLV, error) {
	return tlv.NewTLV(uint8(constants.DATA_TYPE_STRING), []byte(value))
}

// DecodeString decodes a TLV of type DATA_TYPE_STRING into a string.
func DecodeString(tlv tlv.TLV) (string, error) {
	if tlv == nil {
		return "", fmt.Errorf("TLV is nil")
	}
	if tlv.Type() != uint8(constants.DATA_TYPE_STRING) {
		return "", fmt.Errorf("expected string type, got %x", tlv.Type())
	}
	return string(tlv.Value()), nil
}

// EncodeBoolean encodes a boolean into a TLV of type DATA_TYPE_BOOLEAN
func EncodeBoolean(value bool) (tlv.TLV, error) {
	var byteValue byte
	if value {
		byteValue = 0x01
	} else {
		byteValue = 0x00
	}
	return tlv.NewTLV(uint8(constants.DATA_TYPE_BOOLEAN), []byte{byteValue})
}

// DecodeBoolean decodes a TLV of type DATA_TYPE_BOOLEAN into a boolean (0x00=false, 0x01=true).
func DecodeBoolean(tlv tlv.TLV) (bool, error) {
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
func EncodeAny(value any) (tlv.TLV, error) {
	switch v := value.(type) {
	case int32:
		return EncodeInteger(v)
	case float32:
		return EncodeFloat(v)
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
func DecodeAny(tlv tlv.TLV) (any, error) {
	if tlv == nil {
		return nil, fmt.Errorf("TLV is nil")
	}
	switch constants.DataType(tlv.Type()) {
	case constants.DATA_TYPE_INTEGER:
		return DecodeInteger(tlv)
	case constants.DATA_TYPE_FLOAT:
		return DecodeFloat(tlv)
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

// EncodeUint16 encodes a uint16 value in big-endian format.
func EncodeUint16(value uint16) []byte {
	data := make([]byte, 2)
	binary.BigEndian.PutUint16(data, value)
	return data
}

// DecodeUint16 decodes a uint16 value from bytes in big-endian format.
func DecodeUint16(data []byte) (uint16, error) {
	if data == nil {
		return 0, fmt.Errorf("data is nil")
	}
	if len(data) < 2 {
		return 0, fmt.Errorf("insufficient data for uint16")
	}
	return binary.BigEndian.Uint16(data), nil
}

// EncodeByteList encodes a list of bytes.
func EncodeByteList(values []uint8) []byte {
	return values
}

// DecodeByteList decodes a list of bytes.
func DecodeByteList(data []byte) []uint8 {
	return data
}
