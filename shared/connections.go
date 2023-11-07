package shared

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type NodeRole string

const (
	ROLE_MASTER    NodeRole = "Master"
	ROLE_BACKUP    NodeRole = "Backup"
	ROLE_CANDIDATE NodeRole = "Candidate"
	ROLE_WITNESS   NodeRole = "Witness"
)

type Node struct {
	Name string
	Role NodeRole
	Ip   string
	Port string
}

func isValidTopoRole(role NodeRole) bool {
	return role == ROLE_MASTER || role == ROLE_BACKUP || role == ROLE_WITNESS
}

func Parse(filename string, num_nodes int) map[string]*Node {
	file, err := os.Open(filename)
	if err != nil {
		log.Println("failed to open file ", filename)
	}

	reader := bufio.NewScanner(file)
	node_map := make(map[string]*Node)

	for i := 0; i < num_nodes; i++ {
		reader.Scan()
		text := reader.Text()
		raw_args := strings.Split(text, ",")
		cleaned_args := Map(raw_args, strings.TrimSpace)

		if len(cleaned_args) != 4 {
			fmt.Println("Invalid topo file format! Expected: \n<name>, <role>, <ip>, <port>")
			os.Exit(1)
		}

		node := new(Node)
		node.Name = cleaned_args[0]
		node.Role = NodeRole(cleaned_args[1])
		node.Ip = cleaned_args[2]
		node.Port = cleaned_args[3]

		if !isValidTopoRole(node.Role) {
			fmt.Println("Invalid topo file (bad role)! Expected roles: Master, Backup, Witness")
			os.Exit(1)
		}

		node_map[cleaned_args[0]] = node
	}
	return node_map
}

func ParseByRole(node_map map[string]*Node) (*Node, map[string]*Node, map[string]*Node) {
	witness_map := make(map[string]*Node)
	peer_map := make(map[string]*Node)

	var master_node *Node = nil
	for k, v := range node_map {
		switch v.Role {
		case ROLE_MASTER:
			peer_map[k] = v
			master_node = v
		case ROLE_WITNESS:
			witness_map[k] = v
		case ROLE_BACKUP:
			peer_map[k] = v
		}
	}
	return master_node, peer_map, witness_map
}
