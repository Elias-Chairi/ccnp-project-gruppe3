package encoding_test

import (
	"testing"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/entity"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/constants"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/encoding"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/tlv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// -------------------------------------- Positive tests --------------------------------------

var sensorTests = []struct {
	name   string
	sensor entity.Sensor[any]
}{
	{
		name: "Float value",
		sensor: entity.Sensor[any]{
			ID:    0x11,
			Type:  "temperature",
			Unit:  "°C",
			Value: float32(19.75),
		},
	},
	{
		name: "Integer value",
		sensor: entity.Sensor[any]{
			ID:    0x12,
			Type:  "humidity",
			Unit:  "%",
			Value: int32(55),
		},
	},
	{
		name: "String value",
		sensor: entity.Sensor[any]{
			ID:    0x13,
			Type:  "status",
			Unit:  "",
			Value: "ok",
		},
	},
	{
		name: "Boolean value",
		sensor: entity.Sensor[any]{
			ID:    0x14,
			Type:  "motion",
			Unit:  "detected",
			Value: true,
		},
	},
}

func TestEncodeSensorEntry(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	for _, tt := range sensorTests {
		t.Run(tt.name, func(t *testing.T) {
			encoded, err := encoding.EncodeSensorEntry(tt.sensor)
			require.NoError(err)
			require.NotNil(encoded)

			assert.Equal(uint8(constants.SENSOR_ENTRY), encoded.Type())

			innerTLVs, err := tlv.DecodeMultipleTLVs(encoded.Value())
			require.NoError(err)
			require.Len(innerTLVs, 4)

			fields := make(map[constants.SensorField]tlv.TLV)
			for _, inner := range innerTLVs {
				fields[constants.SensorField(inner.Type())] = inner
			}

			idTLV, ok := fields[constants.SENSOR_ID]
			require.True(ok)
			assert.Equal([]byte{tt.sensor.ID}, idTLV.Value())

			typeTLV, ok := fields[constants.SENSOR_TYPE]
			require.True(ok)
			assert.Equal(tt.sensor.Type, string(typeTLV.Value()))

			unitTLV, ok := fields[constants.SENSOR_UNIT]
			require.True(ok)
			assert.Equal(tt.sensor.Unit, string(unitTLV.Value()))

			valueTLV, ok := fields[constants.SENSOR_VALUE]
			require.True(ok)
			valueDataTLV, err := tlv.DecodeTLV(valueTLV.Value())
			require.NoError(err)
			valueData, err := encoding.DecodeAny(valueDataTLV)
			require.NoError(err)
			assert.Equal(tt.sensor.Value, valueData)
		})
	}
}

func TestDecodeSensorEntry(t *testing.T) {
	for _, tt := range sensorTests {
		t.Run(tt.name, func(t *testing.T) {
			idTLV, _ := tlv.NewTLV(uint8(constants.SENSOR_ID), []byte{tt.sensor.ID})
			typeTLV, _ := tlv.NewTLV(uint8(constants.SENSOR_TYPE), []byte(tt.sensor.Type))
			unitTLV, _ := tlv.NewTLV(uint8(constants.SENSOR_UNIT), []byte(tt.sensor.Unit))
			valuePayload, _ := encoding.EncodeAny(tt.sensor.Value)
			valueTLV, _ := tlv.NewTLV(uint8(constants.SENSOR_VALUE), valuePayload.Encode())
			entryTLV, _ := tlv.NewTLV(constants.SENSOR_ENTRY, tlv.EncodeMultipleTLVs([]tlv.TLV{idTLV, typeTLV, unitTLV, valueTLV}))

			decoded, err := encoding.DecodeSensorEntry(entryTLV)
			require.NoError(t, err)
			require.NotNil(t, decoded)

			assert.Equal(t, tt.sensor.ID, decoded.ID)
			assert.Equal(t, tt.sensor.Type, decoded.Type)
			assert.Equal(t, tt.sensor.Unit, decoded.Unit)
			assert.Equal(t, tt.sensor.Value, decoded.Value)
		})
	}
}

// -------------------------------------- Negative tests --------------------------------------

func TestEncodeSensorEntry_InvalidArgument(t *testing.T) {
	// invalid type
	_, err := encoding.EncodeSensorEntry(entity.Sensor[any]{
		ID:    0x21,
		Type:  "custom",
		Unit:  "unit",
		Value: struct{}{},
	})
	assert.Error(t, err)
}

func TestDecodeSensorEntry_InvalidArgument(t *testing.T) {
	// nil TLV
	_, err := encoding.DecodeSensorEntry(nil)
	assert.Error(t, err)

	// wrong type
	wrongType, _ := tlv.NewTLV(uint8(constants.ACTUATOR_ENTRY), []byte{0x01})
	_, err = encoding.DecodeSensorEntry(wrongType)
	assert.Error(t, err)
}
