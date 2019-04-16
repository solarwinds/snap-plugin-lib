package metrictree

type nodeStack []*Node

func (s *nodeStack) Push(v *Node) {
	*s = append(*s, v)
}

func (s *nodeStack) Pop() (*Node, bool) {
	if s.Len() == 0 {
		return nil, false
	}

	idx := s.Len() - 1
	node := (*s)[idx]
	*s = (*s)[:idx]

	return node, true
}

func (s *nodeStack) Len() int {
	return len(*s)
}

func (s *nodeStack) Empty() bool {
	return s.Len() == 0
}
