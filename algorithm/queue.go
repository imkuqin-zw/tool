package algorithm

type QueueNode struct {
	Val  interface{}
	Next *QueueNode
	Prev *QueueNode
}

type Queue struct {
	length int
	root   *QueueNode
	last   *QueueNode
}

func NewQueue() *Queue {
	return &Queue{
		length: 0,
		root:   nil,
		last:   nil,
	}
}

func (q *Queue) Push(val interface{}) {
	node := &QueueNode{
		Val:  val,
		Next: nil,
	}
	if q.root == nil {
		q.root, q.last = node, node
	} else {
		node.Prev = q.last
		q.last.Next = node
		q.last = node
	}
	q.length++
}

func (q *Queue) Pop() interface{} {
	if q.IsEmpty() {
		return nil
	}
	q.length--
	node := q.root
	q.root = node.Next
	if node.Next != nil {
		node.Next.Prev = nil
		node.Next = nil
	} else {
		q.last = nil
	}
	return node.Val
}

func (q *Queue) IsEmpty() bool {
	return q.length == 0
}

func (q *Queue) Length() int {
	return q.length
}
