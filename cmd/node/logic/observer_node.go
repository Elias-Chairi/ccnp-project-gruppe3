package logic

import (
	"log"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/entity"
)

type node struct {
	entity.Node
	Sensors []*observableSensor[any]
}

type nodeOption func(*node)

// WithSensors sets the sensors for the Node and assigns the Node as their observer.
func WithSensors(sensors ...entity.Sensor[any]) nodeOption {
	return func(n *node) {
		for _, sen := range sensors {
			s, err := NewObservableSensor(&sen, n)
			if err != nil {
				log.Println("Error creating observable sensor:", err)
				continue
			}
			n.Sensors = append(n.Sensors, s)
		}
	}
}

// WithActuators sets the actuators for the Node.
func WithActuators(a ...*entity.Actuator[any]) nodeOption {
	return func(n *node) { n.Actuators = a }
}

// NewNode creates a new Node instance with the given options.
func NewNode(opts ...nodeOption) *node {
	node := &node{}
	for _, opt := range opts {
		opt(node)
	}
	return node
}

// UpdateSensorValue updates the value of a Sensor and notifies observers.
func (n *node) UpdateSensorValue(sensor entity.Sensor[any]) {
	log.Println("Sensor updated:", sensor.Type, sensor.Value)
}
