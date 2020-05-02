package main

import (
	"fmt"
	sorting "go-util/algorithms/sorting"
)

func main() {
	numbers := []int{3, 2, 4, 1, 5, 8, 12, 0}

	insertionSort := sorting.NewInsertionSort()
	sortedNumbers := insertionSort.Sort(numbers)
	fmt.Println(sortedNumbers)
}