package tree

import (
	"fmt"
)

type color bool
type direction int

const (
	black, red    color     = true, false
	leftd, rightd direction = 0, 1
)

type RedBlackTree struct {
	Root       *Node
	size       int
	dirtyNodes []float32
}

type Node struct {
	Mean   float32
	Weight float32
	Lowers float32
	Color  color
	Left   *Node
	Right  *Node
	Parent *Node
}

func NewNode(mean float32, weight float32, color color) *Node {
	return &Node{
		Mean:   mean,
		Weight: weight,
		Lowers: weight,
		Color:  color,
	}
}

func (node *Node) UpdateNode(that *Node) {
	node.Weight += that.Weight
	node.Mean += (that.Mean - node.Mean) * that.Weight / node.Weight
}

func (node *Node) Update(x float32, w float32) {
	node.Weight += w
	node.Mean += (x - node.Mean) * w / node.Weight
}

func (node *Node) Copy() *Node {
	return NewNode(node.Mean, node.Weight, node.Color)
}

func NewRedBlackTree() *RedBlackTree {
	return &RedBlackTree{
		Root: nil,
		size: 0,
	}
}

func (tree *RedBlackTree) First() *Node {
	return tree.Root
}

func (tree *RedBlackTree) Set(node *Node, dir direction, child *Node) {
	tree.addDirtyNodes(node.Mean)

	if dir == leftd {
		node.Left = child
	} else {
		node.Right = child
	}
}

func (tree *RedBlackTree) addDirtyNodes(x float32) {
	tree.dirtyNodes = append(tree.dirtyNodes, x)
}

func (tree *RedBlackTree) Insert(x float32, w float32) {
	tree.addDirtyNodes(x)

	var insertedNode *Node
	if tree.Root == nil {
		tree.Root = NewNode(x, w, red)
		insertedNode = tree.Root
	} else {
		node := tree.Root
		loop := true
		for loop {
			compare := x - node.Mean
			switch {
			case compare == 0:
				node.Update(x, w)
				return
			case compare < 0:
				if node.Left == nil {
					node.Left = NewNode(x, w, red)
					insertedNode = node.Left
					loop = false
				} else {
					node = node.Left
				}
			case compare > 0:
				if node.Right == nil {
					node.Right = NewNode(x, w, red)
					insertedNode = node.Right
					loop = false
				} else {
					node = node.Right
				}
			}
		}
		insertedNode.Parent = node
	}

	tree.insertCase1(insertedNode)
	tree.size++

	//tree.updateDirtyNodes()
}

func (tree *RedBlackTree) Remove(x float32) {
	node := tree.Root
loop:
	for node != nil {
		compare := x - node.Mean
		switch {
		case compare == 0:
			break loop
		case compare < 0:
			node = node.Left
		case compare > 0:
			node = node.Right
		}
	}

	if node == nil {
		return
	}
	if node.Left != nil && node.Right != nil {
		pred := node.Left.maximumNode()
		node.Mean = pred.Mean
		node.Weight = pred.Weight
		node.Lowers = pred.Lowers
		node = pred
	}

	var child *Node
	if node.Left == nil || node.Right == nil {
		if node.Right == nil {
			child = node.Left
		} else {
			child = node.Right
		}
		if node.Color == black {
			node.Color = nodeColor(child)
			tree.deleteCase1(node)
		}
		tree.replaceNode(node, child)
		if node.Parent == nil && child != nil {
			child.Color = black
		}
	}
	tree.size--

	//tree.updateDirtyNodes()
}

func (tree *RedBlackTree) updateDirtyNodes() {
	for _, mean := range tree.dirtyNodes {
		path := tree.lookupPath(mean)
		tree.updateLowers(path)
	}
	tree.dirtyNodes = tree.dirtyNodes[:0]
}

func (tree *RedBlackTree) updateLowers(path []*Node) {
	for i := len(path) - 1; i >= 0; i-- {
		node := path[i]
		if node.Left != nil && node.Right != nil {
			node.Lowers = node.Weight + node.Left.Lowers + node.Right.Lowers
		} else if node.Left != nil {
			node.Lowers = node.Weight + node.Left.Lowers
		} else if node.Right != nil {
			node.Lowers = node.Weight + node.Right.Lowers
		} else {
			node.Lowers = node.Weight
		}
	}
}

func (tree *RedBlackTree) GetLowerWeights(x float32) float32 {
	return tree.getLowerWeights(tree.Root, x)
}

func (tree *RedBlackTree) getLowerWeights(node *Node, x float32) float32 {
	if node == nil {
		return 0
	}

	if x > node.Mean {
		return node.Weight + tree.getLowerWeights(node.Right, x) + tree.getFullWeights(node.Left)
	} else {
		return tree.getLowerWeights(node.Left, x)
	}
}

func (tree *RedBlackTree) getFullWeights(node *Node) float32 {
	if node == nil {
		return 0
	} else {
		return node.Lowers
	}
}

func (tree *RedBlackTree) lookupPath(x float32) []*Node {
	path := make([]*Node, 0)
	node := tree.Root
loop:
	for node != nil {
		compare := x - node.Mean
		path = append(path, node)
		switch {
		case compare == 0:
			break loop
		case compare < 0:
			node = node.Left
		case compare > 0:
			node = node.Right
		}
	}
	return path
}

func (tree *RedBlackTree) Empty() bool {
	return tree.size == 0
}

func (tree *RedBlackTree) Size() int {
	return tree.size
}

func (tree *RedBlackTree) GetMin() *Node {
	var parent *Node
	current := tree.Root
	for current != nil {
		parent = current
		current = current.Left
	}
	return parent
}

func (tree *RedBlackTree) GetMax() *Node {
	var parent *Node
	current := tree.Root
	for current != nil {
		parent = current
		current = current.Right
	}
	return parent
}

func (tree *RedBlackTree) Neighbors(x float32) (*Node, *Node) {
	var floor, ceiling *Node = nil, nil
	node := tree.Root
	for node != nil {
		compare := x - node.Mean
		switch {
		case compare == 0:
			return node, node
		case compare < 0:
			ceiling = node
			node = node.Left
		case compare > 0:
			floor = node
			node = node.Right
		}
	}

	return floor, ceiling
}

func (tree *RedBlackTree) Floor(x float32) *Node {
	var floor *Node = nil
	node := tree.Root
	for node != nil {
		compare := x - node.Mean
		switch {
		case compare == 0:
			return node
		case compare < 0:
			node = node.Left
		case compare > 0:
			floor = node
			node = node.Right
		}
	}

	return floor
}

func (tree *RedBlackTree) Ceiling(x float32) *Node {
	var ceiling *Node = nil
	node := tree.Root
	for node != nil {
		compare := x - node.Mean
		switch {
		case compare == 0:
			return node
		case compare < 0:
			ceiling = node
			node = node.Left
		case compare > 0:
			node = node.Right
		}
	}

	return ceiling
}

func (tree *RedBlackTree) Clear() {
	tree.Root = nil
	tree.size = 0
}

func (tree *RedBlackTree) String() string {
	str := "RedBlackTree\n"
	if !tree.Empty() {
		output(tree.Root, "", true, &str)
	}
	return str
}

func (node *Node) String() string {
	return fmt.Sprintf("Mean=%f Weight=%f Lowers=%f", node.Mean, node.Weight, node.Lowers)
}

func output(node *Node, prefix string, isTail bool, str *string) {
	if node.Right != nil {
		newPrefix := prefix
		if isTail {
			newPrefix += "│   "
		} else {
			newPrefix += "    "
		}
		output(node.Right, newPrefix, false, str)
	}
	*str += prefix
	if isTail {
		*str += "└── "
	} else {
		*str += "┌── "
	}
	*str += node.String() + "\n"
	if node.Left != nil {
		newPrefix := prefix
		if isTail {
			newPrefix += "    "
		} else {
			newPrefix += "│   "
		}
		output(node.Left, newPrefix, true, str)
	}
}

func (node *Node) grandparent() *Node {
	if node != nil && node.Parent != nil {
		return node.Parent.Parent
	}
	return nil
}

func (node *Node) uncle() *Node {
	if node == nil || node.Parent == nil || node.Parent.Parent == nil {
		return nil
	}
	return node.Parent.sibling()
}

func (node *Node) sibling() *Node {
	if node == nil || node.Parent == nil {
		return nil
	}
	if node == node.Parent.Left {
		return node.Parent.Right
	}
	return node.Parent.Left
}

func (tree *RedBlackTree) rotateLeft(node *Node) {
	right := node.Right
	tree.replaceNode(node, right)
	tree.Set(node, rightd, right.Left)
	if right.Left != nil {
		right.Left.Parent = node
	}
	tree.Set(right, leftd, node)
	node.Parent = right
}

func (tree *RedBlackTree) rotateRight(node *Node) {
	left := node.Left
	tree.replaceNode(node, left)
	tree.Set(node, leftd, left.Right)
	if left.Right != nil {
		left.Right.Parent = node
	}
	tree.Set(left, rightd, node)
	node.Parent = left
}

func (tree *RedBlackTree) replaceNode(old *Node, new *Node) {
	if old.Parent == nil {
		tree.Root = new
	} else {
		if old == old.Parent.Left {
			tree.Set(old.Parent, leftd, new)
		} else {
			tree.Set(old.Parent, rightd, new)
		}
	}
	if new != nil {
		new.Parent = old.Parent
	}
}

func (tree *RedBlackTree) insertCase1(node *Node) {
	if node.Parent == nil {
		node.Color = black
	} else {
		tree.insertCase2(node)
	}
}

func (tree *RedBlackTree) insertCase2(node *Node) {
	if nodeColor(node.Parent) == black {
		return
	}
	tree.insertCase3(node)
}

func (tree *RedBlackTree) insertCase3(node *Node) {
	uncle := node.uncle()
	if nodeColor(uncle) == red {
		node.Parent.Color = black
		uncle.Color = black
		node.grandparent().Color = red
		tree.insertCase1(node.grandparent())
	} else {
		tree.insertCase4(node)
	}
}

func (tree *RedBlackTree) insertCase4(node *Node) {
	grandparent := node.grandparent()
	if node == node.Parent.Right && node.Parent == grandparent.Left {
		tree.rotateLeft(node.Parent)
		node = node.Left
	} else if node == node.Parent.Left && node.Parent == grandparent.Right {
		tree.rotateRight(node.Parent)
		node = node.Right
	}
	tree.insertCase5(node)
}

func (tree *RedBlackTree) insertCase5(node *Node) {
	node.Parent.Color = black
	grandparent := node.grandparent()
	grandparent.Color = red
	if node == node.Parent.Left && node.Parent == grandparent.Left {
		tree.rotateRight(grandparent)
	} else if node == node.Parent.Right && node.Parent == grandparent.Right {
		tree.rotateLeft(grandparent)
	}
}

func (node *Node) maximumNode() *Node {
	if node == nil {
		return nil
	}
	curr := node
	for curr.Right != nil {
		curr = curr.Right
	}
	return curr
}

func (tree *RedBlackTree) deleteCase1(node *Node) {
	if node.Parent == nil {
		return
	}
	tree.deleteCase2(node)
}

func (tree *RedBlackTree) deleteCase2(node *Node) {
	sibling := node.sibling()
	if nodeColor(sibling) == red {
		node.Parent.Color = red
		sibling.Color = black
		if node == node.Parent.Left {
			tree.rotateLeft(node.Parent)
		} else {
			tree.rotateRight(node.Parent)
		}
	}
	tree.deleteCase3(node)
}

func (tree *RedBlackTree) deleteCase3(node *Node) {
	sibling := node.sibling()
	if nodeColor(node.Parent) == black &&
		nodeColor(sibling) == black &&
		nodeColor(sibling.Left) == black &&
		nodeColor(sibling.Right) == black {
		sibling.Color = red
		tree.deleteCase1(node.Parent)
	} else {
		tree.deleteCase4(node)
	}
}

func (tree *RedBlackTree) deleteCase4(node *Node) {
	sibling := node.sibling()
	if nodeColor(node.Parent) == red &&
		nodeColor(sibling) == black &&
		nodeColor(sibling.Left) == black &&
		nodeColor(sibling.Right) == black {
		sibling.Color = red
		node.Parent.Color = black
	} else {
		tree.deleteCase5(node)
	}
}

func (tree *RedBlackTree) deleteCase5(node *Node) {
	sibling := node.sibling()
	if node == node.Parent.Left &&
		nodeColor(sibling) == black &&
		nodeColor(sibling.Left) == red &&
		nodeColor(sibling.Right) == black {
		sibling.Color = red
		sibling.Left.Color = black
		tree.rotateRight(sibling)
	} else if node == node.Parent.Right &&
		nodeColor(sibling) == black &&
		nodeColor(sibling.Right) == red &&
		nodeColor(sibling.Left) == black {
		sibling.Color = red
		sibling.Right.Color = black
		tree.rotateLeft(sibling)
	}
	tree.deleteCase6(node)
}

func (tree *RedBlackTree) deleteCase6(node *Node) {
	sibling := node.sibling()
	sibling.Color = nodeColor(node.Parent)
	node.Parent.Color = black
	if node == node.Parent.Left && nodeColor(sibling.Right) == red {
		sibling.Right.Color = black
		tree.rotateLeft(node.Parent)
	} else if nodeColor(sibling.Left) == red {
		sibling.Left.Color = black
		tree.rotateRight(node.Parent)
	}
}

func nodeColor(node *Node) color {
	if node == nil {
		return black
	}
	return node.Color
}

func (tree *RedBlackTree) Nodes() []*Node {
	nodes := make([]*Node, tree.size)
	it := tree.Iterator()
	for i := 0; it.Next(); i++ {
		nodes[i] = it.Node()
	}
	return nodes
}

type Iterator struct {
	tree     *RedBlackTree
	node     *Node
	position position
}

type position byte

const (
	begin, between, end position = 0, 1, 2
)

func (tree *RedBlackTree) Iterator() Iterator {
	return Iterator{tree: tree, node: nil, position: begin}
}

func (iterator *Iterator) Next() bool {
	if iterator.position == end {
		goto end
	}
	if iterator.position == begin {
		left := iterator.tree.GetMin()
		if left == nil {
			goto end
		}
		iterator.node = left
		goto between
	}
	if iterator.node.Right != nil {
		iterator.node = iterator.node.Right
		for iterator.node.Left != nil {
			iterator.node = iterator.node.Left
		}
		goto between
	}
	if iterator.node.Parent != nil {
		node := iterator.node
		for iterator.node.Parent != nil {
			iterator.node = iterator.node.Parent
			if node.Mean-iterator.node.Mean <= 0 {
				goto between
			}
		}
	}

end:
	iterator.node = nil
	iterator.position = end
	return false

between:
	iterator.position = between
	return true
}

func (iterator *Iterator) Node() *Node {
	return iterator.node
}

func (iterator *Iterator) Begin() {
	iterator.node = nil
	iterator.position = begin
}

func (iterator *Iterator) End() {
	iterator.node = nil
	iterator.position = end
}

func (iterator *Iterator) First() bool {
	iterator.Begin()
	return iterator.Next()
}
