package datastructures

type Element struct {
	State int
	Cost int
	Predecessor int
}

type MinPriorityQueue []*Element

func NewElement(state, cost, predecessor int) *Element {
	e := new(Element)
	e.State = state
	e.Cost = cost
	e.Predecessor = predecessor

	return e
}

func (mpq MinPriorityQueue) Len() int { 
	return len(mpq) 
}

func (mpq MinPriorityQueue) Less(i, j int) bool {
	if mpq[i].Cost < mpq[j].Cost {
		return true
	} else {
		return false
	}
}

func (mpq MinPriorityQueue) Swap(i, j int) {
	tmp := mpq[i]
	mpq[i] = mpq[j]
	mpq[j] = tmp
}

func (mpq *MinPriorityQueue) Pop() interface{} {
	tmp := *mpq

	l := len(*mpq)
	element := tmp[l - 1]
	*mpq = tmp[0 : l - 1]

	return element
}

func (mpq *MinPriorityQueue) Push(x interface{}) {
	element := x.(*Element)
	*mpq = append(*mpq, element)
}

func (mpq *MinPriorityQueue) Search(x interface{}) int {
	for i, e := range *mpq {
		if e.State == x {
			return i
		}
	}

	return -1
}