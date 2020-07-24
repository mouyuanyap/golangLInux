package main

import (
	"fmt"

	bubblesort "github.com/mouyuan4598/hello/sort"
)

func main() {
	arr := []int{4, 2, 1, 5, 3}
	array := bubblesort.Sort(arr)
	fmt.Println("array")
	fmt.Println(array)
}
