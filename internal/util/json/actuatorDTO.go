package json

import "github.com/Elias-Chairi/ccnp-project-gruppe3/internal/entity"

type ActuatorDTO struct {
	ID    uint8  `json:"id"`
	Type  string `json:"type"`
	Unit  string `json:"unit"`
	State any    `json:"state"`
}

func ToActuator(dto ActuatorDTO) entity.Actuator[any] {
	var state any

	switch dto.Type {
	case "HEATER": // heaters expect int32 or bool
		switch v := dto.State.(type) {
		case float64:
			state = any(int32(v))
		case bool:
			state = any(v)
		}

	case "WINDOW": // windows expect float32
		switch v := dto.State.(type) {
		case float64:
			state = any(float32(v))
		}

	case "FAN": // fans expect int32 or bool
		switch v := dto.State.(type) {
		case float64:
			state = any(int32(v))
		case bool:
			state = any(v)
		}

	case "LIGHT": // lights expect int32 or bool
		switch v := dto.State.(type) {
		case float64:
			state = any(int32(v))
		case bool:
			state = any(v)
		}

	default:
		// fallback, store as-is
		state = dto.State
	}

	return entity.Actuator[any]{
		ID:    dto.ID,
		Type:  dto.Type,
		Unit:  dto.Unit,
		State: state,
	}
}
