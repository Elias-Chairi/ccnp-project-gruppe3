package encoding_test

import (
	"math"
	"testing"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/constants"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/encoding"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// -------------------------------------- Positive tests --------------------------------------

var integerDecimalBytes = map[int32][]byte{
	42:            {0x00, 0x00, 0x00, 0x2A}, // 0x2A == 42
	-100:          {0xFF, 0xFF, 0xFF, 0x9C}, // 0x9C == -100 in two's complement
	0:             {0x00, 0x00, 0x00, 0x00}, // 0x00 == 0
	math.MaxInt32: {0x7F, 0xFF, 0xFF, 0xFF}, // 0x7F FF FF FF == MaxInt32
	math.MinInt32: {0x80, 0x00, 0x00, 0x00}, // 0x80 00 00 00 == MinInt32
}

func TestEncodeInteger(t *testing.T) {
	for integer, expected := range integerDecimalBytes {
		encoded, err := encoding.EncodeInteger(integer)
		require.NoError(t, err)
		assert.Equal(t, uint8(constants.DATA_TYPE_INTEGER), encoded.Type())
		assert.Equal(t, uint16(4), encoded.Length())
		assert.Equal(t, expected, encoded.Value())
	}
}

func TestDecodeInteger(t *testing.T) {
	for expected, bytes := range integerDecimalBytes {
		encoded, _ := encoding.NewTLV(uint8(constants.DATA_TYPE_INTEGER), bytes)
		tlvData, err := encoding.DecodeInteger(encoded)
		require.NoError(t, err)
		assert.Equal(t, expected, tlvData)
	}
}

var floatDecimalBytes = map[float32][]byte{
	22.5:             {0x41, 0xB4, 0x00, 0x00}, // 0x41B40000 == 22.5 in IEEE 754
	-3.14:            {0xC0, 0x48, 0xF5, 0xC3}, // 0xC048F5C3 == -3.14 in IEEE 754
	0.0:              {0x00, 0x00, 0x00, 0x00}, // 0x00000000 == 0.0 in IEEE 754
	math.MaxFloat32:  {0x7F, 0x7F, 0xFF, 0xFF}, // MaxFloat32 in IEEE 754
	-math.MaxFloat32: {0xFF, 0x7F, 0xFF, 0xFF}, // -MaxFloat32 in IEEE 754
}

func TestEncodeFloat(t *testing.T) {
	for floatVal, expected := range floatDecimalBytes {
		encoded, err := encoding.EncodeFloat(floatVal)
		require.NoError(t, err)
		assert.Equal(t, uint8(constants.DATA_TYPE_FLOAT), encoded.Type())
		assert.Equal(t, uint16(4), encoded.Length())
		assert.Equal(t, expected, encoded.Value())
	}
}

func TestDecodeFloat(t *testing.T) {
	for expected, bytes := range floatDecimalBytes {
		encoded, _ := encoding.NewTLV(uint8(constants.DATA_TYPE_FLOAT), bytes)
		tlvData, err := encoding.DecodeFloat(encoded)
		require.NoError(t, err)
		assert.Equal(t, expected, tlvData)
	}
}

var stringBytes = map[string][]byte{
	"hello": {0x68, 0x65, 0x6C, 0x6C, 0x6F}, // "hello" in UTF-8
	"":      {},                             // empty string
	"°C":    {0xC2, 0xB0, 0x43},             // "°C" in UTF-8
}

func TestEncodeString(t *testing.T) {
	for strVal, expected := range stringBytes {
		encoded, err := encoding.EncodeString(strVal)
		require.NoError(t, err)
		assert.Equal(t, uint8(constants.DATA_TYPE_STRING), encoded.Type())
		assert.Equal(t, uint16(len(expected)), encoded.Length())
		assert.Equal(t, expected, encoded.Value())
	}
}

func TestDecodeString(t *testing.T) {
	for expected, bytes := range stringBytes {
		encoded, _ := encoding.NewTLV(uint8(constants.DATA_TYPE_STRING), bytes)
		tlvData, err := encoding.DecodeString(encoded)
		require.NoError(t, err)
		assert.Equal(t, expected, tlvData)
	}
}

var booleanBytes = map[bool][]byte{
	true:  {0x01},
	false: {0x00},
}

func TestEncodeBoolean(t *testing.T) {
	for boolVal, expected := range booleanBytes {
		encoded, err := encoding.EncodeBoolean(boolVal)
		require.NoError(t, err)
		assert.Equal(t, uint8(constants.DATA_TYPE_BOOLEAN), encoded.Type())
		assert.Equal(t, uint16(1), encoded.Length())
		assert.Equal(t, expected, encoded.Value())
	}
}

func TestDecodeBoolean(t *testing.T) {
	for expected, bytes := range booleanBytes {
		encoded, _ := encoding.NewTLV(uint8(constants.DATA_TYPE_BOOLEAN), bytes)
		tlvData, err := encoding.DecodeBoolean(encoded)
		require.NoError(t, err)
		assert.Equal(t, expected, tlvData)
	}
}

var anyBytes = []struct {
	value    any
	dataType uint8
	bytes    []byte
}{
	{int32(42), uint8(constants.DATA_TYPE_INTEGER), []byte{0x00, 0x00, 0x00, 0x2A}},
	{float32(22.5), uint8(constants.DATA_TYPE_FLOAT), []byte{0x41, 0xB4, 0x00, 0x00}},
	{"hello", uint8(constants.DATA_TYPE_STRING), []byte{0x68, 0x65, 0x6C, 0x6C, 0x6F}},
	{true, uint8(constants.DATA_TYPE_BOOLEAN), []byte{0x01}},
	{false, uint8(constants.DATA_TYPE_BOOLEAN), []byte{0x00}},
}

func TestEncodeAny(t *testing.T) {
	for _, val := range anyBytes {
		encoded, err := encoding.EncodeAny(val.value)
		require.NoError(t, err)
		assert.Equal(t, val.dataType, encoded.Type())
		assert.Equal(t, uint16(len(val.bytes)), encoded.Length())
		assert.Equal(t, val.bytes, encoded.Value())
	}
}

func TestDecodeAny(t *testing.T) {
	for _, val := range anyBytes {
		encoded, _ := encoding.NewTLV(val.dataType, val.bytes)
		decoded, err := encoding.DecodeAny(encoded)
		require.NoError(t, err)
		assert.Equal(t, val.value, decoded)
	}
}

func TestEncodeByte(t *testing.T) {
	assert.Equal(t, []byte{0x42}, encoding.EncodeByte(0x42))
}

func TestEncodeByteList(t *testing.T) {
	assert.Equal(t, []byte{0x01, 0x02, 0x03}, encoding.EncodeByteList([]uint8{0x01, 0x02, 0x03}))
}

func TestDecodeByteList(t *testing.T) {
	assert.Equal(t, []uint8{0x01, 0x02, 0x03}, encoding.DecodeByteList([]byte{0x01, 0x02, 0x03}))
}

// -------------------------------------- Negative tests --------------------------------------

func TestDecodeInteger_InvalidArgument(t *testing.T) {
	// nil TLV
	_, err := encoding.DecodeInteger(nil)
	assert.Error(t, err)

	// not an integer TLV
	wrongType, _ := encoding.NewTLV(byte(constants.DATA_TYPE_FLOAT), []byte{0, 0, 0, 0})
	_, err = encoding.DecodeInteger(wrongType)
	assert.Error(t, err)

	// wrong length
	wrongLength, _ := encoding.NewTLV(byte(constants.DATA_TYPE_INTEGER), []byte{0, 0})
	_, err = encoding.DecodeInteger(wrongLength)
	assert.Error(t, err)
}

func TestDecodeFloat_InvalidArgument(t *testing.T) {
	// nil TLV
	_, err := encoding.DecodeFloat(nil)
	assert.Error(t, err)

	// not a float TLV
	wrongType, _ := encoding.NewTLV(byte(constants.DATA_TYPE_INTEGER), []byte{0, 0, 0, 0})
	_, err = encoding.DecodeFloat(wrongType)
	assert.Error(t, err)

	// wrong length
	wrongLength, _ := encoding.NewTLV(byte(constants.DATA_TYPE_FLOAT), []byte{0, 0})
	_, err = encoding.DecodeFloat(wrongLength)
	assert.Error(t, err)
}

func TestDecodeString_InvalidArgument(t *testing.T) {
	// nil TLV
	_, err := encoding.DecodeString(nil)
	assert.Error(t, err)

	// not a string TLV
	wrongType, _ := encoding.NewTLV(byte(constants.DATA_TYPE_INTEGER), []byte{0})
	_, err = encoding.DecodeString(wrongType)
	assert.Error(t, err)
}

func TestDecodeBoolean_InvalidArgument(t *testing.T) {
	// nil TLV
	_, err := encoding.DecodeBoolean(nil)
	assert.Error(t, err)

	// not a boolean TLV
	wrongType, _ := encoding.NewTLV(byte(constants.DATA_TYPE_INTEGER), []byte{0x01})
	_, err = encoding.DecodeBoolean(wrongType)
	assert.Error(t, err)

	// wrong length
	wrongLength, _ := encoding.NewTLV(byte(constants.DATA_TYPE_BOOLEAN), []byte{0, 0})
	_, err = encoding.DecodeBoolean(wrongLength)
	assert.Error(t, err)

	// invalid value
	invalidValue, _ := encoding.NewTLV(byte(constants.DATA_TYPE_BOOLEAN), []byte{0x02})
	_, err = encoding.DecodeBoolean(invalidValue)
	assert.Error(t, err)
}

func TestEncodeAny_InvalidArgument(t *testing.T) {
	// nil value
	_, err := encoding.EncodeAny(nil)
	assert.Error(t, err)

	// unsupported type
	_, err = encoding.EncodeAny(int(42))
	assert.Error(t, err)

	// slice of unsupported type
	_, err = encoding.EncodeAny([]byte{0x01})
	assert.Error(t, err)
}

func TestDecodeAny_InvalidArgument(t *testing.T) {
	// nil TLV
	_, err := encoding.DecodeAny(nil)
	assert.Error(t, err)

	// unsupported TLV type
	unknownType, _ := encoding.NewTLV(0xFF, []byte{0x01})
	_, err = encoding.DecodeAny(unknownType)
	assert.Error(t, err)
}
