package main

import "fmt"

func printArray(arr [5]int)  {
	for i,v := range arr {
		fmt.Println(i,v)
	}
	arr[0] = 100

}

func main() {
	var arr1 [5]int
	arr2 := [3]int{1,2,3}
	arr3 := [...]int{2,4,6,8,10}
	var grid [4][5]int
	fmt.Println(arr1,arr2,arr3)
	fmt.Println(grid)

	printArray(arr1)
	printArray(arr3)
	fmt.Println(arr1,arr2,arr3)
}
