package entity

type Actuator[T any] struct {
	ID    string
	Type  string
	State T
}
