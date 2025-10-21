package main

import (
    "github.com/Elias-Chairi/ccnp-project-gruppe3/internal/entity"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/cmd/control-panel/view"
)

func main() {
    dummyNodes := []*entity.Node{
        {
            ID: 1,
            Sensors: []*entity.Sensor[any]{
                {ID: 1, Type: "Temperature", Unit: "°C", Value: 24.7},
                {ID: 2, Type: "Humidity", Unit: "%", Value: 61.3},
                {ID: 3, Type: "CO₂", Unit: "ppm", Value: 420},
            },
            Actuators: []*entity.Actuator[any]{
                {ID: 1, Type: "Fan", Unit: "On/Off", State: "ON"},
                {ID: 2, Type: "Heater", Unit: "On/Off", State: "OFF"},
            },
        },
    }

    v := view.NewTerminalWithData(dummyNodes)
    v.Start()
}
