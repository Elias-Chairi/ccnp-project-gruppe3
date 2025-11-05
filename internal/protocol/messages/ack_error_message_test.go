package messages_test

import (
	"encoding/binary"
	"testing"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/constants"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/encoding"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/messages"
	"github.com/stretchr/testify/assert"
)

// -------------------------------------- Positive tests --------------------------------------

// ---------------- ACK Message (Success) ----------------
func TestNewAckMessage(t *testing.T) {
	assert := assert.New(t)

	data := "Operation successful"
	msg := messages.NewAckMessage(&data)
	assert.Equal(constants.ACK_ERROR, msg.Type())
	assert.False(msg.IsError())
}

func TestEncode_AckMessage(t *testing.T) {
	assert := assert.New(t)

	data := "success"
	msg := messages.NewAckMessage(&data)
	encoded, err := msg.Encode()
	assert.NoError(err)
	assert.NotNil(encoded)

	assert.Equal(uint8(constants.ACK_ERROR), encoded[0])                          // Type
	assert.Equal(uint16(1+len("success")), binary.BigEndian.Uint16(encoded[1:3])) // length
	assert.Equal(uint8(constants.ACK_SUCCESS), encoded[3])                        // First byte of value is the code
	assert.Equal("success", string(encoded[4:]))                                  // Remaining bytes are the data
}

func TestDecode_AckMessage(t *testing.T) {
	assert := assert.New(t)

	tlvData, _ := encoding.NewTLV(uint8(constants.ACK_ERROR), append([]byte{uint8(constants.ACK_SUCCESS)}, []byte("All good")...))
	decoded, err := messages.DecodeAckErrorMessage(tlvData)
	assert.NoError(err)
	assert.NotNil(decoded)
	assert.Equal(constants.ACK_ERROR, decoded.Type())
	assert.False(decoded.IsError())
	assert.Equal(constants.ACK_SUCCESS, decoded.Code)
	assert.NotNil(decoded.Data)
	assert.Equal("All good", *decoded.Data)
}

// ---------------- Error Message ----------------
func TestNewErrorMessage(t *testing.T) {
	assert := assert.New(t)

	errorMsg := "Invalid value"
	msg, err := messages.NewErrorMessage(constants.ERR_INVALID_VALUE, &errorMsg)
	assert.NoError(err)
	assert.NotNil(msg)
	assert.Equal(constants.ACK_ERROR, msg.Type())
	assert.True(msg.IsError())
}

func TestEncode_ErrorMessage(t *testing.T) {
	assert := assert.New(t)

	errorMsg := "Invalid value"
	msg, _ := messages.NewErrorMessage(constants.ERR_INVALID_VALUE, &errorMsg)
	encoded, err := msg.Encode()
	assert.NoError(err)
	assert.NotNil(encoded)

	assert.Equal(uint8(constants.ACK_ERROR), encoded[0])                                // Type
	assert.Equal(uint16(1+len("Invalid value")), binary.BigEndian.Uint16(encoded[1:3])) // length
	assert.Equal(uint8(constants.ERR_INVALID_VALUE), encoded[3])                        // First byte of value
	assert.Equal("Invalid value", string(encoded[4:]))                                  // Remaining bytes are the data
}

func TestDecode_ErrorMessage(t *testing.T) {
	assert := assert.New(t)

	tlvData, _ := encoding.NewTLV(uint8(constants.ACK_ERROR), append([]byte{uint8(constants.ERR_INVALID_VALUE)}, []byte("Invalid value")...))
	decoded, err := messages.DecodeAckErrorMessage(tlvData)
	assert.NoError(err)
	assert.NotNil(decoded)
	assert.Equal(constants.ACK_ERROR, decoded.Type())
	assert.True(decoded.IsError())
	assert.Equal(constants.ERR_INVALID_VALUE, decoded.Code)
	assert.NotNil(decoded.Data)
	assert.Equal("Invalid value", *decoded.Data)
}

// -------------------------------------- Negative tests --------------------------------------
func TestNewAckMessage_EmptyData(t *testing.T) {
	assert := assert.New(t)

	emptyData := ""
	msg := messages.NewAckMessage(&emptyData)
	assert.NotNil(msg)
	assert.Equal(constants.ACK_SUCCESS, msg.Code)
	assert.NotNil(msg.Data)
	assert.Equal("", *msg.Data)

	msg = messages.NewAckMessage(nil)
	assert.NotNil(msg)
	assert.Equal(constants.ACK_SUCCESS, msg.Code)
	assert.Nil(msg.Data)
}

func TestNewErrorMessage_EmptyData(t *testing.T) {
	assert := assert.New(t)

	emptyData := ""
	msg, err := messages.NewErrorMessage(constants.ERR_INVALID_VALUE, &emptyData)
	assert.NoError(err)
	assert.NotNil(msg)
	assert.Equal(constants.ERR_INVALID_VALUE, msg.Code)
	assert.NotNil(msg.Data)
	assert.Equal("", *msg.Data)

	msg, err = messages.NewErrorMessage(constants.ERR_INVALID_VALUE, nil)
	assert.NoError(err)
	assert.NotNil(msg)
	assert.Equal(constants.ERR_INVALID_VALUE, msg.Code)
	assert.Nil(msg.Data)
}

func TestEncode_AckMessage_NilData(t *testing.T) {
	assert := assert.New(t)

	msg := messages.NewAckMessage(nil)
	encoded, err := msg.Encode()
	assert.NoError(err)
	assert.NotNil(encoded)
	assert.Equal(uint8(constants.ACK_ERROR), encoded[0]) // Type
	assert.Equal(uint16(1), binary.BigEndian.Uint16(encoded[1:3]))
	assert.Equal(uint8(constants.ACK_SUCCESS), encoded[3]) // Code
	assert.Equal(4, len(encoded))                          // Total length should be 4 bytes (Type + Length + Code)
}

func TestNewErrorMessage_InvalidCode(t *testing.T) {
	// Invalid code (not defined)
	_, err := messages.NewErrorMessage(0xFF, nil)
	assert.Error(t, err)

	// Success code not allowed for error message
	_, err = messages.NewErrorMessage(constants.ACK_SUCCESS, nil)
	assert.Error(t, err)
}

func TestDecodeInvalidState_AckErrorMessage(t *testing.T) {
	assert := assert.New(t)

	// Nil TLV
	_, err := messages.DecodeAckErrorMessage(nil)
	assert.Error(err)

	// Wrong type
	wrongType, _ := encoding.NewTLV(uint8(constants.DISCOVERY), []byte{0xA0})
	_, err = messages.DecodeAckErrorMessage(wrongType)
	assert.Error(err)

	// Empty value (no code)
	emptyValue, _ := encoding.NewTLV(uint8(constants.ACK_ERROR), []byte{})
	_, err = messages.DecodeAckErrorMessage(emptyValue)
	assert.Error(err)

	// Invalid error code
	invalidCode, _ := encoding.NewTLV(uint8(constants.ACK_ERROR), []byte{0xFF, 0x01, 0x02})
	_, err = messages.DecodeAckErrorMessage(invalidCode)
	assert.Error(err)
}

func TestDecodeAckErrorMessage_EmptyData(t *testing.T) {
	assert := assert.New(t)

	validNoData, _ := encoding.NewTLV(uint8(constants.ACK_ERROR), []byte{uint8(constants.ERR_INVALID_VALUE)})
	decoded, err := messages.DecodeAckErrorMessage(validNoData)
	assert.NoError(err)
	assert.Equal(constants.ACK_ERROR, decoded.Type())
	assert.True(decoded.IsError())
	assert.Equal(constants.ERR_INVALID_VALUE, decoded.Code)
	assert.Nil(decoded.Data)
}
