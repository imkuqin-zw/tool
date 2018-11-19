package algorithm

type order uint8

const (
	PRE_ORDER order = iota
	IN_ORDER
	POST_ORDER
	LEVEL_ORDER
)

type BTreeArrNode struct {
	Val interface{}
	L   int
	R   int
}

type BTreeNode struct {
	Val interface{}
	L   *BTreeNode
	R   *BTreeNode
}

type BTree struct {
	length int
	root   *BTreeNode
}

func NewTree() *BTree {
	return &BTree{
		length: 0,
		root:   nil,
	}
}

func (bt *BTree) BuildBTree(nodeArr []BTreeArrNode) {
	length := len(nodeArr)
	hasFather := make([]bool, length)
	if length > 0 {
		tmp := make([]*BTreeNode, 0, length)
		for _, item := range nodeArr {
			tmp = append(tmp, &BTreeNode{Val: item.Val})
		}
		for i, item := range nodeArr {
			if item.L != -1 {
				tmp[i].L = tmp[item.L]
				hasFather[item.L] = true
			}
			if item.R != -1 {
				tmp[i].R = tmp[item.R]
				hasFather[item.R] = true
			}
		}
		i := 0
		for i < length && hasFather[i] {
			i++
		}
		bt.root = tmp[i]
	}
	return
}

func (bt *BTree) IsEmpty() bool {
	return bt.length == 0
}

func (bt *BTree) Length() int {
	return bt.length
}

func preOrderTraversal(root *BTreeNode, out []interface{}) {
	t := root
	stack := NewStack()
	for t != nil || !stack.IsEmpty() {
		for t != nil {
			out = append(out, t.Val)
			stack.Push(t)
			t = t.L
		}
		if !stack.IsEmpty() {
			elem := stack.Pop().(*BTreeNode)
			t = elem.R
		}
	}
}

func inOrderTraversal(root *BTreeNode, out []interface{}) {
	t := root
	stack := NewStack()
	for t != nil || !stack.IsEmpty() {
		for t != nil {
			stack.Push(t)
			t = t.L
		}
		if !stack.IsEmpty() {
			elem := stack.Pop().(*BTreeNode)
			out = append(out, elem.Val)
			t = elem.R
		}
	}
}

func postOrderTraversal(root *BTreeNode, out []interface{}) {
	t := root
	stack := NewStack()
	for t != nil || !stack.IsEmpty() {
		for t != nil {
			stack.Push(t)
			t = t.L
		}
		if !stack.IsEmpty() {
			elem := stack.Pop().(*BTreeNode)
			if elem.R != nil {
				t = elem.R
				elem.R = nil
				stack.Push(elem)
			} else {
				out = append(out, elem.Val)
			}
		}
	}
}

func levelOrderTraversal(root *BTreeNode, out []interface{}) {
	if root != nil {
		queue := NewQueue()
		queue.Push(root)
		for !queue.IsEmpty() {
			elem := queue.Pop().(*BTreeNode)
			out = append(out, elem.Val)
			if elem.L != nil {
				queue.Push(elem.L)
			}
			if elem.R != nil {
				queue.Push(elem.R)
			}
		}
	}
}

func (bt *BTree) Traversal(o order) []interface{} {
	out := make([]interface{}, 0, bt.length)
	if bt.IsEmpty() {
		return out
	}
	switch o {
	case PRE_ORDER:
		preOrderTraversal(bt.root, out)
	case IN_ORDER:
		inOrderTraversal(bt.root, out)
	case POST_ORDER:
		postOrderTraversal(bt.root, out)
	case LEVEL_ORDER:
		levelOrderTraversal(bt.root, out)
	}
	return out
}

func traversalLeaveNode(root *BTreeNode, out []interface{}) {
	if root != nil {
		if root.L == nil && root.R == nil {
			out = append(out, root)
		}
		traversalLeaveNode(root.L, out)
		traversalLeaveNode(root.R, out)
	}
}

func (bt *BTree) TraversalLeaveNode() []interface{} {
	if bt.IsEmpty() {
		return nil
	}
	result := make([]interface{}, 0)
	traversalLeaveNode(bt.root, result)
	return result
}

func getBTreeHeight(root *BTreeNode) int {
	if root != nil {
		l := getBTreeHeight(root.L)
		r := getBTreeHeight(root.R)
		if l >= r {
			return l + 1
		} else {
			return r + 1
		}
	}
	return 0
}

func (bt *BTree) GetBTreeHeight() int {
	return getBTreeHeight(bt.root)
}

func find(root *BTreeNode, val interface{}, compare searchCompare) *BTreeNode {
	if root == nil {
		return nil
	}
	var node *BTreeNode
	tmp := compare(val, root.Val)
	if tmp == 0 {
		node = root
	} else if tmp > 0 {
		node = find(root.R, val, compare)
	} else {
		node = find(root.L, val, compare)
	}
	return node
}

func (bt *BTree) Find(val interface{}, compare searchCompare) *BTreeNode {
	return find(bt.root, val, compare)
}

func findMin(root *BTreeNode) *BTreeNode {
	p := root
	for p.L != nil {
		p = p.L
	}
	return p
}

func (bt *BTree) FindMin() *BTreeNode {
	if bt.root == nil {
		return nil
	}
	return findMin(bt.root)
}

func findMax(root *BTreeNode) *BTreeNode {
	p := root
	for p.R != nil {
		p = p.R
	}
	return p
}

func (bt *BTree) FindMax() *BTreeNode {
	if bt.root == nil {
		return nil
	}
	return findMax(bt.root)
}

func bsTreeInsert(root *BTreeNode, val interface{}, compare searchCompare) *BTreeNode {
	if root == nil {
		return &BTreeNode{Val: val}
	}
	tmp := compare(val, root.Val)
	if tmp < 0 {
		root.L = bsTreeInsert(root.L, val, compare)
	} else if tmp > 0 {
		root.R = bsTreeInsert(root.R, val, compare)
	}
	return root
}

func (bt *BTree) BSTreeInsert(val interface{}, compare searchCompare) {
	bt.root = bsTreeInsert(bt.root, val, compare)
}

func bsTreeDelete(root *BTreeNode, val interface{}, compare searchCompare) *BTreeNode {
	if root != nil {
		tmp := compare(val, root.Val)
		if tmp > 0 {
			root.R = bsTreeDelete(root.R, val, compare)
		} else if tmp < 0 {
			root.L = bsTreeDelete(root.L, val, compare)
		} else {
			if root.L != nil && root.R != nil {
				root.Val = findMin(root.R).Val
				root.R = bsTreeDelete(root.R, val, compare)
			} else {
				tmp := root
				if root.L != nil {
					root = root.L
					tmp.L = nil
				} else if root.R != nil {
					root = root.R
					tmp.R = nil
				}
			}
		}
	}
	return root
}

func (bt *BTree) BSTreeDelete(val interface{}, compare searchCompare) {
	bt.root = bsTreeDelete(bt.root, val, compare)
}
