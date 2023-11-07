package main

import (
	"fmt"
	"os"
	. "rtclbedit/shared"
)

func main() {
	args := os.Args
	if len(os.Args) != 3 {
		fmt.Println("[usage]: " + args[0] + " <identifier> <configuration file>")
		os.Exit(1)
	}
	name := args[1]
	filename := args[2]

	node_map := Parse(filename)
	peer_map, witness_map := ParseByRole(node_map)

	for {
	}
}
