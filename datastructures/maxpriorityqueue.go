package datastructures

import (
	"errors"
	"math"
)

type MaxPriorityQueue struct {
	Heap *Heap
}

func NewMaxPriorityQueue(nums []int) *MaxPriorityQueue {
	pq := new(MaxPriorityQueue)
	pq.Heap = NewHeap(nums)

	return pq
}

func (pq *MaxPriorityQueue) Maximum() int {
	return pq.Heap.Array[0]
}

func (pq *MaxPriorityQueue) ExtractMax() (error, int) {
	if pq.Heap.Size < 1 {
		err := errors.New("Heap Underflow")
		return err, -1
	}

	max := pq.Heap.Array[0]
	pq.Heap.Array[0] = pq.Heap.Array[pq.Heap.Size - 1]
	pq.Heap.Size--
	pq.Heap.maxHeapify(0) 

	return nil, max
}

func (pq *MaxPriorityQueue) IncreaseKey(i int, key int) error {
	if key < pq.Heap.Array[i] {
		err := errors.New("New Key is smaller than current key")
		return err
	}

	pq.Heap.Array[i] = key

	for i > 0 && pq.Heap.Array[pq.Heap.Parent(i)] < pq.Heap.Array[i] {
		tmp := pq.Heap.Array[i]
		pq.Heap.Array[i] = pq.Heap.Array[pq.Heap.Parent(i)]
		pq.Heap.Array[pq.Heap.Parent(i)] = tmp
		i = pq.Heap.Parent(i)
	}

	return nil
}

func (pq *MaxPriorityQueue) Insert(key int) {
	pq.Heap.Size++
	pq.Heap.Array = append(pq.Heap.Array, math.MinInt64)
	pq.IncreaseKey(pq.Heap.Size - 1, key)
}