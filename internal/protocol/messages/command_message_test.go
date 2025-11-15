package messages_test

import (
	"testing"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/constants"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/encoding"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/messages"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/selectors"
	"github.com/stretchr/testify/assert"
)

// -------------------------------------- Positive tests --------------------------------------

// ---------------- Command Message Without Node Selector ----------------
func TestNewCommandMessage(t *testing.T) {
	assert := assert.New(t)

	msg := messages.NewCommandMessage(*selectors.NewSingleActuatorSelector(0x05), true)
	assert.NotNil(msg)
	assert.Equal(constants.COMMAND, msg.Type())
	assert.Nil(msg.NodeSelector)
	assert.NotNil(msg.ActuatorSelector)
	assert.Equal(true, msg.ActuatorState)
}

func TestEncode_CommandMessage(t *testing.T) {
	assert := assert.New(t)

	msg := messages.NewCommandMessage(*selectors.NewSingleActuatorSelector(0x05), true)
	encoded, err := msg.Encode()
	assert.NoError(err)
	assert.NotNil(encoded)

	assert.Equal(uint8(constants.COMMAND), encoded.Type())
	// todo test rest
}

func TestDecode_CommandMessage(t *testing.T) {
	assert := assert.New(t)

	aSel, _ := selectors.NewSingleActuatorSelector(0x05).Encode()
	b, _ := encoding.EncodeBoolean(true)
	aState, _ := encoding.NewTLV(uint8(constants.ACTUATOR_STATE), b.Encode())
	encodedTLV, _ := encoding.NewTLV(uint8(constants.COMMAND), encoding.EncodeMultipleTLVs([]encoding.TLV{aSel, aState}))

	decoded, err := messages.DecodeCommandMessage(encodedTLV, false) // expectNode = false
	assert.NoError(err)
	assert.Equal(constants.COMMAND, decoded.Type())
	assert.Nil(decoded.NodeSelector)
	assert.NotNil(decoded.ActuatorSelector)
	assert.Equal([]byte{0x05}, decoded.ActuatorSelector.ActuatorIDs)
	assert.NotNil(decoded.ActuatorState)
	assert.Equal(true, decoded.ActuatorState)
}

// ---------------- Command Message With Node Selector ----------------
func TestNewCommandMessageWithNode(t *testing.T) {
	assert := assert.New(t)

	msg := messages.NewCommandMessageWithNode(
		*selectors.NewSingleNodeSelector(0x03),
		*selectors.NewSingleActuatorSelector(0x05), false)
	assert.NotNil(msg)
	assert.Equal(constants.COMMAND, msg.Type())
	assert.NotNil(msg.NodeSelector)
	assert.NotNil(msg.ActuatorSelector)
	assert.Equal(false, msg.ActuatorState)
}

func TestEncode_CommandMessageWithNode(t *testing.T) {
	assert := assert.New(t)

	msg := messages.NewCommandMessageWithNode(
		*selectors.NewSingleNodeSelector(0x03),
		*selectors.NewSingleActuatorSelector(0x05), false)
	encoded, err := msg.Encode()
	assert.NoError(err)
	assert.NotNil(encoded)

	assert.Equal(uint8(constants.COMMAND), encoded.Type())
	// todo test rest
}

func TestDecode_CommandMessageWithNode(t *testing.T) {
	assert := assert.New(t)

	nSel, _ := selectors.NewSingleNodeSelector(0x03).Encode()
	aSel, _ := selectors.NewSingleActuatorSelector(0x05).Encode()
	b, _ := encoding.EncodeBoolean(false)
	aState, _ := encoding.NewTLV(uint8(constants.ACTUATOR_STATE), b.Encode())
	encodedTLV, _ := encoding.NewTLV(uint8(constants.COMMAND), encoding.EncodeMultipleTLVs([]encoding.TLV{nSel, aSel, aState}))

	decoded, err := messages.DecodeCommandMessage(encodedTLV, true) // expectNode = true
	assert.NoError(err)
	assert.Equal(constants.COMMAND, decoded.Type())
	assert.NotNil(decoded.NodeSelector)
	assert.Equal([]byte{0x03}, decoded.NodeSelector.NodeIDs)
	assert.NotNil(decoded.ActuatorState)
	assert.Equal([]byte{0x05}, decoded.ActuatorSelector.ActuatorIDs)
	assert.Equal(false, decoded.ActuatorState)
}

// -------------------------------------- Negative tests --------------------------------------

func TestDecodeCommandMessage_InvalidArgument(t *testing.T) {
	assert := assert.New(t)

	// Decode with nil TLV
	_, err := messages.DecodeCommandMessage(nil, false)
	assert.Error(err)

	// Decode with wrong TLV type
	wrongType, _ := encoding.NewTLV(uint8(constants.DISCOVERY), []byte{})
	_, err = messages.DecodeCommandMessage(wrongType, false)
	assert.Error(err)

	// Decode with wrong nested TLV type
	wrongNestedTLV, _ := encoding.NewTLV(uint8(constants.ACTUATOR_STATE), []byte{0xFF})
	tl, _ := encoding.NewTLV(uint8(constants.COMMAND), encoding.EncodeMultipleTLVs([]encoding.TLV{wrongNestedTLV}))
	_, err = messages.DecodeCommandMessage(tl, false)
	assert.Error(err)

	// Missing Node Selector
	aSel, _ := selectors.NewSingleActuatorSelector(0x05).Encode()
	b, _ := encoding.EncodeBoolean(true)
	aState, _ := encoding.NewTLV(uint8(constants.ACTUATOR_STATE), b.Encode())
	encodedTLV, _ := encoding.NewTLV(uint8(constants.COMMAND), encoding.EncodeMultipleTLVs([]encoding.TLV{aSel, aState}))
	_, err = messages.DecodeCommandMessage(encodedTLV, true) // expectNode = true
	assert.Error(err)

	// Missing Actuator Selector
	nSel, _ := selectors.NewSingleNodeSelector(0x03).Encode()
	encodedTLV, _ = encoding.NewTLV(uint8(constants.COMMAND), encoding.EncodeMultipleTLVs([]encoding.TLV{nSel, aState}))
	_, err = messages.DecodeCommandMessage(encodedTLV, false)
	assert.Error(err)

	// Missing Actuator State
	encodedTLV, _ = encoding.NewTLV(uint8(constants.COMMAND), encoding.EncodeMultipleTLVs([]encoding.TLV{nSel, aSel}))
	_, err = messages.DecodeCommandMessage(encodedTLV, false)
	assert.Error(err)
}

func TestDecodeCommandMessage_MissingNodeSelector(t *testing.T) {
	assert := assert.New(t)

	aSel, _ := selectors.NewSingleActuatorSelector(0x05).Encode()
	b, _ := encoding.EncodeBoolean(true)
	aState, _ := encoding.NewTLV(uint8(constants.ACTUATOR_STATE), b.Encode())
	encodedTLV, _ := encoding.NewTLV(uint8(constants.COMMAND), encoding.EncodeMultipleTLVs([]encoding.TLV{aSel, aState}))

	decoded, err := messages.DecodeCommandMessage(encodedTLV, false)
	assert.NoError(err)
	assert.NotNil(decoded)
	assert.Nil(decoded.NodeSelector)
	assert.NotNil(decoded.ActuatorState)
	assert.NotNil(decoded.ActuatorSelector)
}
