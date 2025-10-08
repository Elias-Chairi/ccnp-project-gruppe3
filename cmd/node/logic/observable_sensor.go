package logic

import (
	"fmt"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/entity"
)

type SensorObserver[T any] interface {
	UpdateSensorValue(sensor entity.Sensor[T])
}

type observableSensor[T any] struct {
	sensor   *entity.Sensor[T]
	observer SensorObserver[T]
}

func NewObservableSensor[T any](sensor *entity.Sensor[T], observer SensorObserver[T]) (*observableSensor[T], error) {
	if sensor == nil {
		return nil, fmt.Errorf("sensor cannot be nil")
	}
	if observer == nil {
		return nil, fmt.Errorf("observer cannot be nil")
	}
	return &observableSensor[T]{sensor: sensor, observer: observer}, nil
}

// SetValue sets the value of the Sensor and notifies the observer.
func (s *observableSensor[T]) SetValue(newValue T) {
	s.sensor.Value = newValue
	if s.observer != nil {
		s.observer.UpdateSensorValue(*s.sensor)
	}
}
