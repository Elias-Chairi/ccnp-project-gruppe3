package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/cmd/node/connectionhandler"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/cmd/node/greenhouse"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/entity"
	util "github.com/Elias-Chairi/ccnp-project-gruppe3/internal/util/general"
	utilJson "github.com/Elias-Chairi/ccnp-project-gruppe3/internal/util/json"
	utilNet "github.com/Elias-Chairi/ccnp-project-gruppe3/internal/util/net"
)

func main() {
	greenhouseConfigPath := flag.String("config", "", "Path to the greenhouse JSON configuration file")

	var ipList util.IPList
	flag.Var(&ipList, "server", "A comma-separated list of server IP addresses")

	flag.Parse()

	if *greenhouseConfigPath == "" {
		fmt.Println("Please provide a greenhouse configuration file path using the -config flag e.g., -config=examples/greenhouse-configs/greenhouse1.json")
		flag.PrintDefaults() // Prints usage information
		os.Exit(1)
	}

	log.Println("Loading greenhouse configuration...")
	configData, err := os.ReadFile(*greenhouseConfigPath)
	if err != nil {
		log.Fatalf("Failed to read greenhouse configuration file: %v", err)
	}

	node, outdoorConditions, err := utilJson.ParseGreenhouseJSON(string(configData))
	if err != nil {
		log.Fatalf("Failed to parse greenhouse configuration: %v", err)
	}

	if len(ipList) == 0 {
		log.Println("Searching for server... (search duration: 3s)")
		var err error
		ipList, err = utilNet.FindServer(3 * time.Second)
		if err != nil {
			log.Fatalf("Error searching for servers: %v\n", err)
		}
		if len(ipList) == 0 {
			log.Println("No servers found, please specify server IPs manually using the -server flag.")
			flag.PrintDefaults() // Prints usage information
			os.Exit(1)
		}
	}

	c := connectionhandler.ConnectionHandler{
		ServerIPs: ipList,
		Node:      &node,
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

	g := greenhouse.NewGreenhouse(&node, outdoorConditions, onSensorUpdate)

	go c.StartListeningForCommands()

	ticker := time.NewTicker(time.Second)
	for range ticker.C {
		log.Println("Simulating greenhouse step...")
		g.SimulateStep()

	}
}
