package curp

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/rpc"
	. "rtclbedit/shared"
	"sync/atomic"
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
		peer_clients:    ConnectMultiple(peer_map),
		appChan:         appChan,
		timeoutChan:     make(chan struct{}, 1),
		currentTerm:     0,
		votedFor:        HAS_NOT_VOTED,
		syncedIndex:     0,
		nextIndex:       make(map[string]int),
		matchIndex:      make(map[string]int),
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

func (cr *Curp) NodeLifetime() {
	//TODO: sleep or semamore or conditional variable
	//Might redesign
	for !cr.killed() {
		cr.mu.Lock()
		role := cr.role
		cr.mu.Unlock()
		// DPrintf("Server %d %s\n", rf.me, role)
		switch role {
		case ROLE_BACKUP:
			{
				timer := time.Duration(rand.Intn(500-300+cr.indexCounter*50)+300) * time.Millisecond
				select {
				case <-cr.timeoutChan:
					{
						// reset timeout
						break
					}
				case <-time.After(timer):
					{
						// send request vote
						cr.mu.Lock()
						DPrintf("Election Timed out at server %d, switching to Candidate\n", cr.name)
						cr.role = ROLE_CANDIDATE
						cr.mu.Unlock()
						break
					}
				}
				break
			}
		case ROLE_CANDIDATE:
			{
				cr.StartElection()
				break
			}
		case ROLE_MASTER:
			{
				cr.SendHeartbeat()
				//TODO config variable for heartbeat frequency
				time.Sleep(time.Duration(150) * time.Millisecond)
				break
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func (cr *Curp) Kill() {
	atomic.StoreInt32(&cr.dead, 1)
	// Your code here, if desired.
}

func (cr *Curp) killed() bool {
	z := atomic.LoadInt32(&cr.dead)
	return z == 1
}
