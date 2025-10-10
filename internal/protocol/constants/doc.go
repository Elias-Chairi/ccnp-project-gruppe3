// Package constants defines all TLV type codes used by the protocol.
//
// Code space layout (1-byte codes, non-overlapping 16-ID blocks):
//   - MessageType      0x40-0x4F  (DISCOVERY=0x41, REGISTER_NODE=0x42, REGISTER_CONTROL=0x43,
//     SENSOR_UPDATE=0x44, COMMAND=0x45, ACK_ERROR=0x46)
//   - NodeSelector     0x50-0x5F  (SINGLE_NODE=0x51, NODE_LIST=0x52, ALL_NODES=0x53)
//   - ActuatorSelector 0x60-0x6F  (SINGLE_ACTUATOR=0x61, ACTUATOR_LIST=0x62,
//     ACTUATOR_TYPE=0x63, ALL_ACTUATORS=0x64)
//   - SensorField      0x70-0x7F  (SENSOR_ENTRY=0x70, SENSOR_ID=0x71, SENSOR_TYPE=0x72,
//     SENSOR_UNIT=0x73, SENSOR_VALUE=0x74)
//   - ActuatorField    0x80-0x8F  (ACTUATOR_ENTRY=0x80, ACTUATOR_ID=0x81,
//     ACTUATOR_TYPE_FIELD=0x82, ACTUATOR_UNIT=0x83,
//     ACTUATOR_STATE=0x84)
//   - DataType         0x90-0x9F  (INTEGER=0x91, FLOAT=0x92, STRING=0x93, BOOLEAN=0x94)
//   - AckErrorCode     0xA0-0xAF  (ACK_SUCCESS=0xA0, ERR_*=0xA1..)
//
// Each group provides an IsValid method to validate codes at runtime.
// The value 0x00 is globally reserved to avoid zero-initialized memory
// accidentally being treated as a valid code.
package constants
