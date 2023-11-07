package main

import (
	"fmt"
	"os"
	. "rtclbedit/curp"
	. "rtclbedit/shared"
	"strconv"
)

func main() {
	args := os.Args
	if len(os.Args) != 4 {
		fmt.Println("[usage]: " + args[0] + " <identifier> <topo_file> <num_nodes>")
		os.Exit(1)
	}
	identifier := args[1]
	topo_file := args[2]
	num_nodes, err := strconv.Atoi(args[3])

	if err != nil {
		fmt.Println("Must have an integer node count!")
		os.Exit(1)
	}

	node_map := Parse(topo_file, num_nodes)
	master_node, peer_map, witness_map := ParseByRole(node_map)
	curr_node := node_map[identifier]

	if master_node == nil {
		fmt.Println("Invalid topo file (needs at least one master node)!")
		os.Exit(1)
	}

	if curr_node == nil {
		fmt.Println("Invalid topology (identifier not found in topo file)!")
		os.Exit(1)
	}

	DPrintf("%s: pre-init", identifier)
	InitRPC(identifier, node_map)
	DPrintf("%s: post-init", identifier)

	DPrintf("%s: pre-switch", identifier)
	channel := make(chan ExecuteMsg)
	DPrintf("%s: post-channel", identifier)
	switch curr_node.Role {
	case ROLE_MASTER:
		InitCurp(identifier, peer_map, witness_map, channel)
	case ROLE_BACKUP:
		InitCurp(identifier, peer_map, witness_map, channel)
	case ROLE_WITNESS:
		InitWitness(identifier, master_node)
	default:
		panic("Unknown Role! Exiting...")
	}
	DPrintf("%s: post-switch", identifier)

	// Busy-wait forever
	for {
	}
}
