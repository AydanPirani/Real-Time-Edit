package curp

import (
	"net/rpc"
	. "rtclbedit/shared"

	mapset "github.com/deckarep/golang-set/v2"
)

type Witness struct {
	name          string
	unsynced      mapset.Set[string]
	master_client *rpc.Client
}

/**
 * RPC types
 */
type DropArgs struct {
}
type DropReply struct {
}
type RecordArgs struct {
	Command interface{}
}
type RecordReply struct {
}

/**
 * RPC functions
 */
func (w *Witness) Drop(args DropArgs, reply *DropReply) error { // dropRPC called by master
	DPrintf("%s receives Drop RPC %v+", w.name, args)
	return nil
}

func (w *Witness) Record(args RecordArgs, reply *RecordReply) error { // recordRPC called by client
	DPrintf("%s receives Record RPC %v+", w.name, args)
	return nil
}

func (w *Witness) WitnessLifetime() {
	for {

	}
}
