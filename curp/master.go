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
	Document string
}

/**
 * RPC functions
 */
func (cr *Curp) Execute(args ExecuteArgs, reply *ExecuteReply) error { // executeRPC called by clients to master
	// DPrintf("not sutpid")
	executeMessage := ExecuteMsg{
		CommandValid: true,
		Command:      args.Command,
		CommandIndex: cr.syncedIndex + 1,
	}
	cr.appChan <- executeMessage
	cr.applyCmd(args.Command)
	DPrintf("Leader %s executing command %d\n", cr.name, cr.syncedIndex+1)
	// ordering in the background
	go cr.Start(args.Command)
	return nil
}

func (cr *Curp) Sync(args SyncArgs, reply *SyncReply) error { // syncRPC called by clients to master

	DPrintf("Leader %s executing command %d\n", cr.name, cr.syncedIndex+1)
	// no ordering? read does not change state
	return nil
}
