package main

import (
	"fmt"
	"os"
	. "rtclbedit/shared"
	"strconv"
)

func main() {
	args := os.Args
	if len(os.Args) != 4 {
		fmt.Println("[usage]: " + args[0] + " <identifier> <configuration file> <num_nodes>")
		os.Exit(1)
	}
	name := args[1]
	fmt.Println(name)
	filename := args[2]
	node_ct, err := strconv.Atoi(args[3])

	if err != nil {
		fmt.Println("Must have an integer node count!")
		os.Exit(1)
	}

	node_map := Parse(filename, node_ct)
	// master_node, peer_map, witness_map := ParseByRole(node_map)
	master_node, _, _ := ParseByRole(node_map)

	if master_node == nil {
		fmt.Println("Invalid topo file (needs at least one master node)!")
		os.Exit(1)
	}

	fmt.Println(master_node)
	for {

	}
}
