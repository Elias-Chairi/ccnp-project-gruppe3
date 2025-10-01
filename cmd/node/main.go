package main

import (
	"log"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/entity"
)

func main() {
	node := entity.NewNode(entity.WithActuators(
		entity.NewActuator[any](0, "FAN", "RPM", 800),
		entity.NewActuator[any](1, "FAN", "RPM", 0),
	), entity.WithSensors(
		entity.NewSensor[any](0, "TEMPERATURE", "°C", 22.5),
		entity.NewSensor[any](1, "HUMIDITY", "%", 60),
	))
	log.Println(node.GetID())
}
