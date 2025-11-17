package messages_test

import (
	"testing"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/constants"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/encoding"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/messages"
	"github.com/stretchr/testify/assert"
)

// -------------------------------------- Positive tests --------------------------------------

func TestAckSuccessMessage(t *testing.T) {
	assert := assert.New(t)

	msg := messages.AckSuccessMessage()
	assert.False(msg.IsError())
	assert.Equal(constants.ACK_SUCCESS, msg.Code)
	assert.Equal("", msg.Data)
}

func TestEncode(t *testing.T) {
	assert := assert.New(t)

	msg := messages.AckErrorRequestIDMessage{
		Code: constants.ACK_SUCCESS,
		Data: "success",
	}
	encoded, err := msg.Encode()
	assert.NoError(err)
	assert.NotNil(encoded)

	assert.Equal(uint8(constants.ACK_ERROR), encoded.Type())       // Type
	assert.Equal(uint16(1+len("success")), encoded.Length())       // length
	assert.Equal(uint8(constants.ACK_SUCCESS), encoded.Value()[0]) // First byte of value is the code
	assert.Equal("success", string(encoded.Value()[1:]))           // Remaining bytes are the data
}

func TestDecode(t *testing.T) {
	assert := assert.New(t)

	tlvData, _ := encoding.NewTLV(uint8(constants.ACK_ERROR), append([]byte{uint8(constants.ACK_SUCCESS)}, []byte("All good")...))
	decoded, err := messages.DecodeAckErrorMessage(tlvData)
	assert.NoError(err)
	assert.NotNil(decoded)
	assert.False(decoded.IsError())
	assert.Equal(constants.ACK_SUCCESS, decoded.Code)
	assert.Equal("All good", decoded.Data)
}

// -------------------------------------- Negative tests --------------------------------------

func TestEncode_AckMessage_NoData(t *testing.T) {
	assert := assert.New(t)

	msg := messages.AckErrorRequestIDMessage{
		Code: 0,
		Data: "",
	}
	encoded, err := msg.Encode()
	assert.NoError(err)
	assert.NotNil(encoded)
	assert.Equal(uint8(constants.ACK_ERROR), encoded.Type()) // Type
	assert.Equal(uint16(1), encoded.Length())                // length
	assert.Equal(uint8(0), encoded.Value()[0])               // Code
	assert.Equal(4, len(encoded.Encode()))                   // Total length should be 4 bytes (Type + Length + Code)
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
	assert.True(decoded.IsError())
	assert.Equal(constants.ERR_INVALID_VALUE, decoded.Code)
	assert.Equal("", decoded.Data)
}
