package datastructures

import (
	"errors"
)

type Queue struct {
	Head, Tail int
	Array []int
}

func NewQueue() *Queue {
	queue := new(Queue)
	queue.Array = make([]int, 0)
	queue.Head = 0
	queue.Tail = 0

	return queue
}

func (queue *Queue) IsEmpty() bool {
	if queue.Head == queue.Tail {
		return true
	} else {
		return false
	}
}

func (queue *Queue) Enqueue(element int) {
	queue.Array = append(queue.Array, element)
	queue.Tail++
}

func (queue *Queue) Dequeue() (error, int) {
	if queue.IsEmpty() {
		err := errors.New("Underflow")
		return err, -1
	}
	queue.Head++
	return nil, queue.Array[queue.Head-1]
}