package main

import (
	"bufio"
	"fmt"
	"log"
	"net/rpc"
	"os"
	. "rtclbedit/curp"
	. "rtclbedit/shared"
	"strconv"
	"sync"
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
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		cmd := scanner.Text()
		switch cmd {
		case "INSERT", "DELETE":
			args := ExecuteArgs{
				Command: cmd,
			}
			reply := ExecuteReply{}
			master_client.Call("Curp.Execute", args, &reply)
			log.Printf("Client %s finish executed", name)
			var wg sync.WaitGroup
			for witness_name, _ := range witness_clients {
				wg.Add(1)
				go func(name string) {
					defer wg.Done()
					args := RecordArgs{}
					reply := RecordReply{}
					witness_clients[name].Call("Witness.Record", args, &reply)
				}(witness_name)
			}
			wg.Wait()
			log.Printf("Client %s made command durable", name)
		case "DISPLAY":
			args := SyncArgs{}
			reply := SyncReply{}
			master_client.Call("Curp.Sync", args, &reply)
		}

	}
	log.Printf("Done editing session")
}
