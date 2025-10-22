package main

import (
	"log"
	"time"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/cmd/node/greenhouse"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/entity"
)

func main() {
	node := entity.Node{
		ID: 1,
		Actuators: []*entity.Actuator[any]{
			{ID: 1, Type: "HEATER", State: false},                  // off
			{ID: 1, Type: "HEATER", Unit: "°C", State: float32(0)}, // 0 °C (no heating) (heater cannot make it colder)
			{ID: 2, Type: "WINDOW", State: false},                  // closed
			{ID: 2, Type: "WINDOW", State: float32(0.0)},           // 0% open
			{ID: 3, Type: "FAN", State: false},                     // off
			{ID: 4, Type: "FAN", Unit: "RPM", State: int32(0)},     // 0 rpm
			{ID: 5, Type: "LIGHT", State: false},                   // off
			{ID: 5, Type: "LIGHT", Unit: "lx", State: int32(0)},    // 0 lux
		},
		Sensors: []*entity.Sensor[any]{
			{ID: 1, Type: "TEMPERATURE", Unit: "°C"},
			{ID: 2, Type: "HUMIDITY", Unit: "%"},
			{ID: 3, Type: "LIGHT", Unit: "lx"},
		},
	}

	outdoorConditions := greenhouse.OutdoorConditions{
		Temperature: 10.0,
		Humidity:    70.0,
		LightLevel:  200.0,
	}

	onSensorUpdate := func(sensor *entity.Sensor[any]) {
		log.Printf("Sensor %d, type %s, updated: %v%s\n", sensor.ID, sensor.Type, sensor.Value, sensor.Unit)
	}

	g := greenhouse.NewGreenhouse(node, outdoorConditions, onSensorUpdate)

	ticker := time.NewTicker(time.Second)
	for range ticker.C {
		log.Println("Simulating greenhouse step...")
		g.SimulateStep()
	}
}
