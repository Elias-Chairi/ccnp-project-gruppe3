package main

import (
	"log"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/cmd/node/logic"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/entity"
)

func main() {
	node := logic.NewNode(logic.WithActuators(
		&entity.Actuator[any]{0, "FAN", "RPM", 800},
		&entity.Actuator[any]{1, "FAN", "RPM", 0},
	), logic.WithSensors(
		&entity.Sensor[any]{0, "TEMPERATURE", "°C", 22.5},
		&entity.Sensor[any]{1, "HUMIDITY", "%", 60},
	))
	log.Println(node.ID)
}
