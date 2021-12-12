package search

import (
	"container/heap"
	"fmt"
	"os"
	"go-utils/cloud/aws/ec2"
	ds "go-utils/datastructures"
	iac "go-utils/infrastructureascode"
	"io/ioutil"
	"math"
	"math/rand"
	"time"
	"github.com/jmcvetta/randutil"
	"gonum.org/v1/gonum/mat"
)

type Yuma struct {
	roles map[string]int
	configurations map[int]string
	searchTree *mat.Dense
	q *mat.Dense
	policy_EpsilonGreedy map[int][]randutil.Choice
	ansible *iac.Ansible
}

func NewYuma(ansible *iac.Ansible, searchTreeData []float64, qData []float64) *Yuma {
	yuma := new(Yuma)
	yuma.roles = make(map[string]int, 0)
	yuma.configurations = make(map[int]string, 0)
	yuma.ansible = ansible
	if err := yuma.identifyRoles(yuma.ansible.GetRepository()); err != nil {
		return nil
	}
	yuma.searchTree = mat.NewDense(int(math.Exp2(float64(len(yuma.roles)))), int(math.Exp2(float64(len(yuma.roles) - 1)) + 1), searchTreeData)
	yuma.q = mat.NewDense(int(math.Exp2(float64(len(yuma.roles)))), int(math.Exp2(float64(len(yuma.roles) - 1)) + 1), searchTreeData)
	yuma.policy_EpsilonGreedy = make(map[int][]randutil.Choice, 0)

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

func (yuma *Yuma) isTerminal(state int) bool {
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

func (yuma *Yuma) BuildSearchTree(sigma int, omega int, state int, depth int, path []string) (error, *mat.Dense) {
	if yuma.isTerminal(state) {
		return nil, yuma.searchTree
	}

	for _, action := range yuma.actions(state) {
		for i := 0; i < sigma; i++ {
			instance := ec2.NewEC2Instance()
			err := instance.Create()
			if err != nil {
				return err, nil
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

	return nil, yuma.searchTree
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

func (yuma *Yuma) CreateDeploymentPlan(hosts string, path []string) error {
	export := ""
	export += "---\n"
	export += "- hosts: "
	export += hosts + "\n"
	export += "  roles:\n"

	for _, software := range path {
		export += "    - " + software + "\n"
	}

	return yuma.log(export, "playbook.yml")
}

func (yuma *Yuma) LearnActionValues_QLearning(target int, alpha float64, gamma float64, epsilon float64, sigma int, omega int, episodes int) error {
	exportResults := fmt.Sprintln(yuma.GetRoles()) + "\n"
	if err := yuma.log(exportResults, "results.txt"); err != nil {
		return err
	}
	exportResults = ""

	// Initialize Q(s, a)
	yuma.initializeQ(target)
	// Repeat (for each episode)
	for i := 0; i < episodes; i++ {
		exportResults += "++++++++++++++++++++++++++\n"
		exportResults += fmt.Sprintf("Episode: %d\n", i + 1)
		exportResults += "++++++++++++++++++++++++++\n\n"
		if err := yuma.log(exportResults, "results.txt"); err != nil {
			return err
		}
		exportResults = ""

		// Initialize S
		state, path := yuma.initializeState()
		// Repeat (for each step of episode until S is target state)
		for state & target == 0 {
			f := mat.Formatted(yuma.q, mat.Prefix("    "), mat.Squeeze())
			exportResults += fmt.Sprintf("\nQ = %v\n\n\n", f)
			if err := yuma.log(exportResults, "results.txt"); err != nil {
				return err
			}
			exportResults = ""

			// Choose A from S using policy derived from Q (e.g., epsilon-greedy)
			yuma.derivePolicy_EpsilonGreedy(epsilon)
			c, _ := randutil.WeightedChoice(yuma.policy_EpsilonGreedy[state])
			action := c.Item.(int)

			exportResults += fmt.Sprintf("State: %d\n", state)
			exportResults += "Path: "
			exportResults += fmt.Sprintln(path)
			exportResults += "Policy: "
			exportResults += fmt.Sprintln(yuma.policy_EpsilonGreedy)
			exportResults += fmt.Sprintf("Action: %d\n", action)
			if err := yuma.log(exportResults, "results.txt"); err != nil {
				return err
			}
			exportResults = ""

			// Take action A, observe R, S'
			err, success, reward, successor := yuma.takeAction(sigma, omega, path, state, action)
			if err != nil {
				return err
			} else if success {
				path = append(path, yuma.configurations[action])
			}

			// Q(S, A) = Q(S, A) + alpha * [reward + gamma * max_Q(S', a) - Q(S, A)]	
			qUpdate := yuma.q.At(state, action) + alpha * (reward + (gamma * yuma.max_Q(successor)) - yuma.q.At(state, action))
			yuma.q.Set(state, action, qUpdate)

			// S = S'
			state = successor

			exportResults += "Reward: " + fmt.Sprintf("%f", reward) + "\n"
			exportResults += "Max_Q_Successor: " + fmt.Sprintf("%f", yuma.max_Q(successor)) + "\n"
			exportResults += "Q_Update: " + fmt.Sprintf("%f", qUpdate) + "\n"
			exportResults += "Successor: " + fmt.Sprintf("%d\n\n", successor)
			if err := yuma.log(exportResults, "results.txt"); err != nil {
				return err
			}
			exportResults = ""
		}
	}

	exportResults += "++++++++++++++++++++++++++\n"
	exportResults += fmt.Sprintf("Completed Learning\n")
	exportResults += "++++++++++++++++++++++++++\n\n"
	if err := yuma.log(exportResults, "results.txt"); err != nil {
		return err
	}

	return nil
}

func (yuma *Yuma) DeriveOptimalPolicy(target int) error {
	state := yuma.startState()
	path := make([]string, 0)

	for state & target == 0 {
		action := yuma.argMax_Action(state)
		role := yuma.configurations[action]
		path = append(path, role)
		state = state | action
	}

	return yuma.CreateDeploymentPlan("Yuma", path)
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

func (yuma *Yuma) initializeQ(target int) {
	for i := 0; i < int(math.Exp2(float64(len(yuma.roles)))); i++ {
		for j := 0; j < int(math.Exp2(float64(len(yuma.roles) - 1)) + 1); j++ {
			if i & target != 0 || j == 0 {
				yuma.q.Set(i, j, 0.0)
			} else {
				rand.Seed(time.Now().UnixNano())
				yuma.q.Set(i, j, rand.Float64() * -1.0)
			}
		}
	}
}

func (yuma *Yuma) initializeState() (int, []string) {
	return yuma.startState(), make([]string, 0)
}

func (yuma *Yuma) derivePolicy_EpsilonGreedy(epsilon float64) {
	for i := 0; i < int(math.Exp2(float64(len(yuma.roles)))); i++ {
		maxAction := yuma.argMax_Action(i)
		choices := make([]randutil.Choice, 0)
		for _, action := range yuma.actions(i) {
			c := randutil.Choice{}
			c.Item = action
			if action == maxAction {
				c.Weight = int((1.0 - epsilon) * 10.0)
			} else {
				c.Weight = int(epsilon * 10.0)
			}
			choices = append(choices, c)
		}
		yuma.policy_EpsilonGreedy[i] = choices
	}
}

func (yuma *Yuma) takeAction(sigma int, omega int, path []string, state int, action int) (error, bool, float64, int) {
	reward := math.MaxFloat64 * -1.0
	successor := -1
	success := false

	for i := 0; i < sigma; i++ {
		instance := ec2.NewEC2Instance()
		err := instance.Create()
		if err != nil {
			return err, success, math.MaxFloat64 * -1.0, -1
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
			reward = -1.0
			successor = state | action
			success = true
			break
		} else {
			instance.Delete()
			yuma.ansible.DeleteInventory()
			reward = -10.0
			successor = state
		}
	}

	return nil, success, reward, successor
}

func (yuma *Yuma) argMax_Action(state int) int {
	maxQ := math.MaxFloat64 * -1.0
	maxAction := 0
	for _, action := range yuma.actions(state) {
		if yuma.q.At(state, action) > maxQ {
			maxQ = yuma.q.At(state, action)
			maxAction = action
		}
	}

	return maxAction
}

func (yuma *Yuma) max_Q(state int) float64 {
	if yuma.isTerminal(state) {
		return 0.0
	}

	maxQ := math.MaxFloat64 * -1.0
	for _, action := range yuma.actions(state) {
		if yuma.q.At(state, action) > maxQ {
			maxQ = yuma.q.At(state, action)
		}
	}

	return maxQ
}

func (yuma *Yuma) log(exportResults string, filePath string) error {
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