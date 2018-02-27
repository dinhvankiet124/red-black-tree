package tree

import (
	"time"
	"math/rand"
)

type TDigest struct {
	Compression float32
	Count       float32
	DeltaT      float32
	Min         float32
	Max         float32
	LastClean   time.Time
	Tree        *RedBlackTree
}

func NewTDigest(compression float32, deltaT float32) *TDigest {
	return &TDigest{
		Compression: compression,
		Count:       0,
		DeltaT:      deltaT,
		Min:         -1,
		Max:         0,
		LastClean:   time.Now(),
		Tree:        NewRedBlackTree(),
	}
}

func (t *TDigest) Size() int {
	return t.Tree.Size()
}

func (t *TDigest) Update(node *Node, x float32, w float32) {
	newNode := node.Copy()
	t.Tree.Remove(node.Mean)
	newNode.Update(x, w)
	t.Tree.Insert(newNode.Mean, newNode.Weight)
}

func (t *TDigest) FindClosest(x float32) []*Node {
	//floor, ceiling := t.Tree.Neighbors(x)
	//
	//if ceiling == nil {
	//	return []*Node{floor}
	//}
	//
	//if floor == nil {
	//	return []*Node{ceiling}
	//}

	//ceiling := t.Tree.Ceiling(x)
	//if ceiling == nil {
	//	return []*Node{t.Tree.Floor(x)}
	//}

	return []*Node{t.Tree.Floor(x)}

	//floor := t.Tree.Floor(x)
	//if floor == nil {
	//	return []*Node{ceiling}
	//}
	//
	//if Abs(floor.Mean-x) < Abs(ceiling.Mean-x) {
	//	return []*Node{floor}
	//} else if Abs(floor.Mean-x) == Abs(ceiling.Mean-x) && (ceiling.Mean != floor.Mean) {
	//	return []*Node{floor, ceiling}
	//} else {
	//	return []*Node{ceiling}
	//}
}

func (t *TDigest) Add(x float32, w float32) {
	t.Count += w

	if t.Min == -1 || t.Min > x {
		t.Min = x
	} else if t.Max < x {
		t.Max = x
	}

	if t.Tree.Size() == 0 {
		t.Tree.Insert(x, w)
		return
	}

	S := t.FindClosest(x)

	for _, node := range S {
		if w <= 0 || node == nil {
			break
		}

		sum := t.Tree.GetLowerWeights(node.Mean)
		q := ((node.Weight / 2.0) + sum) / t.Count
		k := 24.0 * t.Count * q * (1 - q) / t.Compression

		if node.Weight+w > k {
			continue
		}

		deltaW := Min(k-node.Weight, w)
		t.Update(node, x, deltaW)
		w -= deltaW
	}

	if w > 0 {
		t.Tree.Insert(x, w)
	}

	if float32(t.Size()) > 20*t.Compression {
		t.Compress()
	}
}

func (t *TDigest) Compress() {
	nodes := t.Tree.Nodes()
	for i := t.Size() - 1; i > 0; i-- {
		other := rand.Intn(i + 1)
		tmp := nodes[other]
		nodes[other] = nodes[i]
		nodes[i] = tmp
	}

	t.Tree.Clear()
	for _, node := range nodes {
		t.Add(node.Mean, node.Weight)
	}
}

func (t *TDigest) Percentile(q float32) float32 {
	if q < 0 || q > 1 {
		panic("q should be in [0,1]")
		return -1
	}

	tree := t.Tree
	size := tree.Size()
	if size == 0 {
		return -1
	} else if size == 1 {
		return t.Tree.First().Mean
	}

	index := q * t.Count

	var curr *Node = nil
	var wSoFar float32 = 0

	it := tree.Iterator()
	for i := 0; it.Next(); i++ {
		if i == 0 {
			curr = it.Node()
			wSoFar = curr.Weight / 2.0
			if index < wSoFar {
				return (t.Min*index + curr.Mean*(wSoFar-index)) / wSoFar
			}
			continue
		}

		node := it.Node()

		dw := (curr.Weight + node.Weight) / 2.0
		if wSoFar+dw > index {
			z1 := index - wSoFar
			z2 := wSoFar + dw - index
			return (node.Mean*z1 + curr.Mean*z2) / (z1 + z2)
		}
		wSoFar += dw
		curr = node
	}

	z1 := index - wSoFar
	z2 := curr.Weight/2.0 - z1
	return (curr.Mean*z2 + t.Max*z1) / (z1 + z2)
}
