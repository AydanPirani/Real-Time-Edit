package shared

import (
	"bufio"
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

func Parse(filename string) map[string]*Node {

	file, err := os.Open(filename)
	if err != nil {
		log.Println("failed to open file ", filename)
	}
	reader := bufio.NewScanner(file)
	num_nodes := 5 //hardcoded number
	node_map := make(map[string]*Node)
	for i := 0; i < num_nodes; i++ {
		reader.Scan()
		text := reader.Text()
		test_split := strings.Split(text, " ")
		node := new(Node)
		node.Name = test_split[0]
		node.Role = NodeRole(test_split[1])
		node.Ip = test_split[2]
		node.Port = test_split[3]

		node_map[test_split[0]] = node
	}
	return node_map
}

func ParseByRole(node_map map[string]*Node) (map[string]*Node, map[string]*Node) {
	witness_map := make(map[string]*Node)
	peer_map := make(map[string]*Node)
	for k, v := range node_map {
		switch v.Role {
		case ROLE_MASTER:
			peer_map[k] = v
		case ROLE_WITNESS:
			witness_map[k] = v
		case ROLE_BACKUP:
			peer_map[k] = v
		}
	}
	return peer_map, witness_map
}
