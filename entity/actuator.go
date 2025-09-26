package entity

type actuator[T any] struct {
	id           int
	actuatorType string
	unit         string
	state        T
}

// NewActuator creates a new Actuator instance with the given parameters.
func NewActuator[T any](id int, actuatorType string, unit string, initialState T) *actuator[T] {
	return &actuator[T]{
		id:           id,
		actuatorType: actuatorType,
		unit:         unit,
		state:        initialState,
	}
}

// GetID returns the ID of the Actuator.
func (a *actuator[T]) GetID() int {
	return a.id
}

// SetID sets the ID of the Actuator.
func (a *actuator[T]) GetType() string {
	return a.actuatorType
}

// GetUnit returns the unit of the Actuator.
func (a *actuator[T]) GetUnit() string {
	return a.unit
}

// GetState returns the current state of the Actuator.
func (a *actuator[T]) GetState() T {
	return a.state
}

// SetState sets the state of the Actuator.
func (a *actuator[T]) SetState(newState T) {
	a.state = newState
}
