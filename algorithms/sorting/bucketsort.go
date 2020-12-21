package sorting

type BucketSort struct {
}

func NewBucketSort() *BucketSort {
	bucketSort := new(BucketSort)
	return bucketSort
}

func (bucketSort *BucketSort) Sort(numbers []float64) []float64 {
	n := float64(len(numbers))
	
	b := make([][]float64, len(numbers))
	for i := 0; i < len(numbers); i++ {
		b[i] = make([]float64, 0)
	}

	for i := 0; i < len(numbers); i++ {
		b[int(n*numbers[i])] = append(b[int(n*numbers[i])], numbers[i])
	}

	for i := 0; i < len(numbers); i++ {
		b[i] = bucketSort.insertionSort(b[i])
	}

	result := make([]float64, 0)

	for i := 0; i < len(numbers); i++ {
		result = append(result, b[i]...)
	}

	return result
}

func (bucketSort *BucketSort) insertionSort(numbers []float64) []float64 {
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