package k8sml

import (
	"reflect"
	"strings"
	terraform "KubeArch/kubearch/proletarian/terraform"
)

type NATGateway struct {
	ID string `yaml:"id"`
	RuntimeVariables map[string]string
	Subnet *Subnet
}

func (natGateway *NATGateway) GetID() string {
	return natGateway.ID
}

func (natGateway *NATGateway) GetVariableValue(variable string) interface{} {
	e := reflect.ValueOf(natGateway).Elem()
	var value interface{}

	for i := 0; i < e.NumField(); i++ {
		if strings.ToLower(e.Type().Field(i).Name) == strings.ToLower(variable) {
			value = e.Field(i).Interface()
			return value
		}
	}

	return nil
}

func (natGateway *NATGateway) GetTargetID() string {
	return natGateway.ID
}

func (natGateway *NATGateway) GetRuntimeVariables() map[string]string {
    return natGateway.RuntimeVariables
}

func (natGateway *NATGateway) AddRuntimeVariable(key, value string) {
    natGateway.RuntimeVariables[key] = value
}

func (natGateway *NATGateway) ExportModule() error {
	e := reflect.ValueOf(natGateway).Elem()

	subnet := natGateway.Subnet
	provider := strings.Split(subnet.Kubernetes.Cloud.GetCloudProvider().GetType(), "*k8sml.")[1]
	subnetType := strings.Split(reflect.TypeOf(subnet).String(), "*k8sml.")[1]

	module := terraform.NewModule(natGateway.ID, provider, strings.Split(reflect.TypeOf(natGateway).String(), "*k8sml.")[1])
	for i := 0; i < e.NumField(); i++ {
		key := e.Type().Field(i).Name
		value := e.Field(i).Interface()
		module.AddVariable(key, value)
	}
	module.AddVariable(subnetType, "${module." + subnet.ID + "." + strings.ToLower(subnetType) + "}")
	module.AddVariable("k8sTag", natGateway.Subnet.Kubernetes.ID)

	err := module.Export()

	return err
}