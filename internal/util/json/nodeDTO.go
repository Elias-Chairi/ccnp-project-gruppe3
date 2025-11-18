package json

type NodeDTO struct {
	ID        uint8         `json:"id"`
	Actuators []ActuatorDTO `json:"actuators"`
	Sensors   []SensorDTO   `json:"sensors"`
}
