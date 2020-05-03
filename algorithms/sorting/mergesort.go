package sorting

type MergeSort struct {
}

func NewMergeSort() *MergeSort {
	mergeSort := new(MergeSort)
	return mergeSort
}

func (mergeSort *MergeSort) Sort(numbers []int, p, r int) []int {
	if p < r {
		var q int
		q = (p+r+1)/2
		q -= 1

		numbers = mergeSort.Sort(numbers, p, q)
		numbers = mergeSort.Sort(numbers, q+1, r)
		numbers = merge(numbers, p, q, r)
	}

	return numbers
}

func merge(numbers []int, p, q, r int) []int {
	n1 := q-p
	n2 := r-q-1

	left := make([]int, 0)
	right := make([]int, 0)

	for i := 0; i <= n1; i++ {
		left = append(left, numbers[p+i])
	}
	for j := 0; j <= n2; j++ {
		right = append(right, numbers[q+1+j])
	}

	i := 0
	j := 0

	for k := p; k <= r; k++ {
		if i < len(left) && j < len(right) {
			if left[i] <= right[j] {
				numbers[k] = left[i]
				i++
				continue
			} else {
				numbers[k] = right[j]
				j++
				continue
			}
		}
		if i >= len(left) && j < len(right) {
			numbers[k] = right[j]
			j++
			continue
		}
		if j >= len(right) && i < len(left) {
			numbers[k] = left[i]

			i++
			continue
		}
	}

	return numbers
}