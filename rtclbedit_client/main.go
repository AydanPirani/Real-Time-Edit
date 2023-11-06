package main

import (
	"fmt"
	"log"
	"net/rpc"
	"os"
	. "rtclbedit/curp"
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
	master_map, witness_map, _ := ParseByRole(node_map)
	master := Map_begin(master_map)
	master_client, _ := rpc.Dial("tcp", master.Ip+":"+master.Port)
	var witness_clients []*rpc.Client
	for _, v := range witness_map {
		c, _ := rpc.Dial("tcp", v.Ip+":"+v.Port)
		witness_clients = append(witness_clients, c)
	}
	// sending 1 request
	{
		log.Printf("Client %s making requests...", name)
		args := ExecuteArgs{}
		reply := ExecuteReply{}
		master_client.Call("Master.Execute", args, reply)
		for _, witness_client := range witness_clients {
			go func(witness_client *rpc.Client) {
				args := RecordArgs{}
				reply := RecordReply{}
				witness_client.Call("Witness.Record", args, reply)
			}(witness_client)
		}
	}

	// [TODO] error handling

}
