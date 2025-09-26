package main

import "github.com/Elias-Chairi/ccnp-project-gruppe3/entity"

func main() {
	node := entity.NewNode([]entity.Actuator[any]{
		*entity.NewActuator[any](0, "FAN", "RPM", 800),
		*entity.NewActuator[any](1, "FAN", "RPM", 0),
	})
	sensor1 := entity.NewSensor[any](0, "TEMPERATURE", "°C", 22.5)
	sensor2 := entity.NewSensor[any](1, "HUMIDITY", "%", 60)
	node.AddSensor(*sensor1)
	node.AddSensor(*sensor2)

}
