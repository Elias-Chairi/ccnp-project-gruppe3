package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/util"
)

func main() {
	var ipList util.IPList
	flag.Var(&ipList, "server", "A comma-separated list of server IP addresses")
	flag.Parse()

	if len(ipList) == 0 {
		log.Println("Searching for server... (search duration: 3s)")
		var err error
		ipList, err = util.FindServer(3 * time.Second)
		if err != nil {
			log.Fatalf("Error searching for servers: %v\n", err)
		}
		if len(ipList) == 0 {
			log.Fatalln("No servers found, please specify server IPs manually using the -server flag.")
		}
	}

	fmt.Printf("Parsed server IP addresses: %v\n", ipList)
}
