package datastructures

type Element struct {
	State int
	Cost int
	Predecessor int
}

type MinPriorityQueue struct {
	Queue []*Element
	Elements map[int]int
}

func NewElement(state, cost, predecessor int) *Element {
	e := new(Element)
	e.State = state
	e.Cost = cost
	e.Predecessor = predecessor

	return e
}

func NewMinPriorityQueue() *MinPriorityQueue {
	mpq := new(MinPriorityQueue)
	mpq.Queue = make([]*Element, 0)
	mpq.Elements = make(map[int]int, 0)

	return mpq
}

func (mpq MinPriorityQueue) Len() int { 
	return len(mpq.Queue) 
}

func (mpq MinPriorityQueue) Less(i, j int) bool {
	if mpq.Queue[i].Cost < mpq.Queue[j].Cost {
		return true
	} else {
		return false
	}
}

func (mpq MinPriorityQueue) Swap(i, j int) {
	tmp := mpq.Queue[i]
	mpq.Queue[i] = mpq.Queue[j]
	mpq.Queue[j] = tmp
}

func (mpq *MinPriorityQueue) Pop() interface{} {
	tmp := *mpq

	l := len(mpq.Queue)
	element := tmp.Queue[l - 1]
	mpq.Queue = tmp.Queue[0 : l - 1]

	return element
}

func (mpq *MinPriorityQueue) Push(x interface{}) {
	element := x.(*Element)
	mpq.Queue = append(mpq.Queue, element)
}

func (mpq *MinPriorityQueue) Search(x interface{}) int {
	for i, e := range mpq.Queue {
		if e.State == x {
			return i
		}
	}

	return -1
}