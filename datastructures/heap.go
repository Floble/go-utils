package datastructures

type Heap struct {
	Size int
	Array []int
}

func NewHeap(numbers []int) *Heap {
	heap := new(Heap)
	heap.Size = 0
	heap.Array = numbers

	return heap
}

func (heap *Heap) Parent(i int) int {
	return i / 2
}

func (heap *Heap) Left(i int) int {
	return 2 * i
}

func (heap *Heap) Right(i int) int {
	return 2 * i + 1
}

func (heap *Heap) maxHeapify(i int) {
	l := heap.Left(i)
	r := heap.Right(i)
	var largest int

	if l <= heap.Size - 1 && heap.Array[l] > heap.Array[i] {
		largest = l
	} else {
		largest = i
	}

	if r <= heap.Size - 1 && heap.Array[r] > heap.Array[largest] {
		largest = r
	}

	if largest != i {
		tmp := heap.Array[i]
		heap.Array[i] = heap.Array[largest]
		heap.Array[largest] = tmp
		heap.maxHeapify(largest)
	}
}

func (heap *Heap) BuildMaxHeap() {
	heap.Size = len(heap.Array)

	for i := len(heap.Array) / 2; i >= 0; i-- {
		heap.maxHeapify(i)
	}
}

func (heap *Heap) Sort() {
	heap.BuildMaxHeap()

	for i := len(heap.Array) - 1; i >= 1; i-- {
		tmp := heap.Array[i]
		heap.Array[i] = heap.Array[0]
		heap.Array[0] = tmp

		heap.Size--
		heap.maxHeapify(0)
	}
}