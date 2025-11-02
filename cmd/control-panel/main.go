package main

import (
    "github.com/Elias-Chairi/ccnp-project-gruppe3/internal/entity"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/cmd/control-panel/view"
)

func main() {
	nodeA := &entity.Node{
		ID: 1,
		Sensors: []*entity.Sensor[any]{
			{ID: 1, Type: "Temperature", Unit: "°C", Value: 22.3},
			{ID: 2, Type: "Humidity", Unit: "%", Value: 55},
			{ID: 3, Type: "CO2", Unit: "ppm", Value: 420},
		},
		Actuators: []*entity.Actuator[any]{
			{ID: 1, Type: "Heater", State: "OFF"},
			{ID: 2, Type: "Fan", State: "ON"},
			{ID: 3, Type: "Sprinkler", State: "OFF"},
		},
	}

	// --- Greenhouse B ---
	nodeB := &entity.Node{
		ID: 2,
		Sensors: []*entity.Sensor[any]{
			{ID: 1, Type: "Temperature", Unit: "°C", Value: 25.7},
			{ID: 2, Type: "Humidity", Unit: "%", Value: 48},
			{ID: 3, Type: "CO2", Unit: "ppm", Value: 390},
		},
		Actuators: []*entity.Actuator[any]{
			{ID: 1, Type: "Heater", Unit: "", State: "OFF"},
			{ID: 2, Type: "Fan", Unit: "", State: "OFF"},
			{ID: 3, Type: "Sprinkler", Unit: "", State: "ON"},
		},
	}

	// --- Greenhouse C ---
	nodeC := &entity.Node{
		ID: 3,
		Sensors: []*entity.Sensor[any]{
			{ID: 1, Type: "Temperature", Unit: "°C", Value: 19.5},
			{ID: 2, Type: "Humidity", Unit: "%", Value: 62},
			{ID: 3, Type: "CO2", Unit: "ppm", Value: 450},
		},
		Actuators: []*entity.Actuator[any]{
			{ID: 1, Type: "Heater", Unit: "", State: "ON"},
			{ID: 2, Type: "Fan", Unit: "", State: "ON"},
			{ID: 3, Type: "Sprinkler", Unit: "", State: "OFF"},
		},
	}

	// --- Combine all nodes ---
	nodes := []*entity.Node{nodeA, nodeB, nodeC}
    v := view.NewTerminalWithData(&nodes)
    v.Start()
}
