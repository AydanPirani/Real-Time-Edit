package curp

import "net/rpc"

type Backup struct {
	name          string
	master_client map[string]*rpc.Client
}
