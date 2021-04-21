package search

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"os/exec"
	"bufio"
	"gonum.org/v1/gonum/mat"
)

type Yuma struct {
	inventory, playbook, repository string
	roles map[string]int
	configurations map[int]string
	searchTree *mat.Dense
}

func NewYuma(inventory, playbook, repository string) *Yuma {
	yuma := new(Yuma)
	yuma.inventory = inventory
	yuma.playbook = playbook
	yuma.repository = repository
	yuma.roles = make(map[string]int, 0)
	yuma.configurations = make(map[int]string, 0)
	if err := yuma.identifyRoles(repository); err != nil {
		return nil
	}
	yuma.searchTree = mat.NewDense(int(math.Exp2(float64(len(yuma.roles)))), int(math.Exp2(float64(len(yuma.roles) - 1)) + 1), nil)

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

func (yuma *Yuma) isEndSucc(state int) bool {
	if state == int(math.Exp2(float64(len(yuma.roles)))) - 1 {
		return true
	} else {
		return false
	}
}

func (yuma *Yuma) isEndFail(state, depth int) bool {
	if state < int(math.Exp2(float64(depth))) - 1 {
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

func (yuma *Yuma) playRole(roles []string, lifecycle string) bool {
	if roles[0] == "" {
		return true
	}

	err := createPlaybook(yuma.playbook, yuma.repository, roles, lifecycle)
	if err != nil {
		return false
	}

	cmd := exec.Command("ansible-playbook", "-i", yuma.inventory, yuma.playbook)
  
	out, err := cmd.StdoutPipe()
	if err != nil {
	  return false
	}
  
	scanner := bufio.NewScanner(out)
	go func() {
	  for scanner.Scan() {
		fmt.Println(scanner.Text())
	  }
	}()
  
	err = cmd.Start()
	if err != nil {
	  return false
	}
  
	err = cmd.Wait()
	if err != nil {
	  return false
	}

	return true
}

func (yuma *Yuma) BuildSearchTree(state int, depth int, path []string) *mat.Dense {
	if yuma.isEndSucc(state) {
		return yuma.searchTree
	}
	
	if yuma.isEndFail(state, depth) {
		return yuma.searchTree
	}

	for _, action := range yuma.actions(state) {
		yuma.playRole(path, "install")
		deletePlaybook(yuma.playbook)

		role := make([]string, 1)
		role[0] = yuma.configurations[action]

		if yuma.playRole(role, "install") {
			deletePlaybook(yuma.playbook)
			yuma.searchTree.Set(state, action, float64(state | action))
		} else {
			deletePlaybook(yuma.playbook)
			yuma.searchTree.Set(state, action, float64(state))
		}

		successor := yuma.successor(state, action)
		if successor != state {
			path[depth] = yuma.configurations[action]
		}

		yuma.playRole(path, "remove")
		deletePlaybook(yuma.playbook)

		yuma.BuildSearchTree(successor, depth + 1, path)
		clearPath(path, depth)
	}

	return yuma.searchTree
}

func (yuma *Yuma) DetermineExecutionOrder(state int, depth int, path []string, target int, memDepth map[int]int, memPath map[int][]string) (int, []string) {
	fmt.Println(state)
	fmt.Println(depth)
	fmt.Println()
	
	if _, ok := memDepth[state]; ok {
		/* fmt.Println(state)
		fmt.Println(memDepth[state])
		fmt.Println(memPath[state])
		fmt.Println() */
		return memDepth[state], memPath[state]
	}
	
	if state & target != 0 {
		memDepth[state] = depth
		memPath[state] = make([]string, len(yuma.roles))
		copy(memPath[state], path)
		/* fmt.Println(state)
		fmt.Println(memDepth[state])
		fmt.Println(memPath[state])
		fmt.Println() */
		return memDepth[state], memPath[state]
	}

	if yuma.isEndFail(state, depth) {
		memDepth[state] = len(yuma.roles) + 1
		memPath[state] = make([]string, len(yuma.roles))
		copy(memPath[state], path)
		return memDepth[state], memPath[state]
	}

	minDepth := len(yuma.roles)
	minPath := make([]string, len(yuma.roles))
	for _, action := range yuma.actions(state) {
		successor := yuma.successor(state, action)
		if successor != state {
			clearPath(path, depth)
			path[depth] = yuma.configurations[action]
		}

		tmpDepth, tmpPath := yuma.DetermineExecutionOrder(successor, depth + 1, path, target, memDepth, memPath)
		if tmpDepth <= minDepth {
			minDepth = tmpDepth
			copy(minPath, tmpPath)
		}

		/* fmt.Println(state)
		fmt.Println(memDepth[state])
		fmt.Println(memPath[state])
		fmt.Println() */
		memDepth[state] = minDepth
		memPath[state] = minPath
	}

	/* fmt.Println(memDepth)
	fmt.Println(memPath)
	fmt.Println()
	fmt.Println() */

	return memDepth[state], memPath[state]
}

func (yuma *Yuma) PrintDeploymentPlan(deploymentplans [][]string) {
	for _, deploymentplan := range deploymentplans {
		for _, role := range deploymentplan {
			fmt.Println(role)
		}
		fmt.Println("------------------")
	}
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

func createPlaybook(playbook string, repository string, roles []string, lifecycle string) error {
	file, err := os.OpenFile(playbook, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	defer file.Close()

	if _, err := file.WriteString(stringBuilder(repository, roles, lifecycle)); err != nil {
		return err
	}

	return nil
}

func deletePlaybook(playbook string) error {
	if err := os.Remove(playbook); err != nil {
		return err
	}

	return nil
}

func stringBuilder(repository string, roles []string, lifecycle string) string {
	export := "---\n"
	export += "- hosts: Yuma\n"
	export += "  vars:\n"
	export += "  - lifecycle: \"" + lifecycle + "\"\n"
	export += "  roles:\n"
	switch lifecycle {
	case "install":
		for i := 0; i < len(roles); i++ {
			export += "    - " + repository + roles[i] + "\n"
		}
	case "remove":
		for i := len(roles) - 1; i >= 0; i-- {
			export += "    - " + repository + roles[i] + "\n"
		}
	}

	return export
}

func clearPath(path []string, depth int) {
	for i := depth; i < len(path); i++ {
		path[i] = ""
	}
}