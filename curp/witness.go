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
}
type DropReply struct {
	Status int
}
type RecordArgs struct {
	Name string
}
type RecordReply struct {
	Status int
}

/**
 * RPC functions
 */
func (w *Witness) Drop(args DropArgs, reply *DropReply) error { // dropRPC called by master
	fmt.Println("RECEIVED DROP RPC MESSAGE")
	reply.Status = SUCCESS_STATUS
	return nil
}

func (w *Witness) Record(args RecordArgs, reply *RecordReply) error { // recordRPC called by client
	fmt.Println("RECEIVED RECORD RPC MESSAGE")
	reply.Status = SUCCESS_STATUS
	return nil
}
