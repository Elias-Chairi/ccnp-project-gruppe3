package entity

import "errors"

type Node struct {
	ID        string
	Actuators []Actuator[any]
	Sensors   []Sensor[any]
}

// NewNode creates a new Node instance with initialized slices for Actuators and Sensors.
func NewNode() *Node {
	return &Node{
		Actuators: []Actuator[any]{},
		Sensors:   []Sensor[any]{},
	}
}

// SetID sets the ID of the Node. It returns an error if the provided ID is empty.
func (n *Node) SetID(id string) error {
	if id == "" {
		return errors.New("ID cannot be empty")
	}
	n.ID = id
	return nil
}
