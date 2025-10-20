package entity

type Sensor[T any] struct {
	ID    uint8
	Type  string
	Unit  string
	Value T
}
