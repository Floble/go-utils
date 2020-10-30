package k8sml

import (
	"reflect"
	"strings"
	terraform "KubeArch/kubearch/proletarian/terraform"
)

type IAMRole struct {
	ID string
	RuntimeVariables map[string]string
	Policy *IAMPolicy
	VirtualMachine []VirtualMachine
}

func NewIAMRole(id string) *IAMRole {
	role := new(IAMRole)
	role.ID = id
	role.RuntimeVariables = make(map[string]string, 0)

	return role
}

func (role *IAMRole) GetID() string {
    return role.ID
}

func (role *IAMRole) GetVariableValue(variable string) interface{} {
    e := reflect.ValueOf(role).Elem()
	var value interface{}

	for i := 0; i < e.NumField(); i++ {
		if strings.ToLower(e.Type().Field(i).Name) == strings.ToLower(variable) {
			value = e.Field(i).Interface()
			return value
		}
	}

	return nil
}

func (role *IAMRole) GetRuntimeVariables() map[string]string {
    return role.RuntimeVariables
}

func (role *IAMRole) AddRuntimeVariable(key, value string) {
    role.RuntimeVariables[key] = value
}

func (role *IAMRole) ExportModule() error {
	e := reflect.ValueOf(role).Elem()

	roleType := strings.Split(reflect.TypeOf(role).String(), "*k8sml.")[1]
	policyType := strings.Split(reflect.TypeOf(role.Policy).String(), "*k8sml.")[1]

    module := terraform.NewModule(role.ID, strings.Split(role.Policy.CloudProvider.GetType(), "*k8sml.")[1], strings.Split(reflect.TypeOf(role).String(), "*k8sml.")[1])
    for i := 0; i < e.NumField(); i++ {
		key := e.Type().Field(i).Name
        value := e.Field(i).Interface()
        module.AddVariable(key, value)
	}
	module.AddVariable(policyType, "${module." + role.Policy.ID + "." + strings.ToLower(role.Policy.Type) + "}")

	output := terraform.NewOutput(role.ID + "_" + strings.ToLower(roleType))

    for key, _ := range role.RuntimeVariables {
        output.AddVariable(key)
    }

	if err := module.Export(); err != nil {
        return err
    }

    if err := output.Export(); err != nil {
        return err
    }

	return nil
}