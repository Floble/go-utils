package yuma

import (
	"fmt"
	"math"
	"os"
	"container/heap"
	"gonum.org/v1/gonum/mat"
	ds "go-utils/datastructures"
)

type Search struct {
	yuma *Yuma
}

func NewSearch(yuma *Yuma) *Search {
	search := new(Search)
	search.yuma = yuma

	return search
}

func (search *Search) GetYuma() *Yuma {
	return search.yuma
}

func (search *Search) Learn(terminal int) error {
	search.GetYuma().SetModel(0, mat.NewDense(int(math.Exp2(float64(len(search.GetYuma().GetSubprocesses())))), int(math.Exp2(float64(len(search.GetYuma().GetSubprocesses()) - 1)) + 1), nil))
	if err := search.BuildSearchTree(search.GetYuma().GetStartState(), 0, make([]string, 0)); err != nil {
		return err
	}

	exportResults := "++++++++++++++++++++++++++\n"
	exportResults += fmt.Sprintf("Completed Learning\n")
	exportResults += "++++++++++++++++++++++++++\n\n"
	if err := search.log(exportResults, "results.txt"); err != nil {
		return err
	}

	return nil
}

func (search *Search) Solve(target int) []string {
	frontier := ds.NewMinPriorityQueue()
	heap.Init(frontier)
	explored := make([]*ds.Element, 0)
	exploredElements := make(map[int]int, 0)

	startElement := ds.NewElement(search.GetYuma().GetStartState(), 0, 0)
	heap.Push(frontier, startElement)

	for frontier.Len() > 0 {
		state := heap.Pop(frontier).(*ds.Element)
		explored = append(explored, state)
		exploredElements[state.State] = len(explored) - 1
		if state.State & target != 0 {
			return search.FormatExecutionOrder(explored)
		}

		for _, action := range search.GetYuma().Actions(state.State) {
			successor := int(search.GetYuma().GetModel(0).At(state.State, action))
			if _, ok := exploredElements[successor]; ok {
				continue
			}
			successorElement := ds.NewElement(successor, state.Cost + 1, state.State)
			index := frontier.Search(successor)
			if index != -1 {
				heap.Remove(frontier, index)
				delete(frontier.Elements, successor)
			}
			heap.Push(frontier, successorElement)
		}
	}

	return search.FormatExecutionOrder(explored)
}

func (search *Search) BuildSearchTree(state int, depth int, path []string) error {
	if search.GetYuma().IsTerminal(state) {
		exportResults := "++++++++++++++++++++++++++\n"
		exportResults += "Terminal state: "
		exportResults += fmt.Sprintln(path)
		exportResults += "++++++++++++++++++++++++++\n\n"
		if err := search.log(exportResults, "results.txt"); err != nil {
			return err
		}

		return nil
	}

	for _, action := range search.GetYuma().Actions(state) {
		exportResults := "++++++++++++++++++++++++++\n"
		exportResults += "State: "
		exportResults += fmt.Sprintln(path)
		exportResults += "++++++++++++++++++++++++++\n\n"
		if err := search.log(exportResults, "results.txt"); err != nil {
			return err
		}
		exportResults = ""

		exportResults = fmt.Sprintf("Action: %d\n", action)
		if err := search.log(exportResults, "results.txt"); err != nil {
			return err
		}
		exportResults = ""

		switch search.GetYuma().GetMode() {
		case 0:
			if len(search.GetYuma().GetEnvironment().GetInstances(0)[0]) >= 1 {
				search.GetYuma().GetEnvironment().DeleteInstance(0, 0)
			}
		case 1:
			for i := 0; i < len(search.GetYuma().GetSubprocesses()); i++ {
				if instances, ok := search.GetYuma().GetEnvironment().GetInstances(0)[int(math.Exp2(float64(i)))]; ok {
					for len(instances) > 0 {
						search.GetYuma().GetEnvironment().DeleteInstance(0, int(math.Exp2(float64(i))))
					}
				}
			}
		}

		err, success, _, _ := search.GetYuma().GetEnvironment().TakeAction(0, state, action, path, false)
		if err != nil {
			return err
		}

		if success {
			search.GetYuma().GetModel(0).Set(state, action, float64(state | action))
		} else {
			search.GetYuma().GetModel(0).Set(state, action, float64(state))
		}

		successor := int(search.GetYuma().GetModel(0).At(state, action))
		exportResults = fmt.Sprintf("Successor: %d\n\n", successor)
		if err := search.log(exportResults, "results.txt"); err != nil {
			return err
		}
		exportResults = ""

		if successor == state | action {
			path = append(path, search.GetYuma().GetConfigurations()[action])
			search.BuildSearchTree(successor, depth + 1, path)
			path = path[:depth]
		}
	}

	return nil
}

func (search *Search) FormatExecutionOrder(explored []*ds.Element) []string {
	executionOrder := make([]string, explored[len(explored) - 1].Cost)
	e := explored[len(explored) - 1]
	executionOrder[len(executionOrder) - 1] = search.GetYuma().GetConfigurations()[e.State &^ e.Predecessor]
	j := len(executionOrder) - 2
	predecessor := e.Predecessor

	for i := len(explored) - 2; i >= 0; i-- {
		e = explored[i]
		if e.State == predecessor && e.State != search.GetYuma().GetStartState() {
			executionOrder[j] = search.GetYuma().GetConfigurations()[e.State &^ e.Predecessor]
			j -= 1
			predecessor = e.Predecessor
		}
	}

	return executionOrder
}

func (search *Search) ArgMaxAction (q *mat.Dense, state int) int {
	return 0
}

func (search *Search) log(exportResults string, filePath string) error {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	if _, err := file.WriteString(exportResults); err != nil {
		return err
	}

	defer file.Close()

	return nil
}