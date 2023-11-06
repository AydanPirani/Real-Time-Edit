package curp

import (
	"net/rpc"
	. "rtclbedit/shared"
)

type Master struct {
	name            string
	unsynced        []Operation
	witness_clients map[string]*rpc.Client
	backup_clients  map[string]*rpc.Client
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
func (m *Master) Execute(args ExecuteArgs, reply *ExecuteReply) error { // executeRPC called by clients
	return nil
}

func (m *Master) Sync(args SyncArgs, reply *SyncReply) error { // syncRPC called by clients
	return nil
}

/**
 * Non-RPC functions
 */
func (m *Master) AsyncOrdering() {

}
