package main

import (
	"fmt"
	"os"
	"rtclbedit/curp"
	. "rtclbedit/shared"
)

func main() {
	args := os.Args
	if len(os.Args) != 3 {
		fmt.Println("[usage]: " + args[0] + " <identifier> <configuration file>")
		os.Exit(1)
	}
	name := args[1]
	filename := args[2]

	node_map := Parse(filename)
	master_map, witness_map, backup_map := ParseByRole(node_map)

	switch node_map[name].Role {
	case RoleMaster:
		curp.InitMaster(name, witness_map, backup_map)
	case RoleWitness:
		curp.InitWitness(name, master_map)
	case RoleBackup:
		curp.InitBackup(name, master_map)
	}

	curp.InitRPC(name, node_map)
	for {

	}
}
