package curp

import (
	"fmt"
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
	DropLog []LogEntry
}
type DropReply struct {
	Success bool
}
type RecordArgs struct {
	Command interface{}
}
type RecordReply struct {
	Success bool
}

/**
 * RPC functions
 */
func (w *Witness) Drop(args DropArgs, reply *DropReply) error { // dropRPC called by master
	fmt.Println("RECEIVED DROP RPC MESSAGE")
	for _, log := range args.DropLog {
		w.unsynced.Remove(log.Command.(string))
	}
	reply.Success = true
	return nil
}

func (w *Witness) Record(args RecordArgs, reply *RecordReply) error { // recordRPC called by client
	fmt.Println("RECEIVED RECORD RPC MESSAGE")
	w.unsynced.Append(args.Command.(string))
	reply.Success = true
	return nil
}

func (w *Witness) WitnessLifetime() {
	for {

	}
}
