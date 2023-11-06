package curp

import "net/rpc"

type Backup struct {
	name          string
	master_client map[string]*rpc.Client
}

type ReplicateArgs struct {
}
type ReplicateReply struct {
}

func (b *Backup) Replicate(args ReplicateArgs, reply *ReplicateReply) error { // called by master
	return nil
}
