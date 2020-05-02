package sorting

type InsertionSort struct {
}

func NewInsertionSort() *InsertionSort {
	insertionSort := new(InsertionSort)
	return insertionSort
}

func (insertionSort *InsertionSort) Sort(numbers []int) []int {
	for i := 1; i < len(numbers); i++ {
		key := numbers[i]
		j := i-1
		for j >= 0 && numbers[j] > key {
			numbers[j+1] = numbers[j]
			j--
		}
		numbers[j+1] = key
	}
	return numbers
}