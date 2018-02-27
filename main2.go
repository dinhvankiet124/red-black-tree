package main

import (
	"structure/tree"
	"fmt"
	"time"
	"math/rand"
)

func main() {
	t := tree.NewTDigest(50, 2*60)
	N := 1000000
	start := time.Now()
	for i := 0; i < N; i++ {
		t.Add(rand.Float32() * 1000, 1)
	}
	fmt.Println("Cost", time.Now().Sub(start))
	fmt.Println("Size", t.Size())
	fmt.Println(t.Percentile(0.95))
}
