package infrastructureascode

import (
	//"bufio"
	//"fmt"
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

func (ansible *Ansible) PlayRoles(roles []string, lifecycle string) bool {
	if roles[0] == "" {
		return true
	}

	err := ansible.createPlaybook(roles, lifecycle)
	if err != nil {
		return false
	}

	cmd := exec.Command("ansible-playbook", "-i", ansible.GetInventory(), ansible.GetPlaybook())
  
	_, err = cmd.StdoutPipe()
	if err != nil {
		ansible.deletePlaybook()
		return false
	}
  
	/* scanner := bufio.NewScanner(out)
	go func() {
		for scanner.Scan() {
		  fmt.Println(scanner.Text())
		}
	}() */

	err = cmd.Start()
	if err != nil {
		ansible.deletePlaybook()
		return false
	}
  
	err = cmd.Wait()
	if err != nil {
		ansible.deletePlaybook()
		return false
	}

	ansible.deletePlaybook()

	return true
}

func (ansible *Ansible) CreateInventory(publicIP string) error {
	file, err := os.OpenFile(ansible.GetInventory(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	defer file.Close()

	export := "[Yuma]\n"
	export += "yuma1 ansible_host=" + publicIP + " ansible_user=ubuntu\n\n"
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

func (ansible *Ansible) DeleteInventory() error {
	if err := os.Remove(ansible.GetInventory()); err != nil {
		return err
	}

	return nil
}

func (ansible *Ansible) createPlaybook(roles []string, lifecycle string) error {
	file, err := os.OpenFile(ansible.GetPlaybook(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	defer file.Close()

	if _, err := file.WriteString(ansible.stringBuilder(roles, lifecycle)); err != nil {
		return err
	}

	return nil
}

func (ansible *Ansible) deletePlaybook() error {
	if err := os.Remove(ansible.GetPlaybook()); err != nil {
		return err
	}

	return nil
}

func (ansible *Ansible) stringBuilder(roles []string, lifecycle string) string {
	export := "---\n"
	export += "- hosts: Yuma\n"
	export += "  vars:\n"
	export += "  - lifecycle: \"" + lifecycle + "\"\n"
	export += "  roles:\n"
	switch lifecycle {
	case "install":
		for i := 0; i < len(roles); i++ {
			export += "    - " + ansible.GetRepository() + roles[i] + "\n"
		}
	case "remove":
		for i := len(roles) - 1; i >= 0; i-- {
			export += "    - " + ansible.GetRepository() + roles[i] + "\n"
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

func (ansible *Ansible) GetRepository() string {
	return ansible.repository
}