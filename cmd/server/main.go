package main

import (
	"fmt"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/cmd/server/udpService"
)

func main() {
	fmt.Println("starting service")
	udpservice.StartUDPService()
}
