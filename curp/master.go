package curp

import (
	. "rtclbedit/shared"
)

type ExecuteArgs struct {
	Command interface{}
}
type ExecuteReply struct {
	Success bool
}
type SyncArgs struct {
}
type SyncReply struct {
}

/**
 * RPC functions
 */
func (cr *Curp) Execute(args ExecuteArgs, reply *ExecuteReply) error { // executeRPC called by clients to master
	executeMessage := ExecuteMsg{
		CommandValid: true,
		Command:      cr.log[cr.syncedIndex].Command,
		CommandIndex: cr.syncedIndex + 1,
	}
	cr.appChan <- executeMessage
	DPrintf("Leader %d executing command %d\n", cr.name, executeMessage.CommandIndex)
	// ordering in the background
	go cr.Start(args.Command)
	return nil
}

func (cr *Curp) Sync(args SyncArgs, reply *SyncReply) error { // syncRPC called by clients to master
	return nil
}
