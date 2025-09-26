package main

import "github.com/Elias-Chairi/ccnp-project-gruppe3/entity"

func main() {
	node := entity.NewNode([]entity.Actuator[any]{
		*entity.NewActuator(0, "FAN", "RPM", 800),
		*entity.NewActuator(1, "FAN", "RPM", 0),
	})
	sensor1 := entity.NewSensor(0, "TEMPERATURE", "°C", 22.5)
	sensor2 := entity.NewSensor(1, "HUMIDITY", "%", 60)
	node.AddSensor(*sensor1)
	node.AddSensor(*sensor2)

}
