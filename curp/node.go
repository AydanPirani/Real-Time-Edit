package curp

import (
	"log"
	"net"
	"net/rpc"
	. "rtclbedit/shared"
	"time"
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

func InitMaster(name string, witness_map map[string]*Node, backup_map map[string]*Node) {
	witness_send_map := Connect(name, witness_map)
	backup_send_map := Connect(name, backup_map)

	master := &Master{
		name:            name,
		unsynced:        []Operation{},
		witness_clients: witness_send_map,
		backup_clients:  backup_send_map,
	}

	rpc.Register(master)
}

func InitBackup(name string, master_map map[string]*Node) {
	if len(master_map) != 1 {
		log.Fatalf("Backup %s has more than one master!", name)
	}

	master_send_map := Connect(name, master_map)
	backup := &Backup{
		name:          name,
		master_client: master_send_map,
	}

	rpc.Register(backup)
}

func InitWitness(name string, master_map map[string]*Node) {
	if len(master_map) != 1 {
		log.Fatalf("Backup %s has more than one master!", name)
	}

	master_send_map := Connect(name, master_map)
	witness := &Witness{
		name:          name,
		unsynced:      []Operation{},
		master_client: master_send_map,
	}

	rpc.Register(witness)
}
