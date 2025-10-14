package main

import (
	"log"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/cmd/node/logic"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/entity"
)

func main() {
	node := logic.NewNode(logic.WithActuators(
		&entity.Actuator[any]{ID: 0, Type: "FAN", Unit: "RPM", State: 800},
		&entity.Actuator[any]{ID: 1, Type: "LIGHT", Unit: "on/off", State: true},
	), logic.WithSensors(
		&entity.Sensor[any]{ID: 0, Type: "TEMPERATURE", Unit: "°C", Value: 22.5},
		&entity.Sensor[any]{ID: 1, Type: "HUMIDITY", Unit: "%", Value: 60},
	))
	log.Println(node.ID)
}
