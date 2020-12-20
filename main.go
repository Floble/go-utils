package main

import (
	"fmt"
	sorting "go-utils/algorithms/sorting"
	//string_matching "go-utils/algorithms/string_matching"
	//datastructs "go-utils/datastructures"
	//architecture "go-utils/algorithms/architecture"
)

func main() {
	/* test1 := []int{3, 2, 4, 1, 5, 8, 12, 0}
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
	fmt.Println(shifts) */

	/* yuma := architecture.NewYuma("/home/floble/go/src/go-utils/algorithms/architecture/example/hosts", "/home/floble/go/src/go-utils/algorithms/architecture/example/build.yml", "/home/floble/go/src/go-utils/algorithms/architecture/example/roles/")

	fmt.Println("Role name - configuration:")
	for role, config := range yuma.Roles {
		fmt.Printf("%s", role)
		fmt.Printf(" - ")
		fmt.Printf("%b", config)
		fmt.Println()
	}

	fmt.Println()
	fmt.Println()

	//result := yuma.DetermineDeploymentPlan_Backtracking(0, 0)
	//result := yuma.DetermineDeploymentPlan_Dfs(0, 0)
	//result := yuma.DetermineDeploymentPlan_Greedy_Recursive(0, 0)
	result := yuma.DetermineDeploymentPlan_Greedy_Iterative()

	if result {
		fmt.Println("All Ansible roles are included in the deploymentplan")
	} else {
		fmt.Println("The deploymentplan does not include all Ansible roles")
	}
	fmt.Println()

	fmt.Println("Deploymentplan (configurations):")
	for i := 1; i < len(yuma.DeploymentPlan); i++ {
		fmt.Printf("%d", i)
		fmt.Printf(" - ")
		fmt.Printf("%b", yuma.DeploymentPlan[i])
		fmt.Println()
	}

	fmt.Println()
	fmt.Println()
	fmt.Println("Deploymentplan (role names):")

	yuma.PrintDeploymentPlan() */

	/* test6 := []int{4, 1, 3, 2, 16, 9, 10, 14, 8, 7}

	heap := datastructs.NewHeap(test6)
	heap.Sort()

	for _, n := range heap.Array {
		fmt.Println(n)
	} */

	/* test7 := []int{4, 1, 3, 2, 16, 9, 10, 14, 8, 7}

	pq := datastructs.NewMaxPriorityQueue(test7)
	pq.Heap.BuildMaxHeap()
	
	pq.Insert(21)

	fmt.Println(pq.Maximum()) */

	/* test8 := []int{3, 2, 4, 1, 5, 8, 12, 0}
	quickSort := sorting.NewQuickSort()
	quickSort.Sort(test8)
	fmt.Println("QuickSort:")
	fmt.Println(test8) */

	/* test9 := []int{3, 8, 2, 9, 1, 13, 6, 5, 0, 2, 3, 9}
	countingSort := sorting.NewCountingSort()
	result := countingSort.Sort(test9, 13)
	fmt.Println(result) */

	test10 := []int{43, 12, 21}
	radixSort := sorting.NewRadixSort()
	result := radixSort.Sort(test10, 2)
	fmt.Println(result)
}