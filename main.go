package main

import (
	"fmt"
	"go-basics/sum"
)

func main() {
	nums := []int{1, 2, 3, 4, 5}
	fmt.Println(sum.Sum(nums...))
}
