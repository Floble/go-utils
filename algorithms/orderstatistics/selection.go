package orderstatistics

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