package search

import (
	"container/heap"
	"go-utils/cloud/aws/ec2"
	ds "go-utils/datastructures"
	iac "go-utils/infrastructureascode"
	"io/ioutil"
	"math"
	"time"
	"gonum.org/v1/gonum/mat"
)

type Yuma struct {
	roles map[string]int
	configurations map[int]string
	searchTree *mat.Dense
	ansible *iac.Ansible
}

func NewYuma(ansible *iac.Ansible, searchTreeData []float64) *Yuma {
	yuma := new(Yuma)
	yuma.roles = make(map[string]int, 0)
	yuma.configurations = make(map[int]string, 0)
	yuma.ansible = ansible
	if err := yuma.identifyRoles(yuma.ansible.GetRepository()); err != nil {
		return nil
	}
	yuma.searchTree = mat.NewDense(int(math.Exp2(float64(len(yuma.roles)))), int(math.Exp2(float64(len(yuma.roles) - 1)) + 1), searchTreeData)

	return yuma
}

func (yuma *Yuma) startState() int {
	return 0
}

func (yuma *Yuma) actions(state int) []int {
	actions := make([]int, 0)

	for _, action := range yuma.roles {
		if state & action != 0 {
			continue
		}

		actions = append(actions, action)
	}

	return actions
}

func (yuma *Yuma) successor(state, action int) int {
	return int(yuma.searchTree.At(state, action))
}

func (yuma *Yuma) isEnd(state int) bool {
	if state == int(math.Exp2(float64(len(yuma.roles)))) - 1 {
		return true
	} else {
		return false
	}
}

func (yuma *Yuma) identifyRoles(path string) error {
	roles, err := ioutil.ReadDir(path)
    if err != nil {
        return err
    }

	for i := 0; i < len(roles); i++ {
		var config int = yuma.startState()
		config |= int(math.Exp2(float64(i)))
		yuma.roles[roles[i].Name()] = config
		yuma.configurations[config] = roles[i].Name()
	}
	
	return nil
}

func (yuma *Yuma) BuildSearchTree(sigma int, omega int, state int, depth int, path []string) *mat.Dense {
	if yuma.isEnd(state) {
		return yuma.searchTree
	}

	for _, action := range yuma.actions(state) {
		for i := 0; i < sigma; i++ {
			instance := ec2.NewEC2Instance()
			err := instance.Create()
			if err != nil {
				return nil
			}

			yuma.ansible.CreateInventory(instance.GetPublicIP())
			time.Sleep(time.Duration(omega) * time.Second)
			instance.AddToKnownHosts()
			if len(path) > 0 && !yuma.ansible.PlayRoles(path, "install") {
				instance.Delete()
				yuma.ansible.DeleteInventory()
				continue
			}

			role := make([]string, 0)
			role = append(role, yuma.configurations[action])

			if yuma.ansible.PlayRoles(role, "install") {
				instance.Delete()
				yuma.ansible.DeleteInventory()
				yuma.searchTree.Set(state, action, float64(state | action))
				break
			} else {
				instance.Delete()
				yuma.ansible.DeleteInventory()
				yuma.searchTree.Set(state, action, float64(state))
			}
		}

		successor := yuma.successor(state, action)
		if successor == state | action {
			path = append(path, yuma.configurations[action])
			yuma.BuildSearchTree(sigma, omega, successor, depth + 1, path)
			path = path[:depth]
		}
	}

	return yuma.searchTree
}

func (yuma *Yuma) DetermineExecutionOrder_DynamicProgramming(state int, depth int, path []string, target int, memDepth map[int]int, memPath map[int][]string) (int, []string) {
	if _, ok := memDepth[state]; ok {
		return memDepth[state], memPath[state]
	}
	
	if state & target != 0 {
		memDepth[state] = depth
		memPath[state] = make([]string, len(path))
		copy(memPath[state], path)
		return memDepth[state], memPath[state]
	}

	minDepth := len(yuma.roles)
	minPath := make([]string, len(yuma.roles))
	for _, action := range yuma.actions(state) {
		successor := yuma.successor(state, action)
		if successor == state | action {
			path = append(path, yuma.configurations[action])

			tmpDepth, tmpPath := yuma.DetermineExecutionOrder_DynamicProgramming(successor, depth + 1, path, target, memDepth, memPath)
			if tmpDepth <= minDepth {
				minDepth = tmpDepth
				minPath = minPath[:len(tmpPath)]
				copy(minPath, tmpPath)
			}

			memDepth[state] = minDepth
			memPath[state] = minPath

			path = path[:depth]
		}
	}

	return memDepth[state], memPath[state]
}

func (yuma *Yuma) DetermineExecutionOrder_UniformCostSearch(target int) (int, []string) {
	frontier := ds.NewMinPriorityQueue()
	heap.Init(frontier)
	explored := make([]*ds.Element, 0)
	exploredElements := make(map[int]int, 0)

	startElement := ds.NewElement(yuma.startState(), 0, 0)
	heap.Push(frontier, startElement)

	for frontier.Len() > 0 {
		state := heap.Pop(frontier).(*ds.Element)
		explored = append(explored, state)
		exploredElements[state.State] = len(explored) - 1
		if state.State & target != 0 {
			return yuma.FormatExecutionOrder_UniformCostSearch(explored)
		}

		for _, action := range yuma.actions(state.State) {
			successor := yuma.successor(state.State, action)
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

	return yuma.FormatExecutionOrder_UniformCostSearch(explored)
}

func (yuma *Yuma) FormatExecutionOrder_UniformCostSearch(explored []*ds.Element) (int, []string) {
	executionOrder := make([]string, explored[len(explored) - 1].Cost)
	e := explored[len(explored) - 1]
	minCost := e.Cost
	executionOrder[len(executionOrder) - 1] = yuma.configurations[e.State &^ e.Predecessor]
	j := len(executionOrder) - 2
	predecessor := e.Predecessor

	for i := len(explored) - 2; i >= 0; i-- {
		e = explored[i]
		if e.State == predecessor && e.State != yuma.startState() {
			executionOrder[j] = yuma.configurations[e.State &^ e.Predecessor]
			j -= 1
			predecessor = e.Predecessor
		}
	}

	return minCost, executionOrder
}

func (yuma *Yuma) CreateDeploymentPlan(hosts string, path []string) string {
	export := ""
	export += "---\n"
	export += "- hosts: "
	export += hosts + "\n"
	export += "  roles:\n"

	for _, software := range path {
		export += "    - " + software + "\n"
	}

	return export
}

func (yuma *Yuma) GetRoles() map[string]int {
	return yuma.roles
}

func (yuma *Yuma) GetConfigurations() map[int]string {
	return yuma.configurations
}

func (yuma *Yuma) GetSearchTree() *mat.Dense {
	return yuma.searchTree
}

func clearPath(path []string, depth int) {
	for i := depth; i < len(path); i++ {
		path[i] = ""
	}
}