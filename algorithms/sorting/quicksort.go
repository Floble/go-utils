package sorting

import (
	"math/rand"
	"time"
)

type QuickSort struct {
}

func NewQuickSort() *QuickSort {
	quicksort := new(QuickSort)
	return quicksort
}

func (qs *QuickSort) Sort(numbers []int) {
	quicksort(numbers, 0, len(numbers) - 1)
}

func quicksort(numbers []int, p int, r int) {
	if p < r {
		q := randomPartition(numbers, p, r)
		quicksort(numbers, p, q - 1)
		quicksort(numbers, q + 1, r)
	}
}

func randomPartition(numbers []int, p int, r int) int {
	rand.Seed(time.Now().UnixNano())
	i := r - p + 1
	i = rand.Intn(r - p + 1) + p

	tmp := numbers[r]
	numbers[r] = numbers[i]
	numbers[i] = tmp

	return partition(numbers, p, r)
}

func partition(numbers []int, p int, r int) int {
	x := numbers[r]
	i := p - 1

	for j := p; j < r; j++ {
		if numbers[j] <= x {
			i++
			tmp := numbers[i]
			numbers[i] = numbers[j]
			numbers[j] = tmp
		}
	}

	tmp := numbers[i + 1]
	numbers[i + 1] = numbers[r]
	numbers[r] = tmp

	return i + 1
}