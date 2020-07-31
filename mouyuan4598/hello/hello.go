package main

import (
	"fmt"

	bubblesort "github.com/mouyuan4598/hello/sort"
)

func main() {
	arr := []int{9, 7, 1, 5, 3}
	array := bubblesort.Sort(arr)
	var str string
	fmt.Scanln(&str)
	fmt.Println("array: " + str)
	fmt.Println(array)
}
