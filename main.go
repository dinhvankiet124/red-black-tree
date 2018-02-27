package main

import (
	"fmt"

	"structure/tree"
)

func main() {
	rbt := tree.NewRedBlackTree()
	for i := 0; i < 1000; i++ {
		rbt.Insert(float32(i), 1)
	}
	//fmt.Println(rbt.String())
	//fmt.Println(rbt.Floor(10))
	//fmt.Println(rbt.Ceiling(10))
	var sum float32 = 0
	for i := 0; i < 1000; i++ {
		if i % 3 == 0 {
			rbt.Remove(float32(i))
			rbt.Insert(float32(1000 + i), 1)
			continue
		}
		a := rbt.GetLowerWeights(float32(i))
		fmt.Println(a)
		sum += a
	}
	fmt.Println("Final", sum)
	//fmt.Println("after")
	//rbt.Remove(6)
	//rbt.Remove(9)
	//for _, a := range arr {
	//	fmt.Println(rbt.GetLowerWeights(a))
	//}
}

