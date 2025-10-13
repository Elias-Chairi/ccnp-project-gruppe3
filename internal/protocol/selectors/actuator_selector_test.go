package selectors_test

import (
	"testing"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/constants"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/selectors"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/tlv"
	"github.com/stretchr/testify/assert"
)

// -------------------------------------- Positive tests --------------------------------------

// ---------------- Single Actuator Selector ----------------
func TestNewSingleActuatorSelector(t *testing.T) {
	assert := assert.New(t)

	selector := selectors.NewSingleActuatorSelector(0x42)
	assert.NotNil(selector)
	assert.Equal(constants.SINGLE_ACTUATOR, selector.Type)
	assert.Equal([]uint8{0x42}, selector.ActuatorIDs)
	assert.Nil(selector.TypeName)
}

func TestEncode_SingleActuator(t *testing.T) {
	assert := assert.New(t)

	selector := selectors.NewSingleActuatorSelector(0x42)
	encoded, err := selector.Encode()
	assert.NoError(err)
	assert.NotNil(encoded)
	assert.Equal(uint8(constants.SINGLE_ACTUATOR), encoded.Type())
	assert.Equal(uint16(1), encoded.Length())
	assert.Equal([]byte{0x42}, encoded.Value())
}

func TestDecode_SingleActuator(t *testing.T) {
	assert := assert.New(t)

	encoded, _ := tlv.NewTLV(uint8(constants.SINGLE_ACTUATOR), []byte{0x42})
	selectorDecoded, err := selectors.DecodeActuatorSelector(encoded)
	assert.NoError(err)
	assert.Equal(constants.SINGLE_ACTUATOR, selectorDecoded.Type)
	assert.Equal([]uint8{0x42}, selectorDecoded.ActuatorIDs)
	assert.Nil(selectorDecoded.TypeName)
}

// ---------------- Actuator List Selector ----------------
func TestNewActuatorListSelector(t *testing.T) {
	assert := assert.New(t)

	selector := selectors.NewActuatorListSelector([]uint8{0x01, 0x02, 0x03})
	assert.NotNil(selector)
	assert.Equal(constants.ACTUATOR_LIST, selector.Type)
	assert.Equal([]uint8{0x01, 0x02, 0x03}, selector.ActuatorIDs)
	assert.Nil(selector.TypeName)
}

func TestEncode_ActuatorList(t *testing.T) {
	assert := assert.New(t)

	selector := selectors.NewActuatorListSelector([]uint8{0x01, 0x02, 0x03})
	encodedTLV, err := selector.Encode()
	assert.NoError(err)
	assert.NotNil(encodedTLV)
	assert.Equal(uint8(constants.ACTUATOR_LIST), encodedTLV.Type())
	assert.Equal(uint16(3), encodedTLV.Length())
	assert.Equal([]uint8{0x01, 0x02, 0x03}, encodedTLV.Value())
}

func TestDecode_ActuatorList(t *testing.T) {
	assert := assert.New(t)

	encodedTLV, _ := tlv.NewTLV(uint8(constants.ACTUATOR_LIST), []byte{0x01, 0x02, 0x03})
	selectorDecoded, err := selectors.DecodeActuatorSelector(encodedTLV)
	assert.NoError(err)
	assert.Equal(constants.ACTUATOR_LIST, selectorDecoded.Type)
	assert.Equal([]uint8{0x01, 0x02, 0x03}, selectorDecoded.ActuatorIDs)
	assert.Nil(selectorDecoded.TypeName)
}

// ---------------- Actuator Type Selector ----------------
func TestNewActuatorTypeSelector(t *testing.T) {
	assert := assert.New(t)

	selector := selectors.NewActuatorTypeSelector("temperature")
	assert.NotNil(selector)
	assert.Equal(constants.ACTUATOR_TYPE, selector.Type)
	assert.NotNil(selector.TypeName)
	assert.Equal("temperature", *selector.TypeName)
}

func TestEncode_ActuatorType(t *testing.T) {
	assert := assert.New(t)

	selector := selectors.NewActuatorTypeSelector("temperature")
	encodedTLV, err := selector.Encode()
	assert.NoError(err)
	assert.NotNil(encodedTLV)
	assert.Equal(uint8(constants.ACTUATOR_TYPE), encodedTLV.Type())
	assert.Equal(uint16(len("temperature")), encodedTLV.Length())
	assert.Equal([]byte("temperature"), encodedTLV.Value())
}

func TestDecode_ActuatorType(t *testing.T) {
	assert := assert.New(t)

	encodedTLV, _ := tlv.NewTLV(uint8(constants.ACTUATOR_TYPE), []byte("temperature"))
	selectorDecoded, err := selectors.DecodeActuatorSelector(encodedTLV)
	assert.NoError(err)
	assert.Equal(constants.ACTUATOR_TYPE, selectorDecoded.Type)
	assert.NotNil(selectorDecoded.TypeName)
	assert.Equal("temperature", *selectorDecoded.TypeName)
}

// ---------------- All Actuators Selector ----------------
func TestNewAllActuatorsSelector(t *testing.T) {
	assert := assert.New(t)

	selector := selectors.NewAllActuatorsSelector()
	assert.NotNil(selector)
	assert.Equal(constants.ALL_ACTUATORS, selector.Type)
	assert.Nil(selector.ActuatorIDs)
	assert.Nil(selector.TypeName)
}

func TestEncode_AllActuators(t *testing.T) {
	assert := assert.New(t)

	selector := selectors.NewAllActuatorsSelector()
	encodedTLV, err := selector.Encode()
	assert.NoError(err)
	assert.NotNil(encodedTLV)
	assert.Equal(uint8(constants.ALL_ACTUATORS), encodedTLV.Type())
	assert.Equal(uint16(0), encodedTLV.Length())
	assert.Equal([]byte{}, encodedTLV.Value())
}

func TestDecode_AllActuators(t *testing.T) {
	assert := assert.New(t)

	encodedTLV, _ := tlv.NewTLV(uint8(constants.ALL_ACTUATORS), []byte{})
	selectorDecoded, err := selectors.DecodeActuatorSelector(encodedTLV)
	assert.NoError(err)
	assert.Equal(constants.ALL_ACTUATORS, selectorDecoded.Type)
	assert.Nil(selectorDecoded.ActuatorIDs)
	assert.Nil(selectorDecoded.TypeName)
}

// -------------------------------------- Negative tests --------------------------------------

// ---------------- Single Actuator Selector ----------------
func TestEncodeInvalidState_SingleActuator(t *testing.T) {
	assert := assert.New(t)

	// Empty ActuatorIDs
	selector := &selectors.ActuatorSelector{Type: constants.SINGLE_ACTUATOR, ActuatorIDs: []uint8{}}
	_, err := selector.Encode()
	assert.Error(err)

	// Multiple ActuatorIDs
	selector = &selectors.ActuatorSelector{Type: constants.SINGLE_ACTUATOR, ActuatorIDs: []uint8{0x01, 0x02}}
	_, err = selector.Encode()
	assert.Error(err)

	// Optional fields filled
	unused := "unused"
	selector = &selectors.ActuatorSelector{Type: constants.SINGLE_ACTUATOR, ActuatorIDs: []uint8{0x01}, TypeName: &unused}
	encoded, err := selector.Encode()
	assert.Nil(err)
	assert.Equal(uint8(constants.SINGLE_ACTUATOR), encoded.Type())
	assert.Equal(1, int(encoded.Length()))
	assert.Equal([]uint8{0x01}, encoded.Value())
}

func TestDecodeInvalidState_SingleActuator(t *testing.T) {
	assert := assert.New(t)

	// Decode with length=0
	tlvData, _ := tlv.NewTLV(uint8(constants.SINGLE_ACTUATOR), nil)
	_, err := selectors.DecodeActuatorSelector(tlvData)
	assert.Error(err)

	// Decode with length=2
	tlvData, _ = tlv.NewTLV(uint8(constants.SINGLE_ACTUATOR), []byte{0x01, 0x02})
	_, err = selectors.DecodeActuatorSelector(tlvData)
	assert.Error(err)
}

// ---------------- Actuator List Selector ----------------
func TestEncodeInvalidState_ActuatorList(t *testing.T) {
	assert := assert.New(t)
	// Nil list
	selector := &selectors.ActuatorSelector{Type: constants.ACTUATOR_LIST, ActuatorIDs: nil}
	_, err := selector.Encode()
	assert.Error(err)

	// Empty list
	selector = &selectors.ActuatorSelector{Type: constants.ACTUATOR_LIST, ActuatorIDs: []uint8{}}
	_, err = selector.Encode()
	assert.Error(err)

	// Optional fields filled
	unused := "unused"
	selector = &selectors.ActuatorSelector{Type: constants.ACTUATOR_LIST, ActuatorIDs: []uint8{0x01, 0x02}, TypeName: &unused}
	encoded, err := selector.Encode()
	assert.Nil(err)
	assert.Equal(uint8(constants.ACTUATOR_LIST), encoded.Type())
	assert.Equal(2, int(encoded.Length()))
	assert.Equal([]byte{0x01, 0x02}, encoded.Value())
}

func TestDecodeInvalidState_ActuatorList(t *testing.T) {
	// Decode with length=0
	tlvData, _ := tlv.NewTLV(uint8(constants.ACTUATOR_LIST), nil)
	_, err := selectors.DecodeActuatorSelector(tlvData)
	assert.Error(t, err)
}

// ---------------- Actuator Type Selector ----------------
func TestEncodeInvalidState_ActuatorType(t *testing.T) {
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

	// Optional fields filled
	typeName := "typeName"
	selector = &selectors.ActuatorSelector{Type: constants.ACTUATOR_TYPE, ActuatorIDs: []uint8{0x01, 0x02}, TypeName: &typeName}
	encoded, err := selector.Encode()
	assert.Nil(err)
	assert.Equal(uint8(constants.ACTUATOR_TYPE), encoded.Type())
	assert.Equal(len("typeName"), int(encoded.Length()))
	assert.Equal([]byte("typeName"), encoded.Value())
}

func TestDecodeInvalidState_ActuatorType(t *testing.T) {
	assert := assert.New(t)

	// Decode with length=0
	tlvData, _ := tlv.NewTLV(uint8(constants.ACTUATOR_TYPE), nil)
	_, err := selectors.DecodeActuatorSelector(tlvData)
	assert.Error(err)

	// Decode with whitespace-only value
	tlvData, _ = tlv.NewTLV(uint8(constants.ACTUATOR_TYPE), []byte("   "))
	_, err = selectors.DecodeActuatorSelector(tlvData)
	assert.Error(err)
}

// ---------------- All Actuators Selector ----------------
func TestEncodeInvalidState_AllActuators(t *testing.T) {
	assert := assert.New(t)

	// Optional fields filled
	unused := "unused"
	selector := &selectors.ActuatorSelector{Type: constants.ALL_ACTUATORS, ActuatorIDs: []uint8{0x01}, TypeName: &unused}
	encoded, err := selector.Encode()
	assert.Nil(err)
	assert.Equal(uint8(constants.ALL_ACTUATORS), encoded.Type())
	assert.Equal(0, int(encoded.Length()))
	assert.Equal([]byte{}, encoded.Value())
}

func TestDecodeInvalidState_AllActuators(t *testing.T) {
	assert := assert.New(t)

	// Decode with non-empty value
	tlvData, _ := tlv.NewTLV(uint8(constants.ALL_ACTUATORS), []byte{0x42})
	_, err := selectors.DecodeActuatorSelector(tlvData)
	assert.Error(err)
}

// ---------------- Unknown Selector Type ----------------

func TestActuatorSelectorEncodeInvalidState_UnknownType(t *testing.T) {
	selector := &selectors.ActuatorSelector{Type: constants.ActuatorSelector(0xFF), ActuatorIDs: []uint8{0x01}}
	_, err := selector.Encode()
	assert.Error(t, err)
}

func TestActuatorSelectorDecodeInvalidState_UnknownType(t *testing.T) {
	tlvData, _ := tlv.NewTLV(0xFF, []byte{0x00, 0x01, 0x42})
	_, err := selectors.DecodeActuatorSelector(tlvData)
	assert.Error(t, err)
}
