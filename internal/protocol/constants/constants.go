// Package constants defines all TLV type codes used by the protocol.
//
// Code space layout (1-byte codes, non-overlapping 16-ID blocks):
//   - MessageType      : 0x40-0x4F
//   - NodeSelector     : 0x50-0x5F
//   - ActuatorSelector : 0x60-0x6F
//   - SensorField      : 0x70-0x7F
//   - ActuatorField    : 0x80-0x8F
//   - DataType         : 0x90-0x9F
//   - AckErrorCode     : 0xA0-0xAF
//
// Some groups provide an IsValid method to validate codes at runtime.
package constants

// Message Types, Range: 0x40-0x4F
type MessageType uint8

const (
	DISCOVERY        MessageType = 0x41
	REGISTER_NODE    MessageType = 0x42
	REGISTER_CONTROL MessageType = 0x43
	SENSOR_UPDATE    MessageType = 0x44
	ACTUATOR_UPDATE  MessageType = 0x45
	NODE_ADDED       MessageType = 0x46
	NODE_REMOVED     MessageType = 0x47
	COMMAND          MessageType = 0x48
	ACK_ERROR        MessageType = 0x49
)

func (m MessageType) IsValid() bool {
	switch m {
	case DISCOVERY, REGISTER_NODE, REGISTER_CONTROL, SENSOR_UPDATE, COMMAND, ACK_ERROR:
		return true
	default:
		return false
	}
}

// Node Selector Types, Range: 0x50-0x5F
type NodeSelector uint8

const (
	SINGLE_NODE NodeSelector = 0x51
	NODE_LIST   NodeSelector = 0x52
	ALL_NODES   NodeSelector = 0x53
)

func (s NodeSelector) IsValid() bool {
	switch s {
	case SINGLE_NODE, NODE_LIST, ALL_NODES:
		return true
	default:
		return false
	}
}

// Actuator Selector Types, Range: 0x60-0x6F
type ActuatorSelector uint8

const (
	SINGLE_ACTUATOR ActuatorSelector = 0x61
	ACTUATOR_LIST   ActuatorSelector = 0x62
	ACTUATOR_TYPE   ActuatorSelector = 0x63
	ALL_ACTUATORS   ActuatorSelector = 0x64
)

func (s ActuatorSelector) IsValid() bool {
	switch s {
	case SINGLE_ACTUATOR, ACTUATOR_LIST, ACTUATOR_TYPE, ALL_ACTUATORS:
		return true
	default:
		return false
	}
}

// Node field codes, Range: 0x30-0x3F
type NodeField uint8

const (
	NODE_ENTRY       uint8     = 0x30
	NODE_ID          NodeField = 0x31
	SENSOR_ENTRIES   NodeField = 0x32
	ACTUATOR_ENTRIES NodeField = 0x33
)

// Sensor Field codes, Range: 0x70-0x7F
type SensorField uint8

const (
	SENSOR_ENTRY uint8       = 0x70
	SENSOR_ID    SensorField = 0x71
	SENSOR_TYPE  SensorField = 0x72
	SENSOR_UNIT  SensorField = 0x73
	SENSOR_VALUE SensorField = 0x74
)

// Actuator Field codes, Range: 0x80-0x8F
type ActuatorField uint8

const (
	ACTUATOR_ENTRY      uint8         = 0x80
	ACTUATOR_ID         ActuatorField = 0x81
	ACTUATOR_TYPE_FIELD ActuatorField = 0x82
	ACTUATOR_UNIT       ActuatorField = 0x83
	ACTUATOR_STATE      ActuatorField = 0x84
)

// Data Type TLVs, Range: 0x90-0x9F
type DataType uint8

const (
	DATA_TYPE_INTEGER DataType = 0x91
	DATA_TYPE_FLOAT   DataType = 0x92
	DATA_TYPE_STRING  DataType = 0x93
	DATA_TYPE_BOOLEAN DataType = 0x94
)

func (d DataType) IsValid() bool {
	switch d {
	case DATA_TYPE_INTEGER, DATA_TYPE_FLOAT, DATA_TYPE_STRING, DATA_TYPE_BOOLEAN:
		return true
	default:
		return false
	}
}

// ACK/ERROR Codes, Range: 0xA0-0xAF
type AckErrorCode uint8

const (
	ACK_SUCCESS              AckErrorCode = 0xA0
	ERR_UNKNOWN_NODE_ID      AckErrorCode = 0xA1
	ERR_UNKNOWN_SENSOR_ID    AckErrorCode = 0xA2
	ERR_UNKNOWN_ACTUATOR_ID  AckErrorCode = 0xA3
	ERR_INVALID_ACTION       AckErrorCode = 0xA4
	ERR_INVALID_VALUE        AckErrorCode = 0xA5
	ERR_MALFORMED_MESSAGE    AckErrorCode = 0xA6
	ERR_INVALID_MESSAGE_TYPE AckErrorCode = 0xA7
)

func (c AckErrorCode) IsValid() bool {
	switch c {
	case ACK_SUCCESS, ERR_UNKNOWN_NODE_ID, ERR_UNKNOWN_SENSOR_ID, ERR_UNKNOWN_ACTUATOR_ID, ERR_INVALID_ACTION, ERR_INVALID_VALUE, ERR_MALFORMED_MESSAGE, ERR_INVALID_MESSAGE_TYPE:
		return true
	default:
		return false
	}
}
