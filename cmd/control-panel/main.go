package main

import (
	"fmt"
	"log"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/util"
)

func main() {
	log.Println("Hello, Control Panel!")
	addr, err := util.FindServer()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("I got the address %v", addr)
}
