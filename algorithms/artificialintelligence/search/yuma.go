package search

import (
	"fmt"
	//"strconv"
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

/* func (yuma *Yuma) DetermineDeploymentPlan_Greedy_Iterative() bool {
	yuma.DeploymentPlan[0] = 0

	for d := 1; d <= len(yuma.Roles); d++ {
		if yuma.DeploymentPlan[d - 1] < int(math.Exp2(float64(d - 1))) - 1 {
			return false
		}

		for role, config := range yuma.Roles {
			if yuma.DeploymentPlan[d - 1] & config != 0 {
				continue
			}

			if yuma.playRole(role) {
				deletePlaybook(yuma.Playbook)
				yuma.DeploymentPlan[d] = yuma.DeploymentPlan[d - 1] | config
				break
			} else {
				deletePlaybook(yuma.Playbook)
				yuma.DeploymentPlan[d] = yuma.DeploymentPlan[d - 1]
			}
		}
	}

	return yuma.DeploymentPlan[len(yuma.Roles)] == int(math.Exp2(float64(len(yuma.Roles)))) - 1
}

func (yuma *Yuma) DetermineDeploymentPlan_Greedy_Recursive(state int, depth int) bool {
	if state == int(math.Exp2(float64(len(yuma.Roles)))) - 1 {
		yuma.DeploymentPlan[depth] = state
		return true
	}

	for role, config := range yuma.Roles {
		if state & config != 0 {
			continue
		}

		if yuma.playRole(role) {
			deletePlaybook(yuma.Playbook)
			if yuma.DetermineDeploymentPlan_Greedy_Recursive(state | config, depth + 1) {
				yuma.DeploymentPlan[depth] = state
				return true
			} else {
				return false
			}
		} else {
			deletePlaybook(yuma.Playbook)
		}
	}

	return false
}

func (yuma *Yuma) DetermineDeploymentPlan_Dfs(state int, depth int) bool {
	if state == int(math.Exp2(float64(len(yuma.Roles)))) - 1 {
		yuma.DeploymentPlan[depth] = state
		return true
	}

	if state < int(math.Exp2(float64(depth))) - 1 {
		return false
	}

	result := false

	if state >= int(math.Exp2(float64(depth))) - 1 {
		for role, config := range yuma.Roles {
			if state & config != 0 {
				continue
			}

			if yuma.playRole(role) {
				deletePlaybook(yuma.Playbook)
				if yuma.DetermineDeploymentPlan_Dfs(state | config, depth + 1) {
					yuma.DeploymentPlan[depth] = state
					result = result || true
					return result
				}
			} else {
				deletePlaybook(yuma.Playbook)
				result = result || yuma.DetermineDeploymentPlan_Dfs(state, depth + 1)
			}
		}
	}

	return result
} */

func (yuma *Yuma) DetermineDeploymentPlan_Backtracking(state int, depth int, path []string) *mat.Dense {
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

		yuma.DetermineDeploymentPlan_Backtracking(successor, depth + 1, path)
		clearDeploymentplan(path, depth)
	}

	return yuma.searchTree
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

func clearDeploymentplan(deploymentplan []string, depth int) {
	for i := depth; i < len(deploymentplan); i++ {
		deploymentplan[i] = ""
	}
}