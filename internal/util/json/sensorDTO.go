package json

import "github.com/Elias-Chairi/ccnp-project-gruppe3/internal/entity"

type SensorDTO struct {
	ID    uint8  `json:"id"`
	Type  string `json:"type"`
	Unit  string `json:"unit"`
	Value any    `json:"value"`
}

func ToSensor(dto SensorDTO) entity.Sensor[any] {
	var value any

	switch v := dto.Value.(type) {
	case float64:
		value = any(float32(v))
	default:
		value = dto.Value
	}

	return entity.Sensor[any]{
		ID:    dto.ID,
		Type:  dto.Type,
		Unit:  dto.Unit,
		Value: value,
	}
}
