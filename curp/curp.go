/**
 * Our Curp implementation act as a frontend to Raft.
 * Curp only introduce a fast-path for commutative operations
 * to be done in 1 RTT. Ordering and leader changes happens
 * in the background and follows Raft logic.
 */
package curp

import (
	"net/rpc"
	. "rtclbedit/shared"
	"sync"
)

/**
 * As each Curp peer becomes aware that successive log entries are
 * committed, the peer should send an ExecuteMsg to the application
 * sit on top of it, via the executeCh passed to Make(). set
 * CommandValid to true to indicate that the ApplyMsg contains a newly
 * committed log entry.
 */

// A Go object implementing a single Curp peer.
type Curp struct {
	mu              sync.Mutex
	name            string // Curp node's unique identifier
	witness_clients map[string]*rpc.Client
	peers_clients   map[string]*rpc.Client
	appChan         chan ExecuteMsg
	/**
	 * Raft component member variables
	 */
	timeoutChan chan struct{} // send message to this channel when we are resetting the timer
	currentTerm int
	votedFor    int
	syncedIndex int

	log        []LogEntry
	nextIndex  []int
	matchIndex []int
	role       NodeRole
}

/**
 * RPC types
 */
type ExecuteArgs struct {
}
type ExecuteReply struct {
}
type SyncArgs struct {
}
type SyncReply struct {
}

/**
 * RPC functions
 */
func (c *Curp) Execute(args ExecuteArgs, reply *ExecuteReply) error { // executeRPC called by clients to master
	return nil
}

func (c *Curp) Sync(args SyncArgs, reply *SyncReply) error { // syncRPC called by clients to master
	return nil
}

func (c *Curp) sendOrderAsync(name string, heartbeat bool) {

}
