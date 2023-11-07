package curp

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	. "rtclbedit/shared"
	"time"

	mapset "github.com/deckarep/golang-set/v2"
)

func ConnectMultiple(node_map map[string]*Node) map[string]*rpc.Client {
	connection_map := make(map[string]*rpc.Client)
	for _, node := range node_map {
		connection_map[node.Name] = Connect(node)
	}
	return connection_map
}

func Connect(node *Node) *rpc.Client {
	for {
		addr := node.Ip + ":" + node.Port
		conn, err := rpc.Dial("tcp", addr)
		if err != nil {
			fmt.Println(err)
			time.Sleep(150 * time.Millisecond)
		} else {
			return conn
		}
	}
}

func InitRPC(name string, node_map map[string]*Node) {
	port := node_map[name].Port
	port = ":" + port
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Println("failed to listen on port", port)
		return
	}

	go rpc.Accept(listener)
}

func InitCurp(name string, peer_map map[string]*Node, witness_map map[string]*Node, appChan chan ExecuteMsg) *Curp {
	DPrintf("%s: creating", name)
	c := &Curp{
		name:            name,
		witness_clients: ConnectMultiple(witness_map),
		peers_clients:   ConnectMultiple(peer_map),
		appChan:         appChan,
		timeoutChan:     make(chan struct{}, 1),
		currentTerm:     0,
		votedFor:        -1,
		syncedIndex:     0,
		nextIndex:       make([]int, len(peer_map)),
		matchIndex:      make([]int, len(peer_map)),
		role:            ROLE_BACKUP,
	}
	DPrintf("%s: pre reg", name)
	rpc.Register(c)
	DPrintf("%s: post reg", name)
	return c
}

func InitWitness(name string, master_node *Node) {
	if master_node == nil {
		log.Fatalf("Backup %s has more than one master!", name)
	}

	witness := &Witness{
		name:          name,
		unsynced:      mapset.NewSet[string](),
		master_client: Connect(master_node),
	}

	rpc.Register(witness)
}
