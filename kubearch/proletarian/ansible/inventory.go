package ansible

import (
	"os"
	"strings"
)

type Inventory struct {
	Role, Output string
	Variables map[string]map[string]interface{}
}

func NewInventory(role string) *Inventory {
	inventory := new(Inventory)

	inventory.Output = "[" + role + "]" + "\n"
	inventory.Role = role
	inventory.Variables = make(map[string]map[string]interface{}, 0)

	return inventory
}

func (inventory *Inventory) AddVirtualMachine(id string) {
	inventory.Variables[id] = make(map[string]interface{}, 0)
}

func (inventory *Inventory) AddVariable(instanceID, variable string, value interface{}) {
	inventory.Variables[instanceID][variable] = value
}

func (inventory *Inventory) Export() error {
	file, err := os.OpenFile("hosts", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	defer file.Close()

	inventory.stringBuilder()

	if _, err := file.WriteString(inventory.Output); err != nil {
		return nil
	}

	return nil
}

func (inventory *Inventory) stringBuilder() {
	for instanceID, variables := range inventory.Variables {
		if instanceID == "" {
			inventory.Output += instanceID
		} else {
			inventory.Output += instanceID + " "
		}

		for variable, value := range variables {
			switch variable {
			case "public_ip":
				variable = "ansible_host"
			case "user":
				variable = "ansible_user"
			}
			inventory.Output += variable + "=" + strings.TrimSuffix(value.(string), "\n") + " "

			if inventory.Role == "all:vars" {
				inventory.Output += "\n"
			}
		}
		inventory.Output += "\n"
	}

	inventory.Output += "\n"
}