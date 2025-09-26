package entity

type Actuator[T any] struct {
	id           int
	actuatorType string
	unit         string
	state        T
}

func NewActuator[T any](id int, actuatorType string, unit string, initialState T) *Actuator[T] {
	return &Actuator[T]{
		id:           id,
		actuatorType: actuatorType,
		unit:         unit,
		state:        initialState,
	}
}

func (a *Actuator[T]) GetID() int {
	return a.id
}

func (a *Actuator[T]) GetType() string {
	return a.actuatorType
}

func (a *Actuator[T]) GetUnit() string {
	return a.unit
}

func (a *Actuator[T]) GetState() T {
	return a.state
}

func (a *Actuator[T]) SetState(newState T) {
	a.state = newState
}
