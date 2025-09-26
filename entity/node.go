package entity

import "log"

type node struct {
	id        int
	actuators []*actuator[any]
	sensors   []*sensor[any]
}

type nodeOption func(*node)

// WithSensors sets the sensors for the Node and assigns the Node as their observer.
func WithSensors(sensors ...*sensor[any]) nodeOption {
	return func(n *node) {
		for _, sensor := range sensors {
			sensor.SetObserver(n)
			n.sensors = append(n.sensors, sensor)
		}
	}
}

// WithActuators sets the actuators for the Node.
func WithActuators(a ...*actuator[any]) nodeOption {
	return func(n *node) { n.actuators = a }
}

// NewNode creates a new Node instance with the given options.
func NewNode(opts ...nodeOption) *node {
	node := &node{}
	for _, opt := range opts {
		opt(node)
	}
	return node
}

// SetID sets the ID of the Node.
func (n *node) SetID(id int) {
	n.id = id
}

// GetID returns the ID of the Node.
func (n *node) GetID() int {
	return n.id
}

func (n *node) GetActuators() []*actuator[any] {
	return n.actuators
}

// UpdateSensorValue updates the value of a Sensor and notifies observers.
func (n *node) UpdateSensorValue(sensor *sensor[any]) {
	log.Println("Sensor updated:", sensor.GetType(), sensor.GetValue())
}
