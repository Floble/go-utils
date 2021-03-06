package architecture

import (
	"fmt"
	"strconv"
	"io/ioutil"
	"math"
	"os"
	"os/exec"
	"bufio"
)

type Yuma struct {
	Inventory, Playbook, Examples string
	Roles map[string]int
	Configurations map[int]string
	DeploymentPlan []int
}

func NewYuma(inventory, playbook, examples string) *Yuma {
	yuma := new(Yuma)
	yuma.Inventory = inventory
	yuma.Playbook = playbook
	yuma.Examples = examples
	yuma.Roles = make(map[string]int, 0)
	yuma.Configurations = make(map[int]string, 0)
	if err := yuma.identifyRoles(examples); err != nil {
		return nil
	}
	yuma.DeploymentPlan = make([]int, len(yuma.Roles) + 1)

	return yuma
}

func (yuma *Yuma) identifyRoles(path string) error {
	roles, err := ioutil.ReadDir(path)
    if err != nil {
        return err
    }

	for i := 0; i < len(roles); i++ {
		var config int = 0
		config |= int(math.Exp2(float64(i)))
		yuma.Roles[roles[i].Name()] = config
		yuma.Configurations[config] = roles[i].Name()
	}
	
	return nil
}

func (yuma *Yuma) DetermineDeploymentPlan_Greedy_Iterative() bool {
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
}

func (yuma *Yuma) DetermineDeploymentPlan_Backtracking(state int, depth int) bool {
	if state == int(math.Exp2(float64(len(yuma.Roles)))) - 1 {
		return true
	}

	if state < int(math.Exp2(float64(depth))) - 1 {
		return false
	}

	result := false

	if state >= int(math.Exp2(float64(depth))) - 1 {
		for role, config := range yuma.Roles {
			if yuma.playRole(role) {
				deletePlaybook(yuma.Playbook)
				result = result || yuma.DetermineDeploymentPlan_Backtracking(state | config, depth + 1)
			} else {
				deletePlaybook(yuma.Playbook)
				result = result || yuma.DetermineDeploymentPlan_Backtracking(state, depth + 1)
			}
		}
	}

	return result
}

func (yuma *Yuma) PrintDeploymentPlan() {
	for i := 1; i < len(yuma.DeploymentPlan); i++ {
		config := yuma.DeploymentPlan[i - 1] ^ yuma.DeploymentPlan[i]
		fmt.Println(strconv.Itoa(i) + " - " + yuma.Configurations[config])
	}
}

func createPlaybook(playbook, examples, role string) error {
	file, err := os.OpenFile(playbook, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	defer file.Close()

	if _, err := file.WriteString(stringBuilder(examples, role)); err != nil {
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

func stringBuilder(examples, role string) string {
	export := "---\n"
	export += "- hosts: Yuma\n"
	export += "  roles:\n"
	export += "    - " + examples + role + "\n"

	return export
}

func (yuma *Yuma) playRole(role string) bool {
	err := createPlaybook(yuma.Playbook, yuma.Examples, role)
	if err != nil {
		return false
	}

	cmd := exec.Command("ansible-playbook", "-i", yuma.Inventory, yuma.Playbook)
  
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