package shared

import "log"

func Map_begin(m map[string]*Node) *Node {
	for _, v := range m {
		return v
	}
	return nil
}

// Debugging
const Debug = 0

func DPrintf(format string, a ...interface{}) (n int, err error) {
	if Debug > 0 {
		log.Printf(format, a...)
	}
	return
}
