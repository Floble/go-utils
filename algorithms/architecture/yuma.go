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
	Roles map[string]int
	Configurations map[int]string
	DeploymentPlan []int
}

func NewYuma(path string) *Yuma {
	yuma := new(Yuma)
	yuma.Roles = make(map[string]int, 0)
	yuma.Configurations = make(map[int]string, 0)
	if err := yuma.identifyRoles(path); err != nil {
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

func (yuma *Yuma) DetermineDeploymentPlan() bool {
	deploymentPlan := make([]int, len(yuma.Roles) + 1)
	deploymentPlan[0] = 0

	for i := 1; i <= len(yuma.Roles); i++ {
		for role, config := range yuma.Roles {
			if deploymentPlan[i - 1] & config != 0 {
				continue
			}

			if yuma.playRole(role) {
				deletePlaybook()
				deploymentPlan[i] = deploymentPlan[i - 1] | config
				break
			} else {
				deletePlaybook()
			}
		}
	}

	yuma.DeploymentPlan = deploymentPlan

	if deploymentPlan[len(yuma.Roles)] == int(math.Exp2(float64(len(yuma.Roles)))) - 1 {
		return true
	} else {
		return false
	}
}

func (yuma *Yuma) PrintDeploymentPlan() {
	for i := 1; i < len(yuma.DeploymentPlan); i++ {
		config := yuma.DeploymentPlan[i - 1] ^ yuma.DeploymentPlan[i]
		fmt.Println(strconv.Itoa(i) + " - " + yuma.Configurations[config])
	}
}

func createPlaybook(role string) error {
	file, err := os.OpenFile("/home/floble/yuma-test/build.yml", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	defer file.Close()

	if _, err := file.WriteString(stringBuilder(role)); err != nil {
		return err
	}

	return nil
}

func deletePlaybook() error {
	if err := os.Remove("/home/floble/yuma-test/build.yml"); err != nil {
		return err
	}

	return nil
}

func stringBuilder(role string) string {
	export := "---\n"
	export += "- hosts: Yuma\n"
	export += "  roles:\n"
	export += "    - /home/floble/yuma-test/roles/" + role + "\n"

	return export
}

func (yuma *Yuma) playRole(role string) bool {
	err := createPlaybook(role)
	if err != nil {
		return false
	}

	cmd := exec.Command("ansible-playbook", "-i", "/home/floble/yuma-test/hosts", "/home/floble/yuma-test/build.yml")
  
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