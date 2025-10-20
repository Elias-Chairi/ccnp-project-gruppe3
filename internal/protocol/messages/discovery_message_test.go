package messages_test

import (
	"testing"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/constants"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/messages"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/tlv"
	"github.com/stretchr/testify/assert"
)

// -------------------------------------- Positive tests --------------------------------------
func TestNewDiscoveryMessage(t *testing.T) {
	assert := assert.New(t)

	msg := messages.NewDiscoveryMessage()
	assert.NotNil(msg)
	assert.Equal(constants.DISCOVERY, msg.Type())
}

func TestEncode_DiscoveryMessage(t *testing.T) {
	assert := assert.New(t)

	msg := messages.NewDiscoveryMessage()
	encoded, err := msg.Encode()
	assert.NoError(err)
	assert.NotNil(encoded)

	expectedTLV, _ := tlv.NewTLV(uint8(constants.DISCOVERY), []byte{})
	expectedData := expectedTLV.Encode()

	assert.Equal(expectedData, encoded)
}

func TestDecode_DiscoveryMessage(t *testing.T) {
	assert := assert.New(t)

	discoveryTLV, _ := tlv.NewTLV(uint8(constants.DISCOVERY), []byte{})

	decoded, err := messages.DecodeDiscoveryMessage(discoveryTLV)
	assert.NoError(err)
	assert.Equal(constants.DISCOVERY, decoded.Type())
}

// -------------------------------------- Negative tests --------------------------------------

func TestDecodeDiscoveryMessage_InvalidArgument(t *testing.T) {
	assert := assert.New(t)

	// nil TLV
	_, err := messages.DecodeDiscoveryMessage(nil)
	assert.Error(err)

	// wrong type
	wrongType, _ := tlv.NewTLV(uint8(constants.COMMAND), []byte{})
	_, err = messages.DecodeDiscoveryMessage(wrongType)
	assert.Error(err)

	// non-empty value
	nonEmptyTLV, _ := tlv.NewTLV(uint8(constants.DISCOVERY), []byte{0x01, 0x02})
	_, err = messages.DecodeDiscoveryMessage(nonEmptyTLV)
	assert.Error(err)
}
