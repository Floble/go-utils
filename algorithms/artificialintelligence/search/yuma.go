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

/* 
Definition of the YUMA structure
Attributes:
- roles: The roles contained in a repository
- configurations: The binary representations for the roles
- searchTree: The model of YUMA (i.e., a search tree) represented as a matrix
- ansible: The executor for invoking Ansible commands
*/
type Yuma struct {
	roles map[string]int
	configurations map[int]string
	searchTree *mat.Dense
	ansible *iac.Ansible
}

/*
The constructor for YUMA
Parameters:
- ansible: The executor that must be used by the object
- searchTreeData: An existing search tree that must be used as the model of the object
Returns:
- *Yuma: An object of the YUMA structure
*/
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

/*
Returns the state from which the search algorithm must start
Returns:
- 0
*/
func (yuma *Yuma) startState() int {
	return 0
}

/*
Returns the actions that can be performed at a specific state
Parameters:
- state: The state for which the possible actions must be returned
Returns:
- actions: The set of actions that can be performed at the given state
*/
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

/*
Returns the successor state when performing a specific action at a specific state
Parameters:
- state: The state at which a specific action must be performed
- action: The action that must be performed at that state
Returns:
- The successor state as a result of a lookup to the model
*/
func (yuma *Yuma) successor(state, action int) int {
	return int(yuma.searchTree.At(state, action))
}

/*
Determines whether a given state is the terminal state
Parameters:
- state: The state that must be checked
Returns:
- true, if the state is a terminal state; false, otherwise
*/
func (yuma *Yuma) isEnd(state int) bool {
	if state == int(math.Exp2(float64(len(yuma.roles)))) - 1 {
		return true
	} else {
		return false
	}
}

/*
Determines all the Ansible roles stored in a specified directory as well as the corresponding binary representations
Parameters:
- path: The path to the directory storing all the Ansible roles
Returns:
- nil, or any error that occcurred
*/
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

/*
A backtracking search algorithm determining the search tree as model of YUMA
Parameters:
- sigma: The number of iterations of the algorithm to overcome potential false positives
- omega: Number of seconds to wait for the creation of EC2 instances
- state: The current state that is examined by the algorithm
- depth: The current depth of the algorithm in the search tree
- path: The order of actions that were performed previosly to come to the current state
Returns:
- searchTree: The search tree that was determined by the algorithm
*/
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

/*
A dynamic programming algorithm to determine the shortest path from the start state to the specified target state
Parameters:
- state: The current state that is examined by the algorithm
- depth: The current depth of the algorithm in the search tree
- path: The path that has been explored so far by the algorhtm
- memDepth: The determined minimum detph required to reach the target state from the current state
- memPath: The determined minimum path to reach the target state from the current state
Returns:
- memDepth: The minimum depth required to reach the target state from the start state
- memPath: The minimum path to reach the target state from the start state
*/
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

/*
The implemented Uniform Cost Search to determine the shortest path from the start state to the target state
Parameters:
- target: The target state to which the minimum path must be determined
Returns:
- The set of explored states according to the UCS
*/
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

/*
Formatting the set of explored states to a minimal execution order
Parameters:
- explored: The set of explored states as output from the UCS
Returns:
- minCost: The minimum cost required to reach the target state from the start state
- executionOrder: The minimum execution order formatted from the set of explored states
*/
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

/*
Create the composition plan as a Ansible playbook
Parameters:
- hosts: The hosts on which the Ansible roles must be determined. This is just a placeholder
- path: The determined minimum execution order representing the composition plan
Returns:
- export: The Ansible playbook in the form of a string that can be wrote into a file
*/
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

/*
Returns all the roles stored in the specified directory
Returns:
- The roles stored in the specified directory
*/
func (yuma *Yuma) GetRoles() map[string]int {
	return yuma.roles
}

/*
Returns the binary representations of the Ansible roles
Returns:
- The binary representations of the Ansible roles
*/
func (yuma *Yuma) GetConfigurations() map[int]string {
	return yuma.configurations
}

/*
Returns the model of YUMA
Returns:
- The search tree that is hold by YUMA
*/
func (yuma *Yuma) GetSearchTree() *mat.Dense {
	return yuma.searchTree
}

/*
Cleans the array holding the visited path
Parameters:
- path: The array to be cleaned
- depth: The depth to which the array must be cleaned
*/
func clearPath(path []string, depth int) {
	for i := depth; i < len(path); i++ {
		path[i] = ""
	}
}