package main

import (
	"flag"
	"log"
	"time"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/cmd/node/connectionhandler"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/cmd/node/greenhouse"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/entity"
	util "github.com/Elias-Chairi/ccnp-project-gruppe3/internal/util/general"
	utilNet "github.com/Elias-Chairi/ccnp-project-gruppe3/internal/util/net"
)

func main() {
	node := entity.Node{
		ID: 1,
		Actuators: []entity.Actuator[any]{
			{ID: 1, Type: "HEATER", State: false},                        // off
			{ID: 2, Type: "HEATER", Unit: "°C", State: int32(0)},         // 0 °C (no heating) (heater cannot make it colder)
			{ID: 4, Type: "WINDOW", Unit: "% open", State: float32(0.0)}, // 0% open
			{ID: 5, Type: "FAN", State: false},                           // off
			{ID: 6, Type: "FAN", Unit: "RPM", State: int32(0)},           // 0 rpm
			{ID: 7, Type: "LIGHT", State: false},                         // off
			{ID: 8, Type: "LIGHT", Unit: "lx", State: int32(500)},        // 0 lux
		},
		Sensors: []entity.Sensor[any]{
			{ID: 1, Type: "TEMPERATURE", Unit: "°C", Value: float32(20.0)},
			{ID: 2, Type: "HUMIDITY", Unit: "%", Value: float32(50.0)},
			{ID: 3, Type: "LIGHT", Unit: "lx", Value: float32(100.0)},
		},
	}

	outdoorConditions := greenhouse.OutdoorConditions{
		Temperature: 10.0,
		Humidity:    70.0,
		LightLevel:  200.0,
	}

	var ipList util.IPList
	flag.Var(&ipList, "server", "A comma-separated list of server IP addresses")
	flag.Parse()

	if len(ipList) == 0 {
		log.Println("Searching for server... (search duration: 3s)")
		var err error
		ipList, err = utilNet.FindServer(3 * time.Second)
		if err != nil {
			log.Fatalf("Error searching for servers: %v\n", err)
		}
		if len(ipList) == 0 {
			log.Fatalln("No servers found, please specify server IPs manually using the -server flag.")
		}
	}

	c := connectionhandler.ConnectionHandler{
		ServerIPs: ipList,
	}

	log.Println("Connecting to server...")
	if err := c.Connect(); err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	log.Println("Registering node...")
	assignedID, err := c.Register(node)
	if err != nil {
		log.Fatalf("Failed to register: %v", err)
	}

	node.ID = assignedID
	log.Printf("Node assigned ID %d\n", node.ID)

	onSensorUpdate := func(sensor *entity.Sensor[any]) {
		log.Printf("Sensor %d, type %s, updated: %v%s\n", sensor.ID, sensor.Type, sensor.Value, sensor.Unit)
		if err := c.SendSensorUpdate(*sensor); err != nil {
			log.Fatalf("Failed to send sensor update: %v\n", err)
		}
	}

	g := greenhouse.NewGreenhouse(node, outdoorConditions, onSensorUpdate)

	ticker := time.NewTicker(time.Second)
	for range ticker.C {
		log.Println("Simulating greenhouse step...")
		g.SimulateStep()

	}
}
