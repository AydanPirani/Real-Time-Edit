package curp

import (
	"net/rpc"
	. "rtclbedit/shared"
)

type Backup struct {
	name          string
	master_client *rpc.Client
}

type OrderAsyncArgs struct {
	Term         int
	LeaderName   string
	PrevLogIndex int
	PrevLogTerm  int
	Entries      []LogEntry
	LeaderSynced int
}
type OrderAsyncReply struct {
	Term    int
	Success bool
}

func (cr *Curp) OrderAsync(args *OrderAsyncArgs, reply *OrderAsyncReply) {
	// TODO: reply.term ???
	// 1. Return if term < currentTerm
	cr.mu.Lock()
	reply.Term = cr.currentTerm

	if args.Term < cr.currentTerm {
		reply.Success = false
		cr.mu.Unlock()
		return
	}
	// 2. If term > currentTerm, currentTerm ← term
	if args.Term > cr.currentTerm {
		cr.currentTerm = args.Term
		reply.Term = cr.currentTerm
		cr.votedFor = HAS_NOT_VOTED
	}
	// 3. If candidate or leader, step down
	if args.LeaderName != cr.name && cr.role != ROLE_BACKUP {
		cr.role = ROLE_BACKUP
	}
	cr.mu.Unlock()
	// 4. Reset election timeout
	cr.timeoutChan <- struct{}{}

	cr.mu.Lock()
	defer cr.mu.Unlock()
	// 5. Return failure if log doesn’t contain an entry at prevLogIndex whose term matches prevLogTerm
	if args.PrevLogIndex != -1 && args.PrevLogIndex < len(cr.log) {
		DPrintf("args.PrevLogIndex %d args.PrevLogTerm %d\n", args.PrevLogIndex, args.PrevLogTerm)
		DPrintf("c.log[args.PrevLogIndex] %+v\n", cr.log[args.PrevLogIndex])
	}
	if args.PrevLogIndex >= len(cr.log) || (args.PrevLogIndex != -1 && cr.log[args.PrevLogIndex].Term != args.PrevLogTerm) {
		reply.Success = false
		return
	}
	reply.Success = true
	DPrintf("args.Entries: %+v", args.Entries)
	// 6. If existing entries conflict with new entries, delete all existing entries starting with first conflicting entry
	i := 0
	for i < len(args.Entries) {
		if args.PrevLogIndex+i+1 >= len(cr.log) {
			break
		}
		if cr.log[args.PrevLogIndex+i+1].Term != args.Entries[i].Term {
			cr.log = cr.log[:args.PrevLogIndex+i+1]
			break
		}
		i++
	}
	DPrintf("cr.name: %d cr.log: %+v\n", cr.name, cr.log)
	// 7. Append any new entries not already in the log
	cr.log = append(cr.log, args.Entries[i:]...)

	// 8. Advance state machine with newly committed entries
	for args.LeaderSynced > cr.syncedIndex && cr.syncedIndex < args.PrevLogIndex+len(args.Entries)+1 {
		executeMessage := ExecuteMsg{
			CommandValid: true, Command: cr.log[cr.syncedIndex].Command, CommandIndex: cr.syncedIndex + 1,
		}
		cr.appChan <- executeMessage
		DPrintf("Follower %d executing command %d\n", cr.name, executeMessage.CommandIndex)
		cr.syncedIndex++
	}

}
