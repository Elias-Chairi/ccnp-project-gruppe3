package entity

type SensorObserver[T any] interface {
	UpdateSensorValue(sensor *Sensor[T])
}

type Sensor[T any] struct {
	id         int
	sensorType string
	unit       string
	value      T
	observer   SensorObserver[T]
}

// NewSensor creates a new Sensor instance with the given parameters and an observer.
func NewSensor[T any](id int, sensorType string, unit string, initialValue T) *Sensor[any] {
	return &Sensor[any]{
		id:         id,
		sensorType: sensorType,
		unit:       unit,
		value:      initialValue,
	}
}

// GetID returns the ID of the Sensor.
func (s *Sensor[T]) GetID() int {
	return s.id
}

// SetID sets the ID of the Sensor.
func (s *Sensor[T]) SetID(id int) {
	s.id = id
}

// GetType returns the type of the Sensor.
func (s *Sensor[T]) GetType() string {
	return s.sensorType
}

// GetUnit returns the unit of the Sensor.
func (s *Sensor[T]) GetUnit() string {
	return s.unit
}

// SetValue sets the value of the Sensor and notifies the observer.
func (s *Sensor[T]) SetValue(newValue T) {
	s.value = newValue
	if s.observer != nil {
		s.observer.UpdateSensorValue(s)
	}
}

// GetValue returns the current value of the Sensor.
func (s *Sensor[T]) GetValue() T {
	return s.value
}

// SetObserver sets the observer for the Sensor.
func (s *Sensor[T]) SetObserver(observer SensorObserver[T]) {
	s.observer = observer
}
