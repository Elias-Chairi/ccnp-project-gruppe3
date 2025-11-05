package messages_test

import (
	"testing"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/constants"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/encoding"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/messages"
	"github.com/stretchr/testify/assert"
)

// -------------------------------------- Positive tests --------------------------------------

func TestNewRegisterControlMessage(t *testing.T) {
	assert := assert.New(t)

	msg := messages.NewRegisterControlMessage()
	assert.NotNil(msg)
	assert.Equal(constants.REGISTER_CONTROL, msg.Type())
}

func TestEncode_RegisterControlMessage(t *testing.T) {
	assert := assert.New(t)

	msg := messages.NewRegisterControlMessage()
	encoded, err := msg.Encode()
	assert.NoError(err)
	assert.NotNil(encoded)

	expectedTLV, _ := encoding.NewTLV(uint8(constants.REGISTER_CONTROL), []byte{})
	expectedData := expectedTLV.Encode()

	assert.Equal(expectedData, encoded)
}

func TestDecode_RegisterControlMessage(t *testing.T) {
	assert := assert.New(t)

	registerControlTLV, _ := encoding.NewTLV(uint8(constants.REGISTER_CONTROL), []byte{})

	decoded, err := messages.DecodeRegisterControlMessage(registerControlTLV)
	assert.NoError(err)
	assert.Equal(constants.REGISTER_CONTROL, decoded.Type())
}

// -------------------------------------- Negative tests --------------------------------------

func TestDecodeRegisterControlMessage_InvalidArgument(t *testing.T) {
	assert := assert.New(t)

	// nil TLV
	_, err := messages.DecodeRegisterControlMessage(nil)
	assert.Error(err)

	// wrong type
	wrongType, _ := encoding.NewTLV(uint8(constants.COMMAND), []byte{})
	_, err = messages.DecodeRegisterControlMessage(wrongType)
	assert.Error(err)

	// non-empty value
	nonEmptyTLV, _ := encoding.NewTLV(uint8(constants.REGISTER_CONTROL), []byte{0x01, 0x02})
	_, err = messages.DecodeRegisterControlMessage(nonEmptyTLV)
	assert.Error(err)
}
