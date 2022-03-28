package infrastructureascode

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
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

func (ansible *Ansible) CreateExecutionOrder(path string, pathPrefix string, roles []string, lifecycle string) error {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	defer file.Close()

	if _, err := file.WriteString(ansible.stringBuilder(pathPrefix, roles, lifecycle)); err != nil {
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

func (ansible *Ansible) Execute(target int, pathPrefix string, roles []string, lifecycle string) bool {
	file, err := os.OpenFile("logs_" + strconv.Itoa(target) + ".txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return false
	}
	defer file.Close()

	if err := ansible.CreateExecutionOrder("yuma_" + strconv.Itoa(target) + ".yml", pathPrefix, roles, lifecycle); err != nil {
		return false
	}
	
	if _, err := os.Stat("hosts_" + strconv.Itoa(target)); errors.Is(err, os.ErrNotExist) {
		fmt.Println("ANSIBLE ERROR: INVENTORY IS MISSING")
		return false
	}
	cmd := exec.Command("ansible-playbook", "-i", "hosts_" + strconv.Itoa(target), "yuma_" + strconv.Itoa(target) + ".yml", "--limit", "yuma1")
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

func (ansible *Ansible) stringBuilder(pathPrefix string, roles []string, lifecycle string) string {
	export := "---\n"
	export += "- hosts: yuma1\n"
	export += "  vars:\n"
	export += "  - lifecycle: \"" + lifecycle + "\"\n"
	export += "  roles:\n"
	switch lifecycle {
	case "install":
		for i := 0; i < len(roles); i++ {
			export += "    - " + pathPrefix + ansible.GetRepository() + roles[i] + "\n"
		}
	case "remove":
		for i := len(roles) - 1; i >= 0; i-- {
			export += "    - " + pathPrefix + ansible.GetRepository() + roles[i] + "\n"
		}
	}

	return export
}