package main

import "fmt"

// func chSl(sl *[]int) *[]int {
// 	fmt.Printf("Sl1: %v, len: %d, cap: %d\n", sl, len(*sl), cap(*sl)) //2,2,2 L3 C3
// 	(*sl)[0] = 1
// 	*sl = append(*sl, 10)
// 	fmt.Printf("Sl2: %v, len: %d, cap: %d\n", sl, len(*sl), cap(*sl)) //1,2,2,10 L4 C6
// 	return sl
// }

func main() {
	// sl1 := &[]int{2, 2, 2}
	// fmt.Printf("Sl1-1: %v, len: %d, cap: %d\n", sl1, len(*sl1), cap(*sl1)) //2,2,2 L3 C3
	// sl2 := chSl(sl1)
	// fmt.Println("-----------")
	// fmt.Printf("Sl2: %v, len: %d, cap: %d\n", sl2, len(*sl2), cap(*sl2))   //1,2,2,10 L4 C6
	// fmt.Printf("Sl1-2: %v, len: %d, cap: %d\n", sl1, len(*sl1), cap(*sl1)) //1,2,2 L3 C3

	// slice := make([]int, 16, 16)
	// fmt.Printf("slice: len sl %d, sc %d", len(slice), cap(slice))
	// slice = append(slice, 1)
	// fmt.Printf("slice: len sl %d, sc %d", len(slice), cap(slice))

	arr := [11]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	slice := arr[2:4]
	fmt.Printf("slice: len sl %d, sc %d\n", len(slice), cap(slice))

	slice = arr[2:4:11]
	fmt.Printf("slice: len sl %d, sc %d\n", len(slice), cap(slice))
}
