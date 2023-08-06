package infrastructureascode

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type Ansible struct {
	repository string
}

func NewAnsible(repository string) *Ansible {
	ansible := new(Ansible)
	ansible.repository = repository

	return ansible
}

func (ansible *Ansible) GetRepository() string {
	return ansible.repository
}

/* func (ansible *Ansible) CreateExecutionOrder(target int, path string, pathPrefix string, roles []string, lifecycle string, host string) error {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	defer file.Close()

	if host != "localhost" {
		if _, err := file.WriteString(ansible.stringBuilder(target, pathPrefix, host, roles, lifecycle)); err != nil {
			return err
		}
	} else {
		if _, err := file.WriteString(ansible.stringBuilder(target, pathPrefix, host, roles, lifecycle)); err != nil {
			return err
		}
	}

	return nil
} */

func (ansible *Ansible) CreateExecutionOrder(target int, path string, pathPrefix string, roles []string, lifecycle string, host string) error {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	defer file.Close()

	if _, err := file.WriteString(ansible.stringBuilder(target, pathPrefix, host, roles, lifecycle)); err != nil {
		return err
	}

	return nil
}

func (ansible *Ansible) DeleteExecutionOrder(path string) error {
	if err := os.Remove(path); err != nil {
		return err
	}

	return nil
}

func (ansible *Ansible) CreateEnvironmentDescription(target int, description interface{}) error {
	file, err := os.OpenFile("hosts_" + strconv.Itoa(target), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	defer file.Close()

	export := "[Yuma]\n"
	export += "yuma1 ansible_host=" + fmt.Sprintf("%v", description) + " ansible_user=ubuntu\n\n"
	export += "[all:vars]\n"
	export += "ansible_python_interpreter=/usr/bin/python3\n"
	export += "ansible_ssh_common_args='-o StrictHostKeyChecking=no'\n"
	export += "ansible_ssh_extra_args='-o StrictHostKeyChecking=no'\n"
	export += "ansible_ssh_private_key_file=~/.ssh/id_rsa_floble"

	if _, err := file.WriteString(export); err != nil {
		return err
	}

	return nil
}

func (ansible *Ansible) RemoveEnvironmentDescription(target int) error {
	if err := os.Remove("hosts_" + strconv.Itoa(target)); err != nil {
		return err
	}

	return nil
}

/* func (ansible *Ansible) Execute(target int, pathPrefix string, roles []string, lifecycle string, host string) bool {
	file, err := os.OpenFile("logs_" + strconv.Itoa(target) + ".txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return false
	}
	defer file.Close()

	if host != "localhost" {
		if err := ansible.CreateExecutionOrder(0, "yuma_" + strconv.Itoa(target) + ".yml", pathPrefix, roles, lifecycle, host); err != nil {
			fmt.Println(err)
			return false
		}
	} else {
		if err := ansible.CreateExecutionOrder(target, "yuma_" + strconv.Itoa(target) + ".yml", pathPrefix, roles, lifecycle, host); err != nil {
			fmt.Println(err)
			return false
		}
	}
	
	var cmd *exec.Cmd
	if host != "localhost" {
		if _, err := os.Stat("hosts_" + strconv.Itoa(target)); errors.Is(err, os.ErrNotExist) {
			fmt.Println("ANSIBLE ERROR: INVENTORY IS MISSING")
			return false
		}
		cmd = exec.Command("ansible-playbook", "-i", "hosts_" + strconv.Itoa(target), "yuma_" + strconv.Itoa(target) + ".yml", "--limit", "yuma1")
	} else {
		cmd = exec.Command("ansible-playbook", "yuma_" + strconv.Itoa(target) + ".yml")
	}
	
	out, err := cmd.StdoutPipe()
	if err != nil {
		ansible.DeleteExecutionOrder("yuma_" + strconv.Itoa(target) + ".yml")
		return false
	}
  
	scanner := bufio.NewScanner(out)
	go func() {
		for scanner.Scan() {
			file.WriteString(scanner.Text())
		}
	}()
	file.WriteString("\n\n")

	if err := cmd.Start(); err != nil {
		ansible.DeleteExecutionOrder("yuma_" + strconv.Itoa(target) + ".yml")
		return false
	}
	if err := cmd.Wait(); err != nil {
		ansible.DeleteExecutionOrder("yuma_" + strconv.Itoa(target) + ".yml")
		return false
	}

	ansible.DeleteExecutionOrder("yuma_" + strconv.Itoa(target) + ".yml")

	return true
} */

func (ansible *Ansible) Execute(target int, pathPrefix string, roles []string, lifecycle string, host string) bool {
	file, err := os.OpenFile("logs_" + strconv.Itoa(target) + ".txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return false
	}
	defer file.Close()

	if err := ansible.CreateExecutionOrder(0, "yuma_" + strconv.Itoa(target) + ".yml", pathPrefix, roles, lifecycle, host); err != nil {
		fmt.Println(err)
		return false
	}
	
	var cmd *exec.Cmd
	if host != "localhost" {
		if _, err := os.Stat("hosts_" + strconv.Itoa(target)); errors.Is(err, os.ErrNotExist) {
			fmt.Println("ANSIBLE ERROR: INVENTORY IS MISSING")
			return false
		}
		cmd = exec.Command("ansible-playbook", "-i", "hosts_" + strconv.Itoa(target), "yuma_" + strconv.Itoa(target) + ".yml", "--limit", "yuma1")
	} else {
		cmd = exec.Command("ansible-playbook", "yuma_" + strconv.Itoa(target) + ".yml")
	}
	
	out, err := cmd.StdoutPipe()
	if err != nil {
		ansible.DeleteExecutionOrder("yuma_" + strconv.Itoa(target) + ".yml")
		return false
	}
  
	scanner := bufio.NewScanner(out)
	go func() {
		for scanner.Scan() {
			file.WriteString(scanner.Text())
		}
	}()
	file.WriteString("\n\n")

	if err := cmd.Start(); err != nil {
		ansible.DeleteExecutionOrder("yuma_" + strconv.Itoa(target) + ".yml")
		return false
	}
	if err := cmd.Wait(); err != nil {
		ansible.DeleteExecutionOrder("yuma_" + strconv.Itoa(target) + ".yml")
		return false
	}

	ansible.DeleteExecutionOrder("yuma_" + strconv.Itoa(target) + ".yml")

	return true
}

func (ansible *Ansible) DetermineInputs(role string) (error, []string) {
	inputs := make([]string, 0)
	
	filePath := ansible.GetRepository() + role + "/defaults/main.yml"
    readFile, err := os.Open(filePath)
    if err != nil {
		return err, nil
    }

    fileScanner := bufio.NewScanner(readFile)
    fileScanner.Split(bufio.ScanLines)
    var fileLines []string
  
    for fileScanner.Scan() {
        fileLines = append(fileLines, fileScanner.Text())
    }

    readFile.Close()
  
    for _, line := range fileLines {
        if strings.Contains(line, ":") {
			inputs = append(inputs, strings.Split(line, ":")[0])
		}
    }

    return nil, inputs
}

func (ansible *Ansible) DetermineOutputs(role string) (error, []string) {
	outputs := make([]string, 0)
	
	filePath := ansible.GetRepository() + role + "/tasks/main.yml"
    readFile, err := os.Open(filePath)
    if err != nil {
		return err, nil
    }

    fileScanner := bufio.NewScanner(readFile)
    fileScanner.Split(bufio.ScanLines)
    var fileLines []string
  
    for fileScanner.Scan() {
        fileLines = append(fileLines, fileScanner.Text())
    }

    readFile.Close()
  
    for i, line := range fileLines {
        if strings.Contains(line, "set_fact:") {
			outputs = append(outputs, strings.Split(fileLines[i + 1], ":")[0])
		}
    }

    return nil, outputs
}

/* func (ansible *Ansible) stringBuilder(target int, pathPrefix string, host string, roles []string, lifecycle string) string {
	export := "---\n"
	switch lifecycle {
	case "create":
		for i := 0; i < len(roles); i++ {
			export += "- hosts: " + host + "\n"
			export += "  vars:\n"
			export += "    - lifecycle: \"" + lifecycle + "\"\n"
			err, variables := ansible.DetermineInputs(roles[i])
			if err == nil && host != "all" {
				for j := 0; j < len(variables); j++ {
					if target <= 0 {
						export += "  - " + variables[j] + ": \"{{ hostvars['localhost']['" + variables[j] + "'] }}\"\n"
					} else {
						export += "  - " + variables[j] + strconv.Itoa(target) + ": \"{{ hostvars['localhost']['" + variables[j] + strconv.Itoa(target) + "'] }}\"\n"
					}
				}
			}
			export += "  roles:\n"
			if target <= 0 {
				export += "    - " + pathPrefix + ansible.GetRepository() + roles[i] + "\n\n"
			} else {
				export += "    - " + pathPrefix + strings.Split(ansible.GetRepository(), "/")[0] + "_" + strconv.Itoa(target) + "/" +  roles[i] + "\n\n"
			}
		}
	case "remove":
		for i := len(roles) - 1; i >= 0; i-- {
			export += "- hosts: " + host + "\n"
			export += "  vars:\n"
			export += "  - lifecycle: \"" + lifecycle + "\"\n"
			err, variables := ansible.DetermineInputs(roles[i])
			if err == nil && host != "all" {
				for j := 0; j < len(variables); j++ {
					if target <= 0 {
						export += "  - " + variables[j] + ": \"\"\n"
					} else {
						export += "  - " + variables[j] + strconv.Itoa(target) + ": \"\"\n"
					}
				}
			}
			export += "  roles:\n"
			if target <= 0 {
				export += "    - " + pathPrefix + ansible.GetRepository() + roles[i] + "\n\n"
			} else {
				export += "    - " + pathPrefix + strings.Split(ansible.GetRepository(), "/")[0] + "_" + strconv.Itoa(target) + "/" +  roles[i] + "\n\n"
			}
		}
	}

	return export
} */

func (ansible *Ansible) stringBuilder(target int, pathPrefix string, host string, roles []string, lifecycle string) string {
	export := "---\n"
	switch lifecycle {
	case "create":
		for i := 0; i < len(roles); i++ {
			export += "- hosts: " + host + "\n"
			export += "  vars:\n"
			export += "    - lifecycle: \"" + lifecycle + "\"\n"
			export += "  roles:\n"
			if target <= 0 {
				export += "    - " + pathPrefix + ansible.GetRepository() + roles[i] + "\n\n"
			} else {
				export += "    - " + pathPrefix + strings.Split(ansible.GetRepository(), "/")[0] + "_" + strconv.Itoa(target) + "/" +  roles[i] + "\n\n"
			}
		}
	case "remove":
		for i := len(roles) - 1; i >= 0; i-- {
			export += "- hosts: " + host + "\n"
			export += "  vars:\n"
			export += "  - lifecycle: \"" + lifecycle + "\"\n"
			export += "  roles:\n"
			if target <= 0 {
				export += "    - " + pathPrefix + ansible.GetRepository() + roles[i] + "\n\n"
			} else {
				export += "    - " + pathPrefix + strings.Split(ansible.GetRepository(), "/")[0] + "_" + strconv.Itoa(target) + "/" +  roles[i] + "\n\n"
			}
		}
	}

	return export
}