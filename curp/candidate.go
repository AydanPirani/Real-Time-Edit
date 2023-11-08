package curp

import (
	"math/rand"
	. "rtclbedit/shared"
	"time"
)

// example RequestVote RPC arguments structure.
// field names must start with capital letters!
type RequestVoteArgs struct {
	// Your data here (2A, 2B).
	Term          int
	CandidateName string
	LastLogIndex  int
	LastLogTerm   int
}

// example RequestVote RPC reply structure.
// field names must start with capital letters!
type RequestVoteReply struct {
	// Your data here (2A).
	Term        int
	VoteGranted bool
}

func (cr *Curp) StartElection() {
	// term < currentTerm
	cr.mu.Lock()
	cr.currentTerm++
	cr.votedFor = HAS_NOT_VOTED
	// send requestVotes
	// spawn a goroutine to process them
	electionResponse := make(chan RequestVoteReply)

	var lastTerm int
	lastLogIndex := len(cr.log) - 1
	if lastLogIndex == -1 {
		lastTerm = 0
	} else {
		lastTerm = cr.log[lastLogIndex].Term
	}
	args := &RequestVoteArgs{
		CandidateName: cr.name,
		Term:          cr.currentTerm,
		LastLogIndex:  lastLogIndex,
		LastLogTerm:   lastTerm,
	}
	cr.mu.Unlock()
	for name := range cr.peer_clients {
		go func(name string) {
			var reply RequestVoteReply
			if cr.sendRequestVote(name, args, &reply) != nil {
				electionResponse <- reply
			}
		}(name)
	}
	timer := time.Duration(rand.Intn(5000-1500)+1500) * time.Millisecond
	votes := 0
	electionDone := false
	n := len(cr.peer_clients)
	for !electionDone && !cr.killed() {
		select {
		case reply := <-electionResponse:
			{
				DPrintf("Server %s got response from %+v \n", cr.name, reply)
				cr.mu.Lock()
				if reply.Term > cr.currentTerm {
					cr.role = ROLE_BACKUP
					cr.currentTerm = reply.Term
					cr.votedFor = HAS_NOT_VOTED
					electionDone = true
					cr.mu.Unlock()
					break
				}
				// process reply, terminate if election succeeded
				if reply.VoteGranted {
					votes++
					DPrintf("Server %s got %d votes\n", cr.name, votes)
				}
				if votes > n/2 {
					cr.role = ROLE_MASTER
					// why are these values set like this?
					for name := range cr.nextIndex {
						cr.nextIndex[name] = len(cr.log)
						cr.matchIndex[name] = 0
					}
					electionDone = true
					cr.mu.Unlock()
					break
				}
				cr.mu.Unlock()
				break
			}
		case <-time.After(timer):
			{
				electionDone = true
				break
			}
		}
	}
	cr.mu.Lock()
	DPrintf("Server %s: Election terminated, role: %s\n", cr.name, cr.role)
	cr.mu.Unlock()
}

func (cr *Curp) RequestVote(args *RequestVoteArgs, reply *RequestVoteReply) error {
	// Your code here (2A, 2B).
	// Read the fields in "args",
	// and accordingly assign the values for fields in "reply".

	// TODO: verify this
	cr.mu.Lock()
	reply.Term = cr.currentTerm

	if args.Term < reply.Term {
		reply.Term = cr.currentTerm
		reply.VoteGranted = false
		cr.mu.Unlock()
		return nil
	}

	if args.Term > reply.Term {
		cr.currentTerm = args.Term
		reply.Term = cr.currentTerm
		// step down
		cr.role = ROLE_BACKUP
		cr.votedFor = HAS_NOT_VOTED

	}

	if cr.votedFor == HAS_NOT_VOTED || cr.votedFor == args.CandidateName {
		if len(cr.log) > 0 {
			if cr.log[len(cr.log)-1].Term > args.LastLogTerm ||
				(cr.log[len(cr.log)-1].Term == args.LastLogTerm && len(cr.log)-1 > args.LastLogIndex) {
				DPrintf("%s VOTE NO for %s's election because disagree logs\n", cr.name, args.CandidateName)
				reply.VoteGranted = false
				cr.mu.Unlock()
				return nil
			}
		}
		reply.Term = cr.currentTerm
		reply.VoteGranted = true
		cr.votedFor = args.CandidateName
		DPrintf("Server %s voted for %s \n", cr.name, cr.votedFor)
		cr.mu.Unlock()

		// reset election timeout
		cr.timeoutChan <- struct{}{}
	} else {
		cr.mu.Unlock()
	}
	return nil
}

func (cr *Curp) sendRequestVote(name string, args *RequestVoteArgs, reply *RequestVoteReply) error {
	ok := cr.peer_clients[name].Call("Curp.RequestVote", args, reply)
	return ok
}
