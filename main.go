package main

import (
	//"fmt"
	//"os"
	//sorting "go-utils/algorithms/sorting"
	//string_matching "go-utils/algorithms/string_matching"
	//datastructs "go-utils/datastructures"
	"fmt"
	search "go-utils/algorithms/artificialintelligence/search"
	iac "go-utils/infrastructureascode"
	//orderstatistics "go-utils/algorithms/orderstatistics"
	//machine_learning "go-utils/algorithms/machinelearning"
	//helper "go-utils/helper"
	//ec2 "go-utils/cloud/aws/ec2"
	//"math"
	//"gonum.org/v1/gonum/mat"
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

	/* test10 := []int{43, 12, 21}
	radixSort := sorting.NewRadixSort()
	result := radixSort.Sort(test10, 2)
	fmt.Println(result) */

	/* test11 := []float64{0.78, 0.17, 0.39, 0.26, 0.72, 0.94, 0.21, 0.12, 0.23, 0.68}
	bucketSort := sorting.NewBucketSort()
	result := bucketSort.Sort(test11)
	fmt.Println(result) */

	/* test12 := []int{0, 23, 12, 14, 1723, 35, 99}
	selection := orderstatistics.NewSelection()
	i := selection.Select(test12, 0, len(test12) - 1, 6)
	fmt.Println(i) */

	/* p1 := machine_learning.NewPoint(vector.NewWithValues([]float64 {2.0}), 4)
	p2 := machine_learning.NewPoint(vector.NewWithValues([]float64 {4.0}), 2)
	points := []*machine_learning.Point{p1, p2} */

	/* wTrue := machine_learning.GenerateTestVector(5, 5)
	fmt.Println(wTrue)
	points := machine_learning.GenerateTestData(5, 5, 10000, wTrue)
	gd := machine_learning.NewGradientDescent(points, func(w vector.Vector, p *machine_learning.Point) vector.Vector { derevative := p.X.Clone(); derevative.Scale(2 * (helper.Float64(vector.Dot(w, p.X)) - p.Y)); return derevative }, 0.01, 10000)	
	w := gd.Run(5)
	fmt.Println(w)
	sgd := machine_learning.NewStochasticGradientDescent(points, func(w vector.Vector, p *machine_learning.Point) vector.Vector { derevative := p.X.Clone(); derevative.Scale(2 * (helper.Float64(vector.Dot(w, p.X)) - p.Y)); return derevative }, 0.1, 10000, len(points))	
	w := sgd.Run(5)
	fmt.Println(w) */
	
	/* trainInput, trainResult := helper.ReadCSV("algorithms/machinelearning/data/iris_train.csv", 7, []int{4, 5, 6}, 4, 3)
	testInput, testResult := helper.ReadCSV("algorithms/machinelearning/data/iris_test.csv", 7, []int{4, 5, 6}, 4, 3)
	config := machine_learning.NewNeuralNetworkConfig(4, 3, 3, 225, 0.1, helper.CrossEntropy, helper.Sigmoid, helper.DSigmoid, helper.SoftMax, helper.ArgMax)
	nn := machine_learning.NewNeuralNetwork(config)
	maxAccuracy := 0.0
	for i := 0; i <= 100; i++ {
		nn.Train(trainInput, trainResult)
		prediction := nn.Predict(testInput)
		accuracy := helper.Accuracy(prediction, testResult)
		maxAccuracy = math.Max(accuracy, maxAccuracy)
	}
	
	fmt.Printf("Accuracy = %0.2f\n", maxAccuracy) */

	/* a := mat.NewDense(2, 1, nil)
	a.Set(0, 0, 0)
	a.Set(1, 0, 1)
	b := mat.NewDense(2, 1, nil)
	b.Set(0, 0, 1)
	b.Set(1, 0, 2)
	c := mat.NewDense(2, 1, nil)
	c.Set(0, 0, 1)
	c.Set(1, 0, 3)
	d := mat.NewDense(2, 1, nil)
	d.Set(0, 0, 8)
	d.Set(1, 0, 4)
	e := mat.NewDense(2, 1, nil)
	e.Set(0, 0, 9)
	e.Set(1, 0, 3)
	f := mat.NewDense(2, 1, nil)
	f.Set(0, 0, 8)
	f.Set(1, 0, 2)
	g := mat.NewDense(2, 1, nil)
	g.Set(0, 0, 8)
	g.Set(1, 0, 8)
	h := mat.NewDense(2, 1, nil)
	h.Set(0, 0, 8)
	h.Set(1, 0, 6)
	i := mat.NewDense(2, 1, nil)
	i.Set(0, 0, 5)
	i.Set(1, 0, 8)
	j := mat.NewDense(2, 1, nil)
	j.Set(0, 0, 7)
	j.Set(1, 0, 7)
	k := mat.NewDense(2, 1, nil)
	k.Set(0, 0, 4)
	k.Set(1, 0, 7)
	l := mat.NewDense(2, 1, nil)
	l.Set(0, 0, 3)
	l.Set(1, 0, 8)
	m := mat.NewDense(2, 1, nil)
	m.Set(0, 0, 2)
	m.Set(1, 0, 9)
	n := mat.NewDense(2, 1, nil)
	n.Set(0, 0, 6)
	n.Set(1, 0, 9) */

	/* a := mat.NewDense(1, 1, nil)
	a.Set(0, 0, 0)
	b := mat.NewDense(1, 1, nil)
	b.Set(0, 0, 2)
	c := mat.NewDense(1, 1, nil)
	c.Set(0, 0, 10)
	d := mat.NewDense(1, 1, nil)
	d.Set(0, 0, 12) */

	/* x := []*mat.Dense {a, b, c, d, e, f, g, h, i, j, k, l, m, n}
	km := machine_learning.NewKMeans(3, 1000)
	minLoss := math.MaxFloat64
	for i := 0; i < 100; i++ {
		_, loss := km.Run(x)
		if loss < minLoss {
			minLoss = loss
		}
	}
	fmt.Println(minLoss) */

	/* test := []float64{
		0, 1, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 16,
		0, 0, 1, 0, 5, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 17,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 5, 6, 0, 0, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 20,
		0, 0, 7, 0, 0, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 0, 21,
		0, 7, 0, 0, 0, 0, 0, 0, 6, 0, 0, 0, 0, 0, 0, 0, 22,
		0, 0, 0, 0, 0, 0, 0, 0, 7, 0, 0, 0, 0, 0, 0, 0, 23,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 17, 16, 0, 20, 0, 0, 0, 24, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 17, 0, 21, 0, 0, 0, 25, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 21, 22, 0, 0, 0, 0, 0, 28, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 23, 0, 0, 0, 0, 0, 29, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 23, 0, 0, 0, 0, 0, 0, 30, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 31, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 25, 24, 0, 28, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 25, 0, 29, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 29, 30, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 31, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 31, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	} */

	/* test := []float64 {
		0, 1, 0, 0, 4,
		0, 0, 1, 0, 5,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		0, 5, 6, 0, 0,
		0, 0, 7, 0, 0,
		0, 7, 0, 0, 0,
		0, 0, 0, 0, 0,
	} */

	/* test := []float64 {
		0, 1, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 16,
		0, 0, 1, 0, 5, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 17,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 5, 4, 0, 0, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 20,
		0, 0, 5, 0, 0, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 0, 21,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 17, 18, 0, 20, 0, 0, 0, 16, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 19, 0, 21, 0, 0, 0, 17, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 19, 0, 0, 22, 0, 0, 0, 18, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 23, 0, 0, 0, 19, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 21, 22, 0, 0, 0, 0, 0, 20, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 23, 0, 0, 0, 0, 0, 21, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 23, 0, 0, 0, 0, 0, 0, 30, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 31, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 31, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	} */

	/* test := []float64 {
		0, 1, 0, 0, 4, 0, 0, 0, 8, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 1, 0, 5, 0, 0, 0, 9, 0, 0, 0, 0, 0, 0, 0, 1,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 5, 4, 0, 0, 0, 0, 0, 12, 0, 0, 0, 0, 0, 0, 0, 4,
		0, 0, 5, 0, 0, 0, 0, 0, 13, 0, 0, 0, 0, 0, 0, 0, 5,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 9, 10, 0, 12, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 8,
		0, 0, 11, 0, 13, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 9,
		0, 11, 0, 0, 14, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 10,
		0, 0, 0, 0, 15, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 11,
		0, 13, 14, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 12,
		0, 0, 15, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 13,
		0, 15, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 30,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 31,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 31, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	} */

	/* ansible := iac.NewAnsible("hosts", "yuma.yml", "roles/")
	yuma := search.NewYuma(ansible, nil, nil)
	yuma.BuildSearchTree(3, 30, 0, 0, make([]string, 0))
	searchTree := yuma.GetSearchTree()
	f := mat.Formatted(searchTree, mat.Prefix("             "), mat.Squeeze())
	printedSearchTree := fmt.Sprintf("\nSearchTree = %v\n\n\n", f)
	//minDepth, minPath := yuma.DetermineExecutionOrder(0, 0, make([]string, 0), 16, make(map[int]int, 0), make(map[int][]string, 0))
	minDepth, minPath := yuma.DetermineExecutionOrder_UniformCostSearch(2)

	exportPlaybook := yuma.CreateDeploymentPlan("yuma", minPath)

	exportResults := "\n"
	exportResults += fmt.Sprint(yuma.GetRoles())
	exportResults += "\n" + printedSearchTree
	exportResults += fmt.Sprint(minDepth) + "\n"
	exportResults += fmt.Sprint(minPath) + "\n\n"

	file, err := os.OpenFile("playbook.yml", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
	}

	defer file.Close()

	if _, err := file.WriteString(exportPlaybook); err != nil {
		fmt.Println(err)
	}

	file, err = os.OpenFile("result.txt", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
	}

	defer file.Close()

	if _, err := file.WriteString(exportResults); err != nil {
		fmt.Println(err)
	} */

	ansible := iac.NewAnsible("hosts", "yuma.yml", "roles/")
	yuma := search.NewYuma(ansible, nil, nil)
	if err := yuma.LearnActionValues_QLearning(2, 0.5, 1.0, 0.1, 3, 30, 50); err != nil {
		fmt.Println(err)
	}
	if err := yuma.DeriveOptimalPolicy(2); err != nil {
		fmt.Println(err)
	}
}