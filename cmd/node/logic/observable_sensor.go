package logic

import (
	"fmt"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/entity"
)

// SensorObserver defines the interface for observing sensor value changes.
type SensorObserver[T any] interface {
	UpdateSensorValue(sensor entity.Sensor[T])
}

// observableSensor wraps an entity.Sensor and notifies its observer on value changes.
type observableSensor[T any] struct {
	sensor   *entity.Sensor[T]
	observer SensorObserver[T]
}

// newObservableSensor creates a new observableSensor instance.
func newObservableSensor[T any](sensor *entity.Sensor[T], observer SensorObserver[T]) (*observableSensor[T], error) {
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
