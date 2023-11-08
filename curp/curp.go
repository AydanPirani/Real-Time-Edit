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

var HAS_NOT_VOTED string = "HAS_NOT_VOTED"

/**
 * As each Curp peer becomes aware that successive log entries are
 * committed, the peer should send an ExecuteMsg to the application
 * sit on top of it, via the executeCh passed to Make(). Set
 * CommandValid to true to indicate that the ExecuteMsg contains a newly
 * committed log entry.
 */

// A Go object implementing a single Curp peer.
type Curp struct {
	mu              sync.Mutex
	name            string // Curp node's unique identifier
	witness_clients map[string]*rpc.Client
	peer_clients    map[string]*rpc.Client
	appChan         chan ExecuteMsg
	/**
	 * Raft component member variables
	 */
	timeoutChan chan struct{} // send message to this channel when we are resetting the timer
	currentTerm int
	votedFor    string
	syncedIndex int

	log        []LogEntry
	nextIndex  map[string]int
	matchIndex map[string]int
	role       NodeRole

	dead         int32 // set by Kill
	indexCounter int   // to specify random ranges before starting new election
}

func (cr *Curp) sendOrderAsync(peer_name string, heartbeat bool) {
	cr.mu.Lock()
	if cr.role != ROLE_MASTER {
		cr.mu.Unlock()
		return
	}
	prevIndex := cr.nextIndex[peer_name] - 1

	var prevTerm int
	if prevIndex == -1 {
		prevTerm = 0
	} else {
		prevTerm = cr.log[prevIndex].Term
	}
	args := &OrderAsyncArgs{
		Term:         cr.currentTerm,
		LeaderName:   cr.name,
		PrevLogIndex: prevIndex,
		PrevLogTerm:  prevTerm,
		Entries:      nil,
		LeaderSynced: cr.syncedIndex,
	}
	if !heartbeat {
		args.Entries = cr.log[prevIndex+1:]
	}

	// DPrintf("Leader %s sending message %+v to server %d\n", rf.me, args, server)
	cr.mu.Unlock()
	var reply OrderAsyncReply
	good := cr.peer_clients[peer_name].Call("Curp.OrderAsync", args, &reply)
	cr.mu.Lock()
	if good == nil {
		if reply.Success {
			// DPrintf("Got success %+v\n", reply)
			if prevIndex+len(args.Entries) >= cr.nextIndex[peer_name] {
				cr.nextIndex[peer_name] = prevIndex + len(args.Entries) + 1
				cr.matchIndex[peer_name] = prevIndex + len(args.Entries)

				DPrintf("Server %s updated server %s to %d\n", cr.name, peer_name, cr.nextIndex[peer_name])
			}

			lastEntry := prevIndex + len(args.Entries)
			if lastEntry < len(cr.log) && cr.syncedIndex <= lastEntry && cr.log[lastEntry].Term == cr.currentTerm {
				count := 1
				for name := range cr.peer_clients {
					if name != cr.name && cr.matchIndex[name] >= lastEntry {
						count++
					}
				}
				if count > len(cr.peer_clients)/2 {
					for lastEntry+1 > cr.syncedIndex {
						cr.syncedIndex++
					}
					DPrintf("Leader %s updated synced index to %d\n", cr.name, cr.syncedIndex)
				}
			}
		} else {
			if reply.Term > cr.currentTerm {
				//stepping down because there is a new leader currentTerm is bigger
				cr.currentTerm = reply.Term
				cr.votedFor = HAS_NOT_VOTED
				cr.role = ROLE_BACKUP
			} else {
				// log.Printf("LEADER updating %s nextIndex to %d\n", server, reply.LastMatch+1)
				// rf.nextIndex[server] = reply.LastMatch + 1
				cr.nextIndex[peer_name] -= 2
				if cr.nextIndex[peer_name] < 0 {
					cr.nextIndex[peer_name] = 0
				}
				// go rf.sendAppendEntries(server, false)
			}
		}
	}
	cr.mu.Unlock()
}

func (cr *Curp) SendHeartbeat() {
	for name := range cr.peer_clients {
		if name == cr.name {
			continue
		}
		DPrintf("leader %s sending heartbeat to %s", cr.name, name)
		go cr.sendOrderAsync(name, true)
	}
}

/**
 * Redesign interaction between Start, SendHeartBeat, SendOrderAsync
 * Will probably look something like. send in parallels, and process request
 * coming into / out of a replyCh
 */
func (cr *Curp) Start(command interface{}) (int, int, bool) {
	// Your code here (2B).
	cr.mu.Lock()
	term := cr.currentTerm

	if cr.role != ROLE_MASTER {
		cr.mu.Unlock()
		return -1, term, false
	}
	// DPrintf("Leader %d got command %+v\n", rf.me, command)
	// send append entries to everyone
	cr.log = append(cr.log, LogEntry{Command: command, Term: term})
	len := len(cr.log)
	cr.mu.Unlock()

	for name := range cr.peer_clients {
		if name != cr.name {
			go cr.sendOrderAsync(name, false)
		}
	}

	return len, term, true
}
