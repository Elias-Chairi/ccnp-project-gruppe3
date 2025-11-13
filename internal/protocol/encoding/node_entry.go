package encoding

import (
	"fmt"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/entity"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/constants"
)

func EncodeNodeEntry(n entity.Node) (TLV, error) {
	var tlvs []TLV

	// Local ID
	localIDTLV, err := NewTLV(uint8(constants.NODE_ID), EncodeByte(n.ID))
	if err != nil {
		return nil, fmt.Errorf("failed to encode local ID TLV: %w", err)
	}
	tlvs = append(tlvs, localIDTLV)

	// Sensors
	if len(n.Sensors) > 0 {
		var sensorTLVs []TLV
		for _, sensor := range n.Sensors {
			sensorTLV, err := EncodeSensorEntry(sensor)
			if err != nil {
				return nil, fmt.Errorf("failed to encode sensor entry: %w", err)
			}
			sensorTLVs = append(sensorTLVs, sensorTLV)
		}
		sensorsTLV, err := NewTLV(uint8(constants.SENSOR_ENTRIES), EncodeMultipleTLVs(sensorTLVs))
		if err != nil {
			return nil, fmt.Errorf("failed to encode sensors TLV: %w", err)
		}
		tlvs = append(tlvs, sensorsTLV)
	}

	// Actuators
	if len(n.Actuators) > 0 {
		var actuatorTLVs []TLV
		for _, actuator := range n.Actuators {
			actuatorTLV, err := EncodeActuatorEntry(actuator)
			if err != nil {
				return nil, fmt.Errorf("failed to encode actuator entry: %w", err)
			}
			actuatorTLVs = append(actuatorTLVs, actuatorTLV)
		}
		actuatorsTLV, err := NewTLV(uint8(constants.ACTUATOR_ENTRIES), EncodeMultipleTLVs(actuatorTLVs))
		if err != nil {
			return nil, fmt.Errorf("failed to encode actuators TLV: %w", err)
		}
		tlvs = append(tlvs, actuatorsTLV)
	}

	// Combine all into Node Entry TLV
	return NewTLV(uint8(constants.NODE_ENTRY), EncodeMultipleTLVs(tlvs))
}

func DecodeNodeEntry(t TLV) (*entity.Node, error) {
	if t == nil {
		return nil, fmt.Errorf("data is nil")
	}

	if t.Type() != constants.NODE_ENTRY {
		return nil, fmt.Errorf("invalid TLV type: expected NODE_ENTRY, got %d", t.Type())
	}

	inner, err := DecodeMultipleTLVs(t.Value())
	if err != nil {
		return nil, fmt.Errorf("failed to decode inner TLVs: %w", err)
	}

	node := &entity.Node{}
	for _, innerTLV := range inner {
		switch innerTLV.Type() {
		case uint8(constants.NODE_ID):
			id, err := DecodeByte(innerTLV.Value())
			if err != nil {
				return nil, fmt.Errorf("failed to decode NODE_ID: %w", err)
			}
			node.ID = id
		case uint8(constants.SENSOR_ENTRIES):
			sensorTLVs, err := DecodeMultipleTLVs(innerTLV.Value())
			if err != nil {
				return nil, fmt.Errorf("failed to decode SENSOR_ENTRIES: %w", err)
			}
			for _, sensorTLV := range sensorTLVs {
				sensor, err := DecodeSensorEntry(sensorTLV)
				if err != nil {
					return nil, fmt.Errorf("failed to decode sensor entry: %w", err)
				}
				node.Sensors = append(node.Sensors, *sensor)
			}
		case uint8(constants.ACTUATOR_ENTRIES):
			actuatorTLVs, err := DecodeMultipleTLVs(innerTLV.Value())
			if err != nil {
				return nil, fmt.Errorf("failed to decode ACTUATOR_ENTRIES: %w", err)
			}
			for _, actuatorTLV := range actuatorTLVs {
				actuator, err := DecodeActuatorEntry(actuatorTLV)
				if err != nil {
					return nil, fmt.Errorf("failed to decode actuator entry: %w", err)
				}
				node.Actuators = append(node.Actuators, *actuator)
			}
		default:
			return nil, fmt.Errorf("unknown TLV type in NODE_ENTRY: %d", innerTLV.Type())
		}
	}

	return node, nil
}
