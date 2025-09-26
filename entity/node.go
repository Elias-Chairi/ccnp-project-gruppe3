package entity

import "log"

type Node struct {
	id        int
	actuators []Actuator[any]
	sensors   []Sensor[any]
}

// NewNode creates a new Node instance with initialized slices for Actuators and Sensors.
func NewNode(actuators []Actuator[any]) *Node {
	return &Node{
		actuators: actuators,
		sensors:   []Sensor[any]{},
	}
}

// SetID sets the ID of the Node.
func (n *Node) SetID(id int) {
	n.id = id
}

// GetID returns the ID of the Node.
func (n *Node) GetID() int {
	return n.id
}

// AddSensor adds a new Sensor to the Node's sensors slice.
func (n *Node) AddSensor(sensor Sensor[any]) {
	sensor.SetObserver(n)
	n.sensors = append(n.sensors, sensor)
}

// UpdateSensorValue updates the value of a Sensor and notifies observers.
func (n *Node) UpdateSensorValue(sensor *Sensor[any]) {
	log.Println("Sensor updated:", sensor.GetType(), sensor.GetValue())
}
