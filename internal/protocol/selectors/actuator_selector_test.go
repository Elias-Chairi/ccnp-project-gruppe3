package selectors_test

import (
	"testing"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/constants"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/selectors"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/tlv"
	"github.com/stretchr/testify/assert"
)

// -------------------------------------- Positive tests --------------------------------------

func TestNewSingleActuatorSelector(t *testing.T) {
	assert := assert.New(t)

	selector := selectors.NewSingleActuatorSelector(0x42)
	assert.NotNil(selector)
	assert.Equal(constants.SINGLE_ACTUATOR, selector.Type)
	assert.Equal([]uint8{0x42}, selector.ActuatorIDs)
	assert.Nil(selector.TypeName)

	// test encoding
	encoded, err := selector.Encode()
	assert.NoError(err)
	assert.NotNil(encoded)
	assert.Equal(uint8(constants.SINGLE_ACTUATOR), encoded.Type())
	assert.Equal(uint16(1), encoded.Length())
	assert.Equal([]byte{0x42}, encoded.Value())

	// test decoding
	selectorDecoded, err := selectors.DecodeActuatorSelector(encoded)
	assert.NoError(err)
	assert.Equal(constants.SINGLE_ACTUATOR, selectorDecoded.Type)
	assert.Equal([]uint8{0x42}, selectorDecoded.ActuatorIDs)
	assert.Nil(selectorDecoded.TypeName)
}

func TestNewActuatorListSelector(t *testing.T) {
	assert := assert.New(t)

	selector := selectors.NewActuatorListSelector([]uint8{0x01, 0x02, 0x03})
	assert.NotNil(selector)
	assert.Equal(constants.ACTUATOR_LIST, selector.Type)
	assert.Equal([]uint8{0x01, 0x02, 0x03}, selector.ActuatorIDs)
	assert.Nil(selector.TypeName)

	// test encoding
	encodedTLV, err := selector.Encode()
	assert.NoError(err)
	assert.NotNil(encodedTLV)
	assert.Equal(uint8(constants.ACTUATOR_LIST), encodedTLV.Type())
	assert.Equal(uint16(3), encodedTLV.Length())
	assert.Equal([]uint8{0x01, 0x02, 0x03}, encodedTLV.Value())

	// test decoding
	selectorDecoded, err := selectors.DecodeActuatorSelector(encodedTLV)
	assert.NoError(err)
	assert.Equal(constants.ACTUATOR_LIST, selectorDecoded.Type)
	assert.Equal([]uint8{0x01, 0x02, 0x03}, selectorDecoded.ActuatorIDs)
	assert.Nil(selectorDecoded.TypeName)
}

func TestNewActuatorTypeSelector(t *testing.T) {
	assert := assert.New(t)

	selector := selectors.NewActuatorTypeSelector("temperature")
	assert.NotNil(selector)
	assert.Equal(constants.ACTUATOR_TYPE, selector.Type)
	assert.NotNil(selector.TypeName)
	assert.Equal("temperature", *selector.TypeName)

	// test encoding
	encodedTLV, err := selector.Encode()
	assert.NoError(err)
	assert.NotNil(encodedTLV)
	assert.Equal(uint8(constants.ACTUATOR_TYPE), encodedTLV.Type())
	assert.Equal(uint16(len("temperature")), encodedTLV.Length())
	assert.Equal([]byte("temperature"), encodedTLV.Value())

	// test decoding
	selectorDecoded, err := selectors.DecodeActuatorSelector(encodedTLV)
	assert.NoError(err)
	assert.Equal(constants.ACTUATOR_TYPE, selectorDecoded.Type)
	assert.NotNil(selectorDecoded.TypeName)
	assert.Equal("temperature", *selectorDecoded.TypeName)
}

func TestNewAllActuatorsSelector(t *testing.T) {
	assert := assert.New(t)

	selector := selectors.NewAllActuatorsSelector()
	assert.NotNil(selector)
	assert.Equal(constants.ALL_ACTUATORS, selector.Type)
	assert.Nil(selector.ActuatorIDs)
	assert.Nil(selector.TypeName)

	// test encoding
	encodedTLV, err := selector.Encode()
	assert.NoError(err)
	assert.NotNil(encodedTLV)
	assert.Equal(uint8(constants.ALL_ACTUATORS), encodedTLV.Type())
	assert.Equal(uint16(0), encodedTLV.Length())
	assert.Equal([]byte{}, encodedTLV.Value())

	// test decoding
	selectorDecoded, err := selectors.DecodeActuatorSelector(encodedTLV)
	assert.NoError(err)
	assert.Equal(constants.ALL_ACTUATORS, selectorDecoded.Type)
	assert.Nil(selectorDecoded.ActuatorIDs)
	assert.Nil(selectorDecoded.TypeName)
}

// -------------------------------------- Negative tests --------------------------------------

func TestSingleActuatorSelector_InvalidCases(t *testing.T) {
	assert := assert.New(t)

	// Empty ActuatorIDs
	selector := &selectors.ActuatorSelector{Type: constants.SINGLE_ACTUATOR, ActuatorIDs: []uint8{}}
	_, err := selector.Encode()
	assert.Error(err)

	// Multiple ActuatorIDs
	selector = &selectors.ActuatorSelector{Type: constants.SINGLE_ACTUATOR, ActuatorIDs: []uint8{0x01, 0x02}}
	_, err = selector.Encode()
	assert.Error(err)

	// Decode with length=0
	tlvData, _ := tlv.DecodeTLV([]byte{uint8(constants.SINGLE_ACTUATOR), 0x00, 0x00})
	_, err = selectors.DecodeActuatorSelector(tlvData)
	assert.Error(err)

	// Decode with length=2
	tlvData, _ = tlv.DecodeTLV([]byte{uint8(constants.SINGLE_ACTUATOR), 0x00, 0x02, 0x01, 0x02})
	_, err = selectors.DecodeActuatorSelector(tlvData)
	assert.Error(err)
}

func TestActuatorListSelector_InvalidCases(t *testing.T) {
	assert := assert.New(t)

	// Empty list
	selector := &selectors.ActuatorSelector{Type: constants.ACTUATOR_LIST, ActuatorIDs: []uint8{}}
	_, err := selector.Encode()
	assert.Error(err)

	// Decode with length=0
	tlvData, _ := tlv.DecodeTLV([]byte{uint8(constants.ACTUATOR_LIST), 0x00, 0x00})
	_, err = selectors.DecodeActuatorSelector(tlvData)
	assert.Error(err)
}

func TestActuatorTypeSelector_InvalidCases(t *testing.T) {
	assert := assert.New(t)

	// Nil TypeName
	selector := &selectors.ActuatorSelector{Type: constants.ACTUATOR_TYPE, TypeName: nil}
	_, err := selector.Encode()
	assert.Error(err)

	// Empty TypeName
	emptyName := ""
	selector = &selectors.ActuatorSelector{Type: constants.ACTUATOR_TYPE, TypeName: &emptyName}
	_, err = selector.Encode()
	assert.Error(err)

	// Whitespace-only TypeName
	whitespaceName := "   "
	selector = &selectors.ActuatorSelector{Type: constants.ACTUATOR_TYPE, TypeName: &whitespaceName}
	_, err = selector.Encode()
	assert.Error(err)

	// Decode with length=0
	tlvData, _ := tlv.DecodeTLV([]byte{uint8(constants.ACTUATOR_TYPE), 0x00, 0x00})
	_, err = selectors.DecodeActuatorSelector(tlvData)
	assert.Error(err)
}

func TestAllActuatorsSelector_InvalidCases(t *testing.T) {
	assert := assert.New(t)

	// Decode with non-empty value
	tlvData, _ := tlv.DecodeTLV([]byte{uint8(constants.ALL_ACTUATORS), 0x00, 0x01, 0x42})
	_, err := selectors.DecodeActuatorSelector(tlvData)
	assert.Error(err)
}

func TestActuatorSelector_UnknownType(t *testing.T) {
	assert := assert.New(t)

	// Encode with unknown type
	selector := &selectors.ActuatorSelector{Type: constants.ActuatorSelector(0xFF), ActuatorIDs: []uint8{0x01}}
	_, err := selector.Encode()
	assert.Error(err)

	// Decode with unknown type
	tlvData, _ := tlv.DecodeTLV([]byte{0xFF, 0x00, 0x01, 0x42})
	_, err = selectors.DecodeActuatorSelector(tlvData)
	assert.Error(err)
}
