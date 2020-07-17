package string_matching

import (
	"fmt"
)

type KMP struct {
}

func NewKMP() *KMP {
	kmp := new(KMP)
	return kmp
}

func (kmp *KMP) Match(s, p string) []int {
	shifts := make([]int, 0)
	j := 0

	prefix := kmp.prefixFunc(p)

	for i := 0; i <= len(s); i++ {
		fmt.Println(j)
		if j == len(p) {
			shifts = append(shifts, i-len(p))
			j = 0
		} else if i < len(s) {
			if string(s[i]) == string(p[j]) {
				j++
			} else {
				if j > 0 {
					j = prefix[j-1]
					i--
				}
			}
		}
	}

	return shifts
}

func (kmp *KMP) prefixFunc(p string) []int {
	prefix := make([]int, len(p))
	i := 0

	for j := 1; j < len(p); j++ {
		if i != j {
			if string(p[i]) != string(p[j]) {
				if i == 0 {
					prefix[j] = 0
				} else {
					i = prefix[i-1]
					j--
				}
			} else {
				prefix[j] = i+1
				i++
			}
		}
	}

	return prefix
}