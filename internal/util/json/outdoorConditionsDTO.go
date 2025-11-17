package json

import "github.com/Elias-Chairi/ccnp-project-gruppe3/cmd/node/greenhouse"

type OutdoorConditionsDTO struct {
	Temperature any `json:"temperature"`
	Humidity    any `json:"humidity"`
	LightLevel  any `json:"lightLevel"`
}

func ToOutdoorConditions(dto OutdoorConditionsDTO) greenhouse.OutdoorConditions {
	oc := greenhouse.OutdoorConditions{}

	if v, ok := dto.Temperature.(float64); ok {
		oc.Temperature = float32(v)
	}
	if v, ok := dto.Humidity.(float64); ok {
		oc.Humidity = float32(v)
	}
	if v, ok := dto.LightLevel.(float64); ok {
		oc.LightLevel = int32(v)
	}

	return oc
}
