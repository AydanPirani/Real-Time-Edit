package shared

import "log"

func Map(list []string, f func(string) string) []string {
	result := make([]string, len(list))
	for i, v := range list {
		result[i] = f(v)
	}
	return result
}

// Debugging
const Debug = 1

func DPrintf(format string, a ...interface{}) (n int, err error) {
	if Debug > 0 {
		log.Printf(format, a...)
	}
	return
}

const (
	SUCCESS_STATUS = iota
	FAILURE_STATUS
)
