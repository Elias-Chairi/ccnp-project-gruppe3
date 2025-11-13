package main

import (
	"log"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/cmd/server/tcpservice"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/cmd/server/udpservice"
)

func main() {
	log.Println("starting UDP service...")
	go udpservice.StartUDPService()

	log.Println("starting TCP service...")
	tcpservice.StartTCPService()
}
