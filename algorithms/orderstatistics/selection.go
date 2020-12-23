package orderstatistics

import (
	"math/rand"
	"time"
)

type Selection struct {
}

func NewSelection() *Selection {
	selection := new(Selection)
	return selection
}

func (selection *Selection) Minimum(nums []int) int {
	min := nums[0]

	for i := 1; i < len(nums); i++ {
		if nums[i] < min {
			min = nums[i]
		}
	}

	return min
}

func (selection *Selection) Maximum(nums []int) int {
	max := nums[0]

	for i := 1; i < len(nums); i++ {
		if nums[i] > max {
			max = nums[i]
		}
	}

	return max
}

func (selection *Selection) MinimumMaximum(nums []int) (int, int) {
	var min int
	var max int
	odd := len(nums) % 2

	if odd > 0 {
		min = nums[0]
		max = nums[0]
	} else {
		if nums[0] <= nums[1] {
			min = nums[0]
			max = nums[1]
		} else {
			min = nums[1]
			max = nums[0]
		}
	}

	if odd > 0 {
		for i := 1; i < len(nums) - 1; i+=2 {
			var tmpMin int
			var tmpMax int

			if nums[i] < nums[i + 1] {
				tmpMin = nums[i]
				tmpMax = nums[i + 1]
			} else {
				tmpMin = nums[i + 1]
				tmpMax = nums[i]
			}

			if tmpMin < min {
				min = tmpMin
			}
			if tmpMax > max {
				max = tmpMax
			}
		}
	} else {
		for i := 2; i < len(nums) - 1; i+=2 {
			var tmpMin int
			var tmpMax int

			if nums[i] < nums[i + 1] {
				tmpMin = nums[i]
				tmpMax = nums[i + 1]
			} else {
				tmpMin = nums[i + 1]
				tmpMax = nums[i]
			}

			if tmpMin < min {
				min = tmpMin
			}
			if tmpMax > max {
				max = tmpMax
			}
		}
	}

	return min, max
}

func (selection *Selection) RandomSelect(nums []int, p, r, i int) int {
	if p >= r {
		return nums[p]
	}

	q := selection.randomPartition(nums, p, r)
	k := q - p + 1

	if i == k {
		return nums[q]
	} else if i < k {
		return selection.RandomSelect(nums, p, q - 1, i)
	} else {
		return selection.RandomSelect(nums, q + 1, r, i - k)
	}
}

func (selection *Selection) randomPartition(numbers []int, p, r int) int {
	rand.Seed(time.Now().UnixNano())
	i := r - p + 1
	i = rand.Intn(r - p + 1) + p

	tmp := numbers[r]
	numbers[r] = numbers[i]
	numbers[i] = tmp

	return selection.partition(numbers, p, r)
}

func (selection *Selection) partition(numbers []int, p, r int) int {
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
	numbers[i+ 1] = numbers[r]
	numbers[r] = tmp

	return i + 1
}

func (selection *Selection) Select(nums []int, p, r, i int) int {
	if p >= r {
		return nums[p]
	}

	x := selection.median(nums, p, r)
	q := selection.modifiedPartition(nums, p, r, x)
	k := q - p + 1

	if i == k {
		return nums[q]
	} else if i < k {
		return selection.Select(nums, p, q - 1, i)
	} else {
		return selection.Select(nums, q + 1, r, i - k)
	}
}

func (selection *Selection) median(nums []int, p, r int) int {
	size := 5
	if r - p + 1 <= size {
		result := selection.insertionSort(nums[p:r])
		median := (len(result) - 1) / 2
		return result[median]
	}

	if p >= r {
		return nums[p]
	}

	medians := make([]int, 0)
	for i := p; i <= r; i += size {
		end := i + size - 1

		if end > r {
			end = r
		}

		medians = append(medians, selection.median(nums, i, end))
	}

	return selection.median(medians, 0, len(medians) - 1)
}

func (selection *Selection) insertionSort(numbers []int) []int {
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

func (selection *Selection) modifiedPartition(numbers []int, p, r, x int) int {
	i := p - 1

	for j := p; j < r; j++ {
		if numbers[j] == x {
			tmp := numbers[r]
			numbers[r] = x
			numbers[j] = tmp
		}

		if numbers[j] <= x {
			i++
			tmp := numbers[i]
			numbers[i] = numbers[j]
			numbers[j] = tmp
		}
	}

	tmp := numbers[i + 1]
	numbers[i + 1] = x
	numbers[r] = tmp

	return i + 1
}