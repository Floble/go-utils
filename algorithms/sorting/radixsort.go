package sorting

import (
	"math"
)

type RadixSort struct {
}

func NewRadixSort() *RadixSort {
	radixSort := new(RadixSort)
	return radixSort
}

func (radixSort *RadixSort) Sort(numbers []int, d int) []int {
	result := numbers
	
	for i := 1; i <= d; i++ {
		k := radixSort.findMax(numbers, i)
		result = radixSort.countingSort(result, k, i)
	}

	return result
}

func (radixSort *RadixSort) findMax(numbers []int, i int) int {
	max := math.MinInt64
	
	for _, number := range numbers {
		n := radixSort.getDigit(number, i)
		if n > max {
			max = n
		}
	}

	return max
}

func (radixSort *RadixSort) getDigit(number int, i int) int {
    r := number % int(math.Pow(10, float64(i)))
    return r / int(math.Pow(10, float64(i-1)))
}

func (radixSort *RadixSort) countingSort(numbers []int, k int, i int) []int {
	result := make([]int, len(numbers) + 1)
	tmp := make([]int, k + 1)

	for i := 0; i <= k; i++ {
		tmp[i] = 0
	}

	for j := 0; j < len(numbers); j++ {
		n := radixSort.getDigit(numbers[j], i)
		tmp[n] += 1
	}

	for i := 1; i <= k; i++ {
		tmp[i] = tmp[i] + tmp[i - 1]
	}

	for j := len(numbers) - 1; j >= 0; j-- {
		n := radixSort.getDigit(numbers[j], i)
		result[tmp[n]] = numbers[j]
		tmp[n] -= 1
	}

	return result[1:]
}