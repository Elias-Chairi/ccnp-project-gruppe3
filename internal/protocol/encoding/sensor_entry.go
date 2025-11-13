package encoding

import (
	"fmt"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/entity"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/constants"
)

// EncodeSensorEntry encodes a sensor into a SENSOR_ENTRY TLV with nested fields.
func EncodeSensorEntry(s entity.Sensor[any]) (TLV, error) {
	var tlvs []TLV

	// Sensor ID
	sid, err := NewTLV(uint8(constants.SENSOR_ID), EncodeByte(s.ID))
	if err != nil {
		return nil, fmt.Errorf("failed to create SENSOR_ID TLV: %w", err)
	}
	tlvs = append(tlvs, sid)

	// Sensor Type
	st, err := NewTLV(uint8(constants.SENSOR_TYPE), []byte(s.Type))
	if err != nil {
		return nil, fmt.Errorf("failed to create SENSOR_TYPE TLV: %w", err)
	}
	tlvs = append(tlvs, st)

	// Sensor Unit
	su, err := NewTLV(uint8(constants.SENSOR_UNIT), []byte(s.Unit))
	if err != nil {
		return nil, fmt.Errorf("failed to create SENSOR_UNIT TLV: %w", err)
	}
	tlvs = append(tlvs, su)

	// Sensor Value (with data type TLV)
	svTLV, err := EncodeAny(s.Value)
	if err != nil {
		return nil, fmt.Errorf("failed to encode sensor value: %w", err)
	}
	sv, err := NewTLV(uint8(constants.SENSOR_VALUE), svTLV.Encode())
	if err != nil {
		return nil, fmt.Errorf("failed to create SENSOR_VALUE TLV: %w", err)
	}
	tlvs = append(tlvs, sv)

	// Wrap all fields in SENSOR_ENTRY TLV
	sensorEntryTLV, err := NewTLV(constants.SENSOR_ENTRY, EncodeMultipleTLVs(tlvs))
	if err != nil {
		return nil, err
	}
	return sensorEntryTLV, nil
}

// DecodeSensorEntry decodes a SENSOR_ENTRY TLV into a Sensor and validates field lengths.
func DecodeSensorEntry(t TLV) (*entity.Sensor[any], error) {
	if t == nil {
		return nil, fmt.Errorf("data is nil")
	}

	if t.Type() != constants.SENSOR_ENTRY {
		return nil, fmt.Errorf("expected SENSOR_ENTRY type, got %x", t.Type())
	}

	tlvs, err := DecodeMultipleTLVs(t.Value())
	if err != nil {
		return nil, fmt.Errorf("failed to decode multiple TLVs: %w", err)
	}

	s := &entity.Sensor[any]{}
	for _, innerTLV := range tlvs {
		switch constants.SensorField(innerTLV.Type()) {
		case constants.SENSOR_ID:
			id, err := DecodeByte(innerTLV.Value())
			if err != nil {
				return nil, fmt.Errorf("failed to decode sensor ID: %w", err)
			}
			s.ID = id
		case constants.SENSOR_TYPE:
			s.Type = string(innerTLV.Value())
		case constants.SENSOR_UNIT:
			s.Unit = string(innerTLV.Value())
		case constants.SENSOR_VALUE:
			valueTLV, err := DecodeTLV(innerTLV.Value())
			if err != nil {
				return nil, fmt.Errorf("failed to decode sensor value TLV: %w", err)
			}
			s.Value, err = DecodeAny(valueTLV)
			if err != nil {
				return nil, fmt.Errorf("failed to decode sensor value: %w", err)
			}
		}
	}

	return s, nil
}
