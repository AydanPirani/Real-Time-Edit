package main

import (
	"fmt"
	"log"
	"net/rpc"
	"os"
	. "rtclbedit/curp"
	. "rtclbedit/shared"
	"strconv"
)

func main() {
	args := os.Args

	if len(os.Args) != 5 {
		fmt.Println("[usage]: " + args[0] + " <identifier> <configuration file> <num_nodes> <command>")
		os.Exit(1)
	}

	name := args[1]
	filename := args[2]
	node_ct, err := strconv.Atoi(args[3])
	cmd := args[4]

	if err != nil {
		fmt.Println("Must have an integer node count!")
		os.Exit(1)
	}

	node_map := Parse(filename, node_ct)
	master_node, _, witness_map := ParseByRole(node_map)

	// var witness_clients []*rpc.Client

	witness_clients := make(map[string]*rpc.Client)

	fmt.Println(witness_map)
	log.Printf("Client %s attempting to connect...", name)
	master_client, _ := rpc.Dial("tcp", master_node.Ip+":"+master_node.Port)
	for k, v := range witness_map {
		c, _ := rpc.Dial("tcp", v.Ip+":"+v.Port)
		// witness_clients = append(witness_clients, c)
		witness_clients[k] = c
	}
	log.Printf("Client %s connected Curp system", name)
	// Example of sending 1 request
	{
		args := ExecuteArgs{
			Command: cmd,
		}
		reply := ExecuteReply{}
		master_client.Call("Curp.Execute", args, &reply)
		log.Printf("Client %s finish executed", name)
		for _, witness_client := range witness_clients {
			args := RecordArgs{}
			reply := RecordReply{}
			// TODO: CONVERT THIS TO GOROUTINES
			witness_client.Call("Witness.Record", args, &reply)
		}
		log.Printf("Client %s made command durable", name)
	}

}
