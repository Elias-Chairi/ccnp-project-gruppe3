package encoding_test

import (
	"testing"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/entity"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/constants"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/encoding"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// -------------------------------------- Positive tests --------------------------------------

var tests = []struct {
	name     string
	actuator entity.Actuator[any]
}{
	{
		name: "Float state",
		actuator: entity.Actuator[any]{
			ID:    0x01,
			Type:  "heater",
			Unit:  "°C",
			State: float32(22.5),
		},
	},
	{
		name: "Integer state",
		actuator: entity.Actuator[any]{
			ID:    0x02,
			Type:  "fan",
			Unit:  "rpm",
			State: int32(1200),
		},
	},
	{
		name: "String state",
		actuator: entity.Actuator[any]{
			ID:    0x03,
			Type:  "light",
			Unit:  "mode",
			State: "bright",
		},
	},
	{
		name: "Boolean state",
		actuator: entity.Actuator[any]{
			ID:    0x04,
			Type:  "switch",
			Unit:  "on/off",
			State: true,
		},
	},
}

func TestEncodeActuatorEntry(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded, err := encoding.EncodeActuatorEntry(tt.actuator)
			require.NoError(err)
			require.NotNil(encoded)

			assert.Equal(uint8(constants.ACTUATOR_ENTRY), encoded.Type())

			innerTLVs, err := encoding.DecodeMultipleTLVs(encoded.Value())
			require.NoError(err)
			require.Len(innerTLVs, 4)

			fieldTLVs := make(map[constants.ActuatorField]encoding.TLV)
			for _, inner := range innerTLVs {
				fieldTLVs[constants.ActuatorField(inner.Type())] = inner
			}

			idTLV, ok := fieldTLVs[constants.ACTUATOR_ID]
			require.True(ok)
			assert.Equal([]byte{tt.actuator.ID}, idTLV.Value())

			typeTLV, ok := fieldTLVs[constants.ACTUATOR_TYPE_FIELD]
			require.True(ok)
			assert.Equal(tt.actuator.Type, string(typeTLV.Value()))

			unitTLV, ok := fieldTLVs[constants.ACTUATOR_UNIT]
			require.True(ok)
			assert.Equal(tt.actuator.Unit, string(unitTLV.Value()))

			stateTLV, ok := fieldTLVs[constants.ACTUATOR_STATE]
			require.True(ok)
			stateValueTLV, err := encoding.DecodeTLV(stateTLV.Value())
			require.NoError(err)
			stateValue, err := encoding.DecodeAny(stateValueTLV)
			require.NoError(err)
			assert.Equal(tt.actuator.State, stateValue)
		})
	}
}

func TestDecodeActuatorEntry(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			idTLV, _ := encoding.NewTLV(uint8(constants.ACTUATOR_ID), []byte{tt.actuator.ID})
			typeTLV, _ := encoding.NewTLV(uint8(constants.ACTUATOR_TYPE_FIELD), []byte(tt.actuator.Type))
			unitTLV, _ := encoding.NewTLV(uint8(constants.ACTUATOR_UNIT), []byte(tt.actuator.Unit))
			StateValueTLV, _ := encoding.EncodeAny(tt.actuator.State)
			stateTLV, _ := encoding.NewTLV(uint8(constants.ACTUATOR_STATE), StateValueTLV.Encode())
			entryTLV, _ := encoding.NewTLV(uint8(constants.ACTUATOR_ENTRY), encoding.EncodeMultipleTLVs([]encoding.TLV{idTLV, typeTLV, unitTLV, stateTLV}))

			decoded, err := encoding.DecodeActuatorEntry(entryTLV)
			require.NoError(t, err)
			require.NotNil(t, decoded)

			assert.Equal(t, tt.actuator.ID, decoded.ID)
			assert.Equal(t, tt.actuator.Type, decoded.Type)
			assert.Equal(t, tt.actuator.Unit, decoded.Unit)
			assert.Equal(t, tt.actuator.State, decoded.State)
		})
	}
}

// -------------------------------------- Negative tests --------------------------------------

func TestEncodeActuatorEntry_InvalidArgument(t *testing.T) {
	// unsupported state type
	_, err := encoding.EncodeActuatorEntry(entity.Actuator[any]{
		ID:    0x01,
		Type:  "custom",
		Unit:  "unit",
		State: struct{}{},
	})
	assert.Error(t, err)
}

func TestDecodeActuatorEntry_InvalidArgument(t *testing.T) {
	// nil TLV
	_, err := encoding.DecodeActuatorEntry(nil)
	assert.Error(t, err)

	// wrong TLV type
	wrongType, _ := encoding.NewTLV(byte(constants.SENSOR_ENTRY), []byte{0x01})
	_, err = encoding.DecodeActuatorEntry(wrongType)
	assert.Error(t, err)
}
