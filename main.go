package main

import (
	"fmt"

	"structure/tree"
)

func main() {
	rbt := tree.NewRedBlackTree()
	arr := []float32{5, 6, 12, 1, 7, 9, 13}
	for _, a := range arr {
		rbt.Insert(a, 1)
	}
	for _, a := range arr {
		fmt.Println(rbt.GetLowerWeights(a))
	}
	fmt.Println("after")
	rbt.Remove(6)
	for _, a := range arr {
		fmt.Println(rbt.GetLowerWeights(a))
	}
}

