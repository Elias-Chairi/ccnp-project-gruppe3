package entity

type Sensor[T any] struct {
	ID    string
	Type  string
	Value T
}
