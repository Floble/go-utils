package k8sml

import (
	"reflect"
	"strings"
	terraform "KubeArch/kubearch/proletarian/terraform"
)

type IAMPolicy struct {
	ID string
	Type string
	RuntimeVariables map[string]string
	Role []*IAMRole
	CloudProvider CloudProvider
}

func NewPolicy(id, kind string) *IAMPolicy {
	policy := new(IAMPolicy)
	policy.ID = id
	policy.Type = kind
	policy.RuntimeVariables = make(map[string]string, 0)
	policy.Role = make([]*IAMRole, 0)

	return policy
}

func (policy *IAMPolicy) GetID() string {
    return policy.ID
}

func (policy *IAMPolicy) GetVariableValue(variable string) interface{} {
    e := reflect.ValueOf(policy).Elem()
	var value interface{}

	for i := 0; i < e.NumField(); i++ {
		if strings.ToLower(e.Type().Field(i).Name) == strings.ToLower(variable) {
			value = e.Field(i).Interface()
			return value
		}
	}

	return nil
}

func (policy *IAMPolicy) GetRuntimeVariables() map[string]string {
    return policy.RuntimeVariables
}

func (policy *IAMPolicy) AddRuntimeVariable(key, value string) {
    policy.RuntimeVariables[key] = value
}

func (policy *IAMPolicy) ExportModule() error {
	e := reflect.ValueOf(policy).Elem()

    module := terraform.NewModule(policy.ID, strings.Split(policy.CloudProvider.GetType(), "*k8sml.")[1], strings.Split(reflect.TypeOf(policy).String(), "*k8sml.")[1])
    for i := 0; i < e.NumField(); i++ {
		key := e.Type().Field(i).Name
        value := e.Field(i).Interface()
        module.AddVariable(key, value)
	}

	var output *terraform.Output

	if policy.Type == "master" {
		output = terraform.NewOutput(policy.ID + "_" + policy.Type + "[0]")
	} else {
		output = terraform.NewOutput(policy.ID + "_" + policy.Type + "[0]")
	}

    for key, _ := range policy.RuntimeVariables {
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