package messages_test

import (
	"testing"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/entity"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/constants"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/encoding"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/messages"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/tlv"
	"github.com/stretchr/testify/assert"
)

// -------------------------------------- Positive tests --------------------------------------

func TestNewRegisterNodeMessage(t *testing.T) {
	assert := assert.New(t)

	msg, err := messages.NewRegisterNodeMessage([]entity.Sensor[any]{
		{ID: 0x01, Type: "TEMPERATURE", Unit: "°C", Value: float32(22.5)},
		{ID: 0x02, Type: "HUMIDITY", Unit: "%", Value: int32(60)},
	}, []entity.Actuator[any]{
		{ID: 0x01, Type: "FAN", Unit: "RPM", State: uint8(0x01)},
		{ID: 0x02, Type: "HEATER", Unit: "°C", State: uint8(0x00)},
	})
	assert.NoError(err)
	assert.NotNil(msg)
	assert.Equal(constants.REGISTER_NODE, msg.Type())
	assert.Len(msg.Sensors, 2)
	assert.Len(msg.Actuators, 2)
}

func TestEncode_RegisterNodeMessage(t *testing.T) {
	assert := assert.New(t)

	sensors := []entity.Sensor[any]{{ID: 0x01, Type: "TEMPERATURE", Unit: "°C", Value: float32(22.5)}}
	actuators := []entity.Actuator[any]{{ID: 0x01, Type: "FAN", Unit: "RPM", State: false}}
	msg, _ := messages.NewRegisterNodeMessage(sensors, actuators)
	encoded, err := msg.Encode()
	assert.NoError(err)
	assert.NotNil(encoded)

	sensorTLV, _ := encoding.EncodeSensorEntry(sensors[0])
	actuatorTLV, _ := encoding.EncodeActuatorEntry(actuators[0])
	expectedTLV, _ := tlv.NewTLV(uint8(constants.REGISTER_NODE), tlv.EncodeMultipleTLVs([]tlv.TLV{sensorTLV, actuatorTLV}))
	expectedData := expectedTLV.Encode()

	assert.Equal(expectedData, encoded)
}

func TestDecode_RegisterNodeMessage(t *testing.T) {
	assert := assert.New(t)

	sensors := []entity.Sensor[any]{{ID: 0x01, Type: "TEMPERATURE", Unit: "°C", Value: float32(22.5)}}
	actuators := []entity.Actuator[any]{{ID: 0x01, Type: "FAN", Unit: "RPM", State: false}}
	sensorTLV, _ := encoding.EncodeSensorEntry(sensors[0])
	actuatorTLV, _ := encoding.EncodeActuatorEntry(actuators[0])
	encodedTLV, _ := tlv.NewTLV(uint8(constants.REGISTER_NODE), tlv.EncodeMultipleTLVs([]tlv.TLV{sensorTLV, actuatorTLV}))
	encodedData := encodedTLV.Encode()

	mainTLV, _ := tlv.DecodeTLV(encodedData)
	decoded, err := messages.DecodeRegisterNodeMessage(mainTLV)
	assert.NoError(err)
	assert.NotNil(decoded)
	assert.Len(decoded.Sensors, 1)
	assert.Len(decoded.Actuators, 1)
	assert.Equal(uint8(0x01), decoded.Sensors[0].ID)
	assert.Equal("TEMPERATURE", decoded.Sensors[0].Type)
	assert.Equal("°C", decoded.Sensors[0].Unit)
	assert.Equal(float32(22.5), decoded.Sensors[0].Value)
	assert.Equal(uint8(0x01), decoded.Actuators[0].ID)
	assert.Equal("FAN", decoded.Actuators[0].Type)
	assert.Equal("RPM", decoded.Actuators[0].Unit)
	assert.Equal(false, decoded.Actuators[0].State)
}

// -------------------------------------- Negative tests --------------------------------------

func TestNewRegisterNodeMessage_InvalidArgument(t *testing.T) {
	assert := assert.New(t)

	// Both nil
	_, err := messages.NewRegisterNodeMessage(nil, nil)
	assert.Error(err)

	// Both empty slices
	_, err = messages.NewRegisterNodeMessage([]entity.Sensor[any]{}, []entity.Actuator[any]{})
	assert.Error(err)
}

// ---------------- Decode Validation ----------------
func TestDecodeRegisterNodeMessage_InvalidArgument(t *testing.T) {
	assert := assert.New(t)

	// nil TLV
	_, err := messages.DecodeRegisterNodeMessage(nil)
	assert.Error(err)

	// wrong type
	wrongType, _ := tlv.NewTLV(uint8(constants.COMMAND), []byte{})
	_, err = messages.DecodeRegisterNodeMessage(wrongType)
	assert.Error(err)

	// no sensors or actuators
	emptyTLV, _ := tlv.NewTLV(uint8(constants.REGISTER_NODE), []byte{})
	_, err = messages.DecodeRegisterNodeMessage(emptyTLV)
	assert.Error(err)

	// invalid nested TLV
	invalidTLV, _ := tlv.NewTLV(0xFF, []byte{0x01})
	value := tlv.EncodeMultipleTLVs([]tlv.TLV{invalidTLV})
	registerNodeTLV, _ := tlv.NewTLV(uint8(constants.REGISTER_NODE), value)
	_, err = messages.DecodeRegisterNodeMessage(registerNodeTLV)
	assert.Error(err)
}
