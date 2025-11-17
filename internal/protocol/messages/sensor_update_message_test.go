package messages_test

import (
	"testing"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/entity"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/constants"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/encoding"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/messages"
	"github.com/stretchr/testify/assert"
)

// -------------------------------------- Positive tests --------------------------------------

// ---------------- Sensor Update Message Without Node ID (Node → Server) ----------------
func TestNewSensorUpdateMessage(t *testing.T) {
	assert := assert.New(t)

	msg := messages.NewSensorUpdateMessage(entity.Sensor[any]{
		ID:    0x01,
		Type:  "TEMPERATURE",
		Unit:  "°C",
		Value: float32(22.5),
	})
	assert.NotNil(msg)
	assert.Equal(constants.SENSOR_UPDATE, msg.Type())
	assert.Nil(msg.NodeID)
	assert.Equal(uint8(0x01), msg.Sensor.ID)
	assert.Equal("TEMPERATURE", msg.Sensor.Type)
	assert.Equal("°C", msg.Sensor.Unit)
	assert.Equal(float32(22.5), msg.Sensor.Value)
}

func TestEncode_SensorUpdateMessage(t *testing.T) {
	assert := assert.New(t)

	sensor := entity.Sensor[any]{
		ID:    0x01,
		Type:  "TEMPERATURE",
		Unit:  "°C",
		Value: float32(22.5),
	}
	msg := messages.NewSensorUpdateMessage(sensor)
	encoded, err := msg.Encode()
	assert.NoError(err)
	assert.NotNil(encoded)

	sensorEntryTLV, _ := encoding.EncodeSensorEntry(sensor)
	expectedTLV, _ := encoding.NewTLV(uint8(constants.SENSOR_UPDATE), encoding.EncodeMultipleTLVs([]encoding.TLV{sensorEntryTLV}))
	expectedData := expectedTLV.Encode()

	assert.Equal(expectedData, encoded.Encode())
}

func TestDecode_SensorUpdateMessage(t *testing.T) {
	assert := assert.New(t)

	sensorEntryTLV, _ := encoding.EncodeSensorEntry(entity.Sensor[any]{
		ID:    0x01,
		Type:  "TEMPERATURE",
		Unit:  "°C",
		Value: float32(22.5),
	})
	encodedTLV, _ := encoding.NewTLV(uint8(constants.SENSOR_UPDATE), encoding.EncodeMultipleTLVs([]encoding.TLV{sensorEntryTLV}))

	decoded, err := messages.DecodeSensorUpdateMessage(encodedTLV, false)
	assert.NoError(err)
	assert.NotNil(decoded)
	assert.Nil(decoded.NodeID)
	assert.Equal(uint8(0x01), decoded.Sensor.ID)
	assert.Equal("TEMPERATURE", decoded.Sensor.Type)
	assert.Equal("°C", decoded.Sensor.Unit)
	assert.Equal(float32(22.5), decoded.Sensor.Value)
}

// ---------------- Sensor Update Message With Node ID (Server → Control) ----------------
func TestNewSensorUpdateMessageWithNode(t *testing.T) {
	assert := assert.New(t)

	msg := messages.NewSensorUpdateMessageWithNode(0x07, entity.Sensor[any]{
		ID:    0x02,
		Type:  "HUMIDITY",
		Unit:  "%",
		Value: int32(60),
	})
	assert.NotNil(msg)
	assert.Equal(constants.SENSOR_UPDATE, msg.Type())
	assert.NotNil(msg.NodeID)
	assert.Equal(uint8(0x07), *msg.NodeID)
	assert.Equal(uint8(0x02), msg.Sensor.ID)
	assert.Equal("HUMIDITY", msg.Sensor.Type)
	assert.Equal("%", msg.Sensor.Unit)
	assert.Equal(int32(60), msg.Sensor.Value)
}

func TestEncode_SensorUpdateMessageWithNode(t *testing.T) {
	assert := assert.New(t)

	sensor := entity.Sensor[any]{
		ID:    0x02,
		Type:  "HUMIDITY",
		Unit:  "%",
		Value: int32(60),
	}
	msg := messages.NewSensorUpdateMessageWithNode(0x07, sensor)
	encoded, err := msg.Encode()
	assert.NoError(err)
	assert.NotNil(encoded)

	nodeIDTLV, _ := encoding.NewTLV(uint8(constants.SINGLE_NODE), []byte{0x07})
	sensorEntryTLV, _ := encoding.EncodeSensorEntry(sensor)
	expectedTLV, _ := encoding.NewTLV(uint8(constants.SENSOR_UPDATE), encoding.EncodeMultipleTLVs([]encoding.TLV{nodeIDTLV, sensorEntryTLV}))
	expectedData := expectedTLV.Encode()

	assert.Equal(expectedData, encoded.Encode())
}

func TestDecode_SensorUpdateMessageWithNode(t *testing.T) {
	assert := assert.New(t)

	nodeIDTLV, _ := encoding.NewTLV(uint8(constants.SINGLE_NODE), []byte{0x07})
	sensorEntryTLV, _ := encoding.EncodeSensorEntry(entity.Sensor[any]{
		ID:    0x02,
		Type:  "HUMIDITY",
		Unit:  "%",
		Value: int32(60),
	})
	encodedTLV, _ := encoding.NewTLV(uint8(constants.SENSOR_UPDATE), encoding.EncodeMultipleTLVs([]encoding.TLV{nodeIDTLV, sensorEntryTLV}))

	decoded, err := messages.DecodeSensorUpdateMessage(encodedTLV, false)
	assert.NoError(err)
	assert.NotNil(decoded)
	assert.NotNil(decoded.NodeID)
	assert.Equal(uint8(0x07), *decoded.NodeID)
	assert.Equal(uint8(0x02), decoded.Sensor.ID)
	assert.Equal("HUMIDITY", decoded.Sensor.Type)
	assert.Equal("%", decoded.Sensor.Unit)
	assert.Equal(int32(60), decoded.Sensor.Value)
}

// -------------------------------------- Negative tests --------------------------------------

func TestDecodeSensorUpdateMessage_InvalidArgument(t *testing.T) {
	assert := assert.New(t)

	// nil TLV
	_, err := messages.DecodeSensorUpdateMessage(nil, false)
	assert.Error(err)

	// wrong type
	wrongType, _ := encoding.NewTLV(uint8(constants.COMMAND), []byte{})
	_, err = messages.DecodeSensorUpdateMessage(wrongType, false)
	assert.Error(err)

	// missing sensor entry
	nodeIDTLV, _ := encoding.NewTLV(uint8(constants.SINGLE_NODE), []byte{0x07})
	value := encoding.EncodeMultipleTLVs([]encoding.TLV{nodeIDTLV})
	sensorUpdateTLV, _ := encoding.NewTLV(uint8(constants.SENSOR_UPDATE), value)
	_, err = messages.DecodeSensorUpdateMessage(sensorUpdateTLV, false)
	assert.Error(err)

	// extra unknown TLV
	unknownTLV, _ := encoding.NewTLV(0xFF, []byte{0x01})
	value = encoding.EncodeMultipleTLVs([]encoding.TLV{unknownTLV})
	sensorUpdateTLV, _ = encoding.NewTLV(uint8(constants.SENSOR_UPDATE), value)
	_, err = messages.DecodeSensorUpdateMessage(sensorUpdateTLV, false)
	assert.Error(err)

	// invalid node ID length (should be 1 byte)
	invalidNodeIDTLV, _ := encoding.NewTLV(uint8(constants.SINGLE_NODE), []byte{0x07, 0x08})
	sensorEntryTLV, _ := encoding.EncodeSensorEntry(entity.Sensor[any]{
		ID:    0x01,
		Type:  "TEMPERATURE",
		Unit:  "°C",
		Value: float32(22.5),
	})
	value = encoding.EncodeMultipleTLVs([]encoding.TLV{invalidNodeIDTLV, sensorEntryTLV})
	sensorUpdateTLV, _ = encoding.NewTLV(uint8(constants.SENSOR_UPDATE), value)
	_, err = messages.DecodeSensorUpdateMessage(sensorUpdateTLV, false)
	assert.Error(err)

}
