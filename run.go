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
	fmt.Println(name)
	filename := args[2]

	node_map := Parse(filename)
	// master_node, peer_map, witness_map := ParseByRole(node_map)
	master_node, _, _ := ParseByRole(node_map)

	fmt.Println(master_node)
	for {

	}
}
