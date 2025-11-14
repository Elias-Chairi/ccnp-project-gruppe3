package main

import (
	"log"
	"time"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/cmd/control-panel/connectionhandler"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/cmd/node/greenhouse"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/entity"
)

func main() {
	node := entity.Node{
		ID: 1,
		Actuators: []entity.Actuator[any]{
			{ID: 1, Type: "HEATER", State: false},                  // off
			{ID: 2, Type: "HEATER", Unit: "°C", State: float32(0)}, // 0 °C (no heating) (heater cannot make it colder)
			{ID: 3, Type: "WINDOW", State: false},                  // closed
			{ID: 4, Type: "WINDOW", State: float32(0.0)},           // 0% open
			{ID: 5, Type: "FAN", State: false},                     // off
			{ID: 6, Type: "FAN", Unit: "RPM", State: int32(0)},     // 0 rpm
			{ID: 7, Type: "LIGHT", State: false},                   // off
			{ID: 8, Type: "LIGHT", Unit: "lx", State: int32(0)},    // 0 lux
		},
		Sensors: []entity.Sensor[any]{
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

	c := connectionhandler.ConnectionHandler{}

	log.Println("Connecting to server...")
	c.Connect()

	log.Println("Registering node...")
	c.Register()

	onSensorUpdate := func(sensor *entity.Sensor[any]) {
		log.Printf("Sensor %d, type %s, updated: %v%s\n", sensor.ID, sensor.Type, sensor.Value, sensor.Unit)
		// c.SendSensorUpdate(node.ID, *sensor)
	}

	g := greenhouse.NewGreenhouse(node, outdoorConditions, onSensorUpdate)

	ticker := time.NewTicker(time.Second)
	for range ticker.C {
		log.Println("Simulating greenhouse step...")
		g.SimulateStep()

	}
}
