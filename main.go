package main

import (
	"fmt"
	sorting "go-util/algorithms/sorting"
	string_matching "go-util/algorithms/string_matching"
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

	test4 := []int{3, 2, 4, 1, 5, 8, 12, 0}
	stack := datastructs.NewStack()
	for _, e := range test4 {
		stack.Push(e)
	}
	fmt.Println("Stack:")
	for i := 0; i < 5; i++ {
		err, e := stack.Pop()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(e)
	}

	test5 := []int{3, 2, 4, 1, 5, 8, 12, 0}
	queue := datastructs.NewQueue()
	for _, e := range test5 {
		queue.Enqueue(e)
	}
	fmt.Println("Queue:")
	for i := 0; i < 5; i++ {
		err, e := queue.Dequeue()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(e)
	}

	test6 := []int{3, 2, 4, 1, 5, 8, 12, 0}
	list := datastructs.NewLinkedList()
	for _, e := range test6 {
		list.Insert(e)
	}
	list.Delete(8)
	fmt.Println("LinkedList:")
	list.Print()

	s := "adsgwadsxdsgwadsgz"
	p := "dsgwadsgz"
	fmt.Println("HFSDFSDFSDAFSFSFSFSF")
	kmp := string_matching.NewKMP()
	shifts := kmp.Match(s, p)
	fmt.Println("Knuth-Morris-Pratt:")
	fmt.Println(shifts)
}