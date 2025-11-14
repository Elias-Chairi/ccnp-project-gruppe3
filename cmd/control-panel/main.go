package main

import (
	"flag"

	// "fmt"
	"log"
	"time"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/cmd/control-panel/connectionhandler"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/cmd/control-panel/view"

	// "github.com/Elias-Chairi/ccnp-project-gruppe3/cmd/control-panel/view"
	// "github.com/Elias-Chairi/ccnp-project-gruppe3/internal/entity"
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

	t := &view.TerminalView{}

	c := &connectionhandler.ConnectionHandler{}
	c.ServerIPs = ipList

	go func() {
		time.Sleep(time.Second * 3) // simulate loading time
		err := c.Connect()
		if err != nil {
			t.FailedToConnectToServer(c.ServerIPs[0])
			return
		}

		nodes, err := c.Register()
		if err != nil {
			t.FailedToRegisterToServer(c.ServerIPs[0])
			return
		}
		t.SetInitialNodes(*nodes)
		t.EndLoading()
	}()

	t.Start()
}
