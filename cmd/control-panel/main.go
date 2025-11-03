package main

import (
	"fmt"
	"log"
	"runtime"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/util"
)

func main() {
	log.Println("Hello, Control Panel!")
	addr, err := util.FindServer()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("I got the address %v\n", addr)
	fmt.Println("Goroutines:", runtime.NumGoroutine())
}
