package messages

import (
	"fmt"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/entity"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/constants"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/encoding"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/tlv"
)

// registerNodeMessage represents a REGISTER_NODE message (Type REGISTER_NODE).
//
// Purpose: Node registers it self with the server over TCP, providing its sensors and actuators.
// TLV: [Type:REGISTER_NODE][Length:N][Value: SensorEntry/ActuatorEntry TLVs]
type registerNodeMessage struct {
	Sensors   []entity.Sensor[any]
	Actuators []entity.Actuator[any]
}

// NewRegisterNodeMessage creates a new REGISTER_NODE message.
func NewRegisterNodeMessage(sensors []entity.Sensor[any], actuators []entity.Actuator[any]) (*registerNodeMessage, error) {
	if len(sensors) == 0 && len(actuators) == 0 {
		return nil, fmt.Errorf("at least one of sensors or actuators must be provided")
	}
	return &registerNodeMessage{
		Sensors:   sensors,
		Actuators: actuators,
	}, nil
}

// Encode encodes the REGISTER_NODE message to bytes.
func (m *registerNodeMessage) Encode() ([]byte, error) {

	// Encode sensor list
	var tlvs []tlv.TLV
	for _, sensor := range m.Sensors {
		tlv, err := encoding.EncodeSensorEntry(sensor)
		if err != nil {
			return nil, fmt.Errorf("failed to encode sensor entry: %w", err)
		}
		tlvs = append(tlvs, tlv)
	}

	// Encode actuator list
	for _, actuator := range m.Actuators {
		tlv, err := encoding.EncodeActuatorEntry(actuator)
		if err != nil {
			return nil, fmt.Errorf("failed to encode actuator entry: %w", err)
		}
		tlvs = append(tlvs, tlv)
	}

	value := tlv.EncodeMultipleTLVs(tlvs)
	mainTLV, err := tlv.NewTLV(uint8(constants.REGISTER_NODE), value)
	if err != nil {
		return nil, fmt.Errorf("failed to create REGISTER_NODE TLV: %w", err)
	}
	return mainTLV.Encode(), nil
}

// Type returns the message type.
func (m *registerNodeMessage) Type() constants.MessageType {
	return constants.REGISTER_NODE
}

// DecodeRegisterNodeMessage decodes a REGISTER_NODE message from TLV.
// Validates type (REGISTER_NODE) and decodes nested sensor/actuator entries.
func DecodeRegisterNodeMessage(t tlv.TLV) (*registerNodeMessage, error) {
	if t == nil {
		return nil, fmt.Errorf("nil TLV provided")
	}
	if t.Type() != uint8(constants.REGISTER_NODE) {
		return nil, fmt.Errorf("expected REGISTER_NODE type, got %x", t.Type())
	}

	tlvs, err := tlv.DecodeMultipleTLVs(t.Value())
	if err != nil {
		return nil, fmt.Errorf("failed to decode TLVs: %w", err)
	}

	msg := &registerNodeMessage{}

	for _, innerTLV := range tlvs {
		switch innerTLV.Type() {
		case constants.SENSOR_ENTRY:
			sensor, err := encoding.DecodeSensorEntry(innerTLV)
			if err != nil {
				return nil, fmt.Errorf("failed to decode sensor entry: %w", err)
			}
			msg.Sensors = append(msg.Sensors, *sensor)
		case constants.ACTUATOR_ENTRY:
			actuator, err := encoding.DecodeActuatorEntry(innerTLV)
			if err != nil {
				return nil, fmt.Errorf("failed to decode actuator entry: %w", err)
			}
			msg.Actuators = append(msg.Actuators, *actuator)
		default:
			return nil, fmt.Errorf("unexpected TLV type %x in REGISTER_NODE message", innerTLV.Type())
		}
	}

	if len(msg.Sensors) == 0 && len(msg.Actuators) == 0 {
		return nil, fmt.Errorf("REGISTER_NODE message must contain at least one sensor or actuator")
	}

	return msg, nil
}
