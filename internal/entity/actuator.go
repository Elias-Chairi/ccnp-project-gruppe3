package entity

type Actuator[T any] struct {
	ID    uint8
	Type  string
	Unit  string
	State T
}
