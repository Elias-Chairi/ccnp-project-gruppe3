package entity

type Node struct {
	ID        uint8
	Actuators []*Actuator[any]
	Sensors   []*Sensor[any]
}
