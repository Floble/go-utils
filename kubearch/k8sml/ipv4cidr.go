package k8sml

import (
	"reflect"
	"strings"
	"net"
	terraform "KubeArch/kubearch/proletarian/terraform"
)

type IPv4Cidr struct {
	Cloud Cloud
	ID string `yaml:"id"`
	Cidr string `yaml:"cidr"`
	RuntimeVariables map[string]string
}

func (cidr *IPv4Cidr) GetID() string {
	return cidr.ID
}

func (cidr *IPv4Cidr) GetVariableValue(variable string) interface{} {
	e := reflect.ValueOf(cidr).Elem()
	var value interface{}

	for i := 0; i < e.NumField(); i++ {
		if strings.ToLower(e.Type().Field(i).Name) == strings.ToLower(variable) {
			value = e.Field(i).Interface()
			return value
		}
	}

	return nil
}

func (cidr *IPv4Cidr) Contains(check string) (bool, error) {
	_, net1, err := net.ParseCIDR(cidr.Cidr)
	if err != nil {
		return false, err
	}

	_, net2, err := net.ParseCIDR(check)
	if err != nil {
		return false, err
	}

	if net1.Contains(net2.IP) {
		return true, nil
	}

	return false, nil
}

func (cidr *IPv4Cidr) GetRuntimeVariables() map[string]string {
    return cidr.RuntimeVariables
}

func (cidr *IPv4Cidr) AddRuntimeVariable(key, value string) {
    cidr.RuntimeVariables[key] = value
}

func (cidr *IPv4Cidr) ExportModule() error {
	e := reflect.ValueOf(cidr).Elem()

	cloud := cidr.Cloud
	provider := strings.Split(cloud.GetCloudProvider().GetType(), "*k8sml.")[1]
	cloudType := strings.Split(reflect.TypeOf(cloud).String(), "*k8sml.")[1]

	module := terraform.NewModule(cidr.ID, provider, strings.Split(reflect.TypeOf(cidr).String(), "*k8sml.")[1])
	for i := 0; i < e.NumField(); i++ {
		key := e.Type().Field(i).Name
		value := e.Field(i).Interface()
		module.AddVariable(key, value)
	}
	module.AddVariable(cloudType, "${module." + cloud.GetID() + "." + strings.ToLower(cloudType) + "}")

	err := module.Export()

	return err
}