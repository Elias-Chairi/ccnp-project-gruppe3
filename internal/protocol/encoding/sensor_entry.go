package encoding

import (
	"fmt"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/entity"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/constants"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/tlv"
)

// Encode encodes the sensor entry as nested TLVs
func EncodeSensorEntry(s entity.Sensor[any]) (tlv.TLV, error) {
	var tlvs []tlv.TLV

	// Sensor ID
	sid, err := tlv.NewTLV(uint8(constants.SENSOR_ID), EncodeByte(s.ID))
	if err != nil {
		return nil, fmt.Errorf("failed to create SENSOR_ID TLV: %w", err)
	}
	tlvs = append(tlvs, sid)

	// Sensor Type
	st, err := tlv.NewTLV(uint8(constants.SENSOR_TYPE), []byte(s.Type))
	if err != nil {
		return nil, fmt.Errorf("failed to create SENSOR_TYPE TLV: %w", err)
	}
	tlvs = append(tlvs, st)

	// Sensor Unit
	su, err := tlv.NewTLV(uint8(constants.SENSOR_UNIT), []byte(s.Unit))
	if err != nil {
		return nil, fmt.Errorf("failed to create SENSOR_UNIT TLV: %w", err)
	}
	tlvs = append(tlvs, su)

	// Sensor Value (with data type TLV)
	svTLV, err := EncodeAny(s.Value)
	if err != nil {
		return nil, fmt.Errorf("failed to encode sensor value: %w", err)
	}
	sv, err := tlv.NewTLV(uint8(constants.SENSOR_VALUE), svTLV.Encode())
	if err != nil {
		return nil, fmt.Errorf("failed to create SENSOR_VALUE TLV: %w", err)
	}
	tlvs = append(tlvs, sv)

	// Wrap all fields in SENSOR_ENTRY TLV
	sensorEntryTLV, err := tlv.NewTLV(constants.SENSOR_ENTRY, tlv.EncodeMultipleTLVs(tlvs))
	if err != nil {
		return nil, err
	}
	return sensorEntryTLV, nil
}

// DecodeSensorEntry decodes a sensor entry from TLV data
func DecodeSensorEntry(t tlv.TLV) (*entity.Sensor[any], error) {
	if t == nil {
		return nil, fmt.Errorf("data is nil")
	}

	if t.Type() != constants.SENSOR_ENTRY {
		return nil, fmt.Errorf("expected SENSOR_ENTRY type, got %x", t.Type())
	}

	tlvs, err := tlv.DecodeMultipleTLVs(t.Value())
	if err != nil {
		return nil, fmt.Errorf("failed to decode multiple TLVs: %w", err)
	}

	s := &entity.Sensor[any]{}
	for _, t := range tlvs {
		switch constants.SensorField(t.Type()) {
		case constants.SENSOR_ID:
			if t.Length() != 1 {
				return nil, fmt.Errorf("invalid sensor ID length")
			}
			s.ID = t.Value()[0]
		case constants.SENSOR_TYPE:
			s.Type = string(t.Value())
		case constants.SENSOR_UNIT:
			s.Unit = string(t.Value())
		case constants.SENSOR_VALUE:
			s.Value, err = DecodeAny(t)
			if err != nil {
				return nil, fmt.Errorf("failed to decode sensor value: %w", err)
			}
		}
	}

	return s, nil
}
