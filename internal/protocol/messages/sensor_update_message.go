package messages

import (
	"fmt"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/entity"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/constants"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/encoding"
)

// sensorUpdateMessage represents a SENSOR_UPDATE message (Type SENSOR_UPDATE).
//
// Purpose: Transmit sensor readings.
//   - Node → Server: Sensor entry without node ID
//   - Server → Control: NodeID (0x51) + Sensor entry
//
// TLV: [Type:SENSOR_UPDATE][Length:N][Value: (Optional NodeID TLV) + SensorEntry TLV]
type sensorUpdateMessage struct {
	NodeID *uint8 // Only present in server->control direction
	Sensor entity.Sensor[any]
}

// NewSensorUpdateMessage creates a new SENSOR_UPDATE message without node ID (node → server).
func NewSensorUpdateMessage(sensor entity.Sensor[any]) *sensorUpdateMessage {
	return &sensorUpdateMessage{
		NodeID: nil,
		Sensor: sensor,
	}
}

// NewSensorUpdateMessageWithNode creates a new SENSOR_UPDATE message with node ID (server → control).
func NewSensorUpdateMessageWithNode(nodeID uint8, sensor entity.Sensor[any]) *sensorUpdateMessage {
	return &sensorUpdateMessage{
		NodeID: &nodeID,
		Sensor: sensor,
	}
}

// Type returns the message type.
func (m *sensorUpdateMessage) Type() constants.MessageType {
	return constants.SENSOR_UPDATE
}

// Encode encodes the SENSOR_UPDATE message to bytes.
func (m *sensorUpdateMessage) Encode() ([]byte, error) {
	var tlvs []encoding.TLV

	// Encode node ID if present
	if m.NodeID != nil {
		tlv, err := encoding.NewTLV(uint8(constants.SINGLE_NODE), []byte{*m.NodeID})
		if err != nil {
			return nil, fmt.Errorf("failed to create SINGLE_NODE TLV: %w", err)
		}
		tlvs = append(tlvs, tlv)
	}

	// Required sensor entry
	sensorTLV, err := encoding.EncodeSensorEntry(m.Sensor)
	if err != nil {
		return nil, fmt.Errorf("failed to encode sensor entry: %w", err)
	}
	tlvs = append(tlvs, sensorTLV)

	// Wrap in SENSOR_UPDATE TLV
	value := encoding.EncodeMultipleTLVs(tlvs)
	mainTLV, err := encoding.NewTLV(uint8(constants.SENSOR_UPDATE), value)
	if err != nil {
		return nil, fmt.Errorf("failed to create SENSOR_UPDATE TLV: %w", err)
	}
	return mainTLV.Encode(), nil
}

// DecodeSensorUpdateMessage decodes a SENSOR_UPDATE message from TLV.
// Validates type (SENSOR_UPDATE), optional NodeID TLV length (1), and presence of sensor entry.
func DecodeSensorUpdateMessage(t encoding.TLV) (*sensorUpdateMessage, error) {
	if t == nil {
		return nil, fmt.Errorf("data is nil")
	}
	if t.Type() != uint8(constants.SENSOR_UPDATE) {
		return nil, fmt.Errorf("expected SENSOR_UPDATE type, got %x", t.Type())
	}

	inner, err := encoding.DecodeMultipleTLVs(t.Value())
	if err != nil {
		return nil, err
	}

	msg := &sensorUpdateMessage{}
	var sensor *entity.Sensor[any]

	for _, innerTLV := range inner {
		switch innerTLV.Type() {
		case uint8(constants.SINGLE_NODE):
			if innerTLV.Length() != 1 {
				return nil, fmt.Errorf("invalid node ID length")
			}
			id := innerTLV.Value()[0]
			msg.NodeID = &id
		case uint8(constants.SENSOR_ENTRY):
			if sensor != nil {
				return nil, fmt.Errorf("duplicate SENSOR_ENTRY")
			}
			sensor, err = encoding.DecodeSensorEntry(innerTLV)
			if err != nil {
				return nil, fmt.Errorf("failed to decode sensor entry: %w", err)
			}
		default:
			return nil, fmt.Errorf("unknown TLV %x inside SENSOR_UPDATE", innerTLV.Type())
		}
	}

	if sensor == nil {
		return nil, fmt.Errorf("sensor entry not found in SENSOR_UPDATE message")
	}
	msg.Sensor = *sensor

	return msg, nil
}
