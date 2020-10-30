package ansible

import (
	"os"
)

type Playbook struct {
	Software map[string][]string
}

func NewPlaybook() *Playbook {
	playbook := new(Playbook)
	playbook.Software = make(map[string][]string, 0)

	return playbook
}

func (playbook *Playbook) AddSoftwareStack(role string, stack []string) {
	if _, ok := playbook.Software[role]; ok {
		playbook.Software[role] = append(playbook.Software[role], stack...)
	} else {
		playbook.Software[role] = stack
	}
}

func (playbook *Playbook) Export() error {
	file, err := os.OpenFile("build.yml", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	defer file.Close()

	if _, err := file.WriteString(playbook.stringBuilder()); err != nil {
		return nil
	}

	return nil
}

func (playbook *Playbook) stringBuilder() string {
	export := "---\n"

	local := playbook.Software["localhost"]
	export += "- hosts: localhost\n"
	export += "  roles:\n"

	for _, software := range local {
		export += "    - ./factors_of_production/software/" + software + "\n"
	}

	export += "\n"

	jumphost := playbook.Software["Jumphost"]
	export += "- hosts: Jumphost\n"
	export += "  roles:\n"

	for _, software := range jumphost {
		export += "    - ./factors_of_production/software/" + software + "\n"
	}

	export += "\n"

	all := playbook.Software["all"]
	export += "- hosts: all\n"
	export += "  roles:\n"

	for _, software := range all {
		export += "    - ./factors_of_production/software/" + software + "\n"
	}

	export += "\n"

	etcd := playbook.Software["Etcd"]
	export += "- hosts: Etcd\n"
	export += "  roles:\n"

	for _, software := range etcd {
		export += "    - ./factors_of_production/software/" + software + "\n"
	}

	export += "\n"

	loadbalancer := playbook.Software["LoadBalancer"]
	export += "- hosts: LoadBalancer\n"
	export += "  roles:\n"

	for _, software := range loadbalancer {
		export += "    - ./factors_of_production/software/" + software + "\n"
	}

	export += "\n"

	controlplane := playbook.Software["Controlplane"]
	export += "- hosts: Controlplane\n"
	export += "  roles:\n"

	for _, software := range controlplane {
		export += "    - ./factors_of_production/software/" + software + "\n"
	}

	export += "\n"

	worker := playbook.Software["Worker"]
	export += "- hosts: Worker\n"
	export += "  roles:\n"

	for _, software := range worker {
		export += "    - ./factors_of_production/software/" + software + "\n"
	}

	export += "\n"

	for role, stack := range playbook.Software {
		if role != "localhost" && role != "Jumphost" && role != "LoadBalancer" && role != "all" && role != "Etcd" && role != "Controlplane" && role != "Worker" {
			export += "- hosts: " + role + "\n"
			export += "  roles:\n"
			for _, software := range stack {
				export += "    - ./factors_of_production/software/" + software + "\n"
			}
		}
		
		export += "\n"
	}

	return export
}