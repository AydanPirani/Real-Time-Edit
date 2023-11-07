package curp

import (
	"net/rpc"

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
}
type RecordReply struct {
}

/**
 * RPC functions
 */
func (w *Witness) Drop(args DropArgs, reply *DropReply) error { // dropRPC called by master
	return nil
}

func (w *Witness) Record(args RecordArgs, reply *RecordReply) error { // recordRPC called by client
	return nil
}
