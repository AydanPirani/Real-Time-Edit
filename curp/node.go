package curp

import (
	"log"
	"net"
	"net/rpc"
	. "rtclbedit/shared"
	"time"

	mapset "github.com/deckarep/golang-set/v2"
)

func Connect(name string, node_map map[string]*Node) map[string]*rpc.Client {
	connection_map := make(map[string]*rpc.Client)
	for _, v := range node_map {
		for {
			addr := v.Ip + ":" + v.Port
			conn, err := rpc.Dial("tcp", addr)
			if err != nil {
				// log.Println("Failed to connect to " + v.name)
				time.Sleep(150 * time.Millisecond)
			} else {

				connection_map[v.Name] = conn
				break
			}
		}

	}
	return connection_map
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
	c := &Curp{
		name:            name,
		witness_clients: Connect(name, witness_map),
		peers_clients:   Connect(name, peer_map),
		appChan:         appChan,
		timeoutChan:     make(chan struct{}, 1),
		currentTerm:     0,
		votedFor:        -1,
		syncedIndex:     0,
		nextIndex:       make([]int, len(peer_map)),
		matchIndex:      make([]int, len(peer_map)),
		role:            ROLE_BACKUP,
	}

	rpc.Register(c)
	return c
}

func InitWitness(name string, master_map map[string]*Node) {
	if len(master_map) != 1 {
		log.Fatalf("Backup %s has more than one master!", name)
	}

	witness := &Witness{
		name:          name,
		unsynced:      mapset.NewSet[string](),
		master_client: Connect(name, master_map),
	}

	rpc.Register(witness)
}
