package infrastructureascode

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
)

type Ansible struct {
	inventory, playbook, repository string
}

func NewAnsible(inventory, playbook, repository string) *Ansible {
	ansible := new(Ansible)
	ansible.inventory = inventory
	ansible.playbook = playbook
	ansible.repository = repository

	return ansible
}

func (ansible *Ansible) GetRepository() string {
	return ansible.repository
}

func (ansible *Ansible) SetPlaybook(playbook string) {
	ansible.playbook = playbook
}

func (ansible *Ansible) CreateExecutionOrder(pathPrefix string, roles []string, lifecycle string) error {
	file, err := os.OpenFile(ansible.GetPlaybook(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	defer file.Close()

	if _, err := file.WriteString(ansible.stringBuilder(pathPrefix, roles, lifecycle)); err != nil {
		return err
	}

	return nil
}

func (ansible *Ansible) DeleteExecutionOrder() error {
	if err := os.Remove(ansible.GetPlaybook()); err != nil {
		return err
	}

	return nil
}

func (ansible *Ansible) CreateEnvironmentDescription(description interface{}) error {
	file, err := os.OpenFile(ansible.GetInventory(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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

func (ansible *Ansible) RemoveEnvironmentDescription() error {
	if err := os.Remove(ansible.GetInventory()); err != nil {
		return err
	}

	return nil
}

func (ansible *Ansible) Execute(pathPrefix string, roles []string, lifecycle string) bool {
	file, _ := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer file.Close()

	if len(roles) == 0 {
		return true
	}

	err := ansible.CreateExecutionOrder(pathPrefix, roles, lifecycle)
	if err != nil {
		return false
	}

	cmd := exec.Command("ansible-playbook", "-i", ansible.GetInventory(), ansible.GetPlaybook())
  
	out, err := cmd.StdoutPipe()
	if err != nil {
		ansible.DeleteExecutionOrder()
		return false
	}
  
	scanner := bufio.NewScanner(out)
	go func() {
		for scanner.Scan() {
			file.WriteString(scanner.Text())
		}
	}()
	file.WriteString("\n\n")

	err = cmd.Start()
	if err != nil {
		ansible.DeleteExecutionOrder()
		return false
	}
  
	err = cmd.Wait()
	if err != nil {
		ansible.DeleteExecutionOrder()
		return false
	}

	ansible.DeleteExecutionOrder()

	return true
}

func (ansible *Ansible) stringBuilder(pathPrefix string, roles []string, lifecycle string) string {
	export := "---\n"
	export += "- hosts: all\n"
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

func (ansible *Ansible) GetInventory() string {
	return ansible.inventory
}

func (ansible *Ansible) GetPlaybook() string {
	return ansible.playbook
}