package main

import (
	"fmt"
	sorting "go-util/algorithms/sorting"
	datastructs "go-util/datastructures"
)

func main() {
	test1 := []int{3, 2, 4, 1, 5, 8, 12, 0}
	insertionSort := sorting.NewInsertionSort()
	result1 := insertionSort.Sort(test1)
	fmt.Println("InsertionSort:")
	fmt.Println(result1)

	test2 := []int{3, 2, 4, 1, 5, 8, 12, 0}
	mergeSort := sorting.NewMergeSort()
	result2 := mergeSort.Sort(test2, 0, len(test2)-1)
	fmt.Println("MergeSort:")
	fmt.Println(result2)

	test3 := []int{3, 2, 4, 1, 5, 8, 12, 0}
	bst := datastructs.NewBinarySearchTree()
	nodes := make(map[int]*datastructs.Node)
	for _, element := range test3 {
		node := datastructs.NewNode(element, nil)
		nodes[element] = node
		bst.Insert(node)
	}
	fmt.Println("BinarySearchTree:")
	fmt.Println("Tree after Insertion:")
	bst.Print(bst.Root)
	fmt.Println("Tree after Deletion:")
	bst.Delete(nodes[3])
	bst.Print(bst.Root)
}