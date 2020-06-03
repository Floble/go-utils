package datastructures

import (
	"errors"
)

type Stack struct {
	Top int
	Array []int
}

func NewStack() *Stack {
	stack := new(Stack)
	stack.Top = 0
	stack.Array = make([]int, 0)

	return stack
}

func (stack *Stack) IsEmpty() bool {
	if stack.Top == 0 {
		return true
	}
	return false
}

func (stack *Stack) Push(e int) {
	stack.Array = append(stack.Array, e)
	stack.Top++
}

func (stack *Stack) Pop() (error, int) {
	if stack.IsEmpty() {
		err := errors.New("Underflow")
		return err, -1
	}
	stack.Top--
	return nil, stack.Array[stack.Top]
}