package sorting

type CountingSort struct {
}

func NewCountingSort() *CountingSort {
	countingSort := new(CountingSort)
	return countingSort
}

func (countingSort *CountingSort) Sort(numbers []int, k int) []int {
	result := make([]int, len(numbers) + 1)
	tmp := make([]int, k + 1)

	for i := 0; i <= k; i++ {
		tmp[i] = 0
	}

	for j := 0; j < len(numbers); j++ {
		tmp[numbers[j]] += 1
	}

	for i := 1; i <= k; i++ {
		tmp[i] = tmp[i] + tmp[i - 1]
	}

	for j := len(numbers) - 1; j >= 0; j-- {
		result[tmp[numbers[j]]] = numbers[j]
		tmp[numbers[j]] -= 1
	}

	return result
}