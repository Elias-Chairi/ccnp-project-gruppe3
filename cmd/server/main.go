package main

import (
	"log"

	udpservice "github.com/Elias-Chairi/ccnp-project-gruppe3/cmd/server/udpService"
)

func main() {
	log.Println("starting service")
	udpservice.StartUDPService()
}
