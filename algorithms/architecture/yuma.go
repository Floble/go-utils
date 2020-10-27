package architecture

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"os/exec"
	"bufio"
)

type Yuma struct {
	Roles map[string]uint8
	States []uint8
}

func NewYuma(path string) *Yuma {
	yuma := new(Yuma)
	yuma.Roles = make(map[string]uint8, 0)
	if err := yuma.identifyRoles(path); err != nil {
		return nil
	}
	yuma.States = make([]uint8, len(yuma.Roles) + 1)

	return yuma
}

func (yuma *Yuma) identifyRoles(path string) error {
	roles, err := ioutil.ReadDir(path)
    if err != nil {
        return err
    }

	for i := 0; i < len(roles); i++ {
		var config uint8 = 0
		config |= uint8(math.Exp2(float64(i)))
		yuma.Roles[roles[i].Name()] = config
	}
	
	return nil
}

func (yuma *Yuma) Dfs(mask uint8, depth int) bool {
	if mask == uint8(math.Exp2(float64(len(yuma.Roles)))) - 1 {
		return true
	}

	for role, config := range yuma.Roles {
		if mask & config != 0 {
			continue
		}

		if yuma.playRole(role) {
			deletePlaybook()
			mask |= config
			yuma.States[depth] = mask

			fmt.Println(role)

			/* fmt.Printf("Mask: ")
			fmt.Printf("%b", mask)
			fmt.Println()
			fmt.Printf("Detph: ")
			fmt.Printf("%d", depth)
			fmt.Println()
			fmt.Println() */

			if yuma.Dfs(mask, depth + 1) {
				return true
			}
		} else {
			deletePlaybook()
		}
	}

	return false
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

	cmd := exec.Command("ansible-playbook", "-i", "/home/floble/yuma-test/hosts", "/home/floble/yuma-test/build.yml", "-vvv")
  
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