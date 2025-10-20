package main

import "github.com/Elias-Chairi/ccnp-project-gruppe3/cmd/control-panel/view"

func main() {
	v := view.NewTerminal()
	v.Start()
}
