package algorithm

type StackNode struct {
	Val  interface{}
	Next *StackNode
}

type Stack struct {
	length int
	root   *StackNode
}

func NewStack() *Stack {
	return &Stack{
		length: 0,
		root:   nil,
	}
}

func (s *Stack) Push(val interface{}) {
	node := &StackNode{Val: val, Next: s.root}
	s.root = node
	s.length++
}

func (s *Stack) Pop() interface{} {
	if s.IsEmpty() {
		return nil
	}
	s.length--
	node := s.root
	s.root = node.Next
	node.Next = nil
	return node.Val
}

func (s *Stack) IsEmpty() bool {
	return s.length == 0
}

func (s *Stack) Length() int {
	return s.length
}
