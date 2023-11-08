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

	if len(os.Args) != 4 {
		fmt.Println("[usage]: " + args[0] + " <identifier> <configuration file> <num_nodes>")
		os.Exit(1)
	}

	name := args[1]
	filename := args[2]
	node_ct, err := strconv.Atoi(args[3])

	if err != nil {
		fmt.Println("Must have an integer node count!")
		os.Exit(1)
	}

	node_map := Parse(filename, node_ct)
	master_node, _, witness_map := ParseByRole(node_map)
	log.Println(master_node)
	master_client, _ := rpc.Dial("tcp", master_node.Ip+":"+master_node.Port)
	log.Println(master_client)
	// var witness_clients []*rpc.Client

	witness_clients := make(map[string]*rpc.Client)

	fmt.Println(witness_map)
	log.Printf("Client %s attempting to connect...", name)
	for k, v := range witness_map {
		c, _ := rpc.Dial("tcp", v.Ip+":"+v.Port)
		// witness_clients = append(witness_clients, c)
		witness_clients[k] = c
	}
	log.Printf("Client %s connected!", name)

	// Example of sending 1 request
	log.Printf("Sending request...")
	{
		args := ExecuteArgs{}
		reply := ExecuteReply{}
		master_client.Call("Curp.Execute", args, &reply)
		fmt.Println("RECEIVED EXECUTE REPLY")
		fmt.Println(reply)
		if reply.Status == FAILURE_STATUS {
			fmt.Println("ERROR WITH CALLING MASTER!")
		}

		for name, witness_client := range witness_clients {
			args := RecordArgs{Name: name}
			reply := RecordReply{}
			// TODO: CONVERT THIS TO GOROUTINES
			witness_client.Call("Witness.Record", args, &reply)
			if reply.Status == FAILURE_STATUS {
				fmt.Printf("ERROR! Unable to record to node %s", args.Name)
			}
		}
	}
	log.Printf("Finished operations!")

	// [TODO] error handling

}
