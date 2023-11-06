package shared

func Map_begin(m map[string]*Node) *Node {
	for _, v := range m {
		return v
	}
	return nil
}
