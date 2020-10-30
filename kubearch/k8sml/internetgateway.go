package k8sml

import (
	"reflect"
	"strings"
	terraform "KubeArch/kubearch/proletarian/terraform"
)

type InternetGateway struct {
	ID string `yaml:"id"`
	RuntimeVariables map[string]string
	Cloud Cloud
}

func (igw *InternetGateway) GetID() string {
	return igw.ID
}

func (igw *InternetGateway) GetVariableValue(variable string) interface{} {
	e := reflect.ValueOf(igw).Elem()
	var value interface{}

	for i := 0; i < e.NumField(); i++ {
		if strings.ToLower(e.Type().Field(i).Name) == strings.ToLower(variable) {
			value = e.Field(i).Interface()
			return value
		}
	}

	return nil
}

func (internetGateway *InternetGateway) GetTargetID() string {
	return internetGateway.ID
}

func (igw *InternetGateway) GetRuntimeVariables() map[string]string {
    return igw.RuntimeVariables
}

func (igw *InternetGateway) AddRuntimeVariable(key, value string) {
    igw.RuntimeVariables[key] = value
}

func (internetGateway *InternetGateway) ExportModule() error {
	e := reflect.ValueOf(internetGateway).Elem()

	cloud := internetGateway.Cloud
	provider := strings.Split(cloud.GetCloudProvider().GetType(), "*k8sml.")[1]
	cloudType := strings.Split(reflect.TypeOf(cloud).String(), "*k8sml.")[1]

	module := terraform.NewModule(internetGateway.ID, provider, strings.Split(reflect.TypeOf(internetGateway).String(), "*k8sml.")[1])
	for i := 0; i < e.NumField(); i++ {
		key := e.Type().Field(i).Name
		value := e.Field(i).Interface()
		module.AddVariable(key, value)
	}
	module.AddVariable(cloudType, "${module." + cloud.GetID() + "." + strings.ToLower(cloudType) + "}")
	module.AddVariable("k8sTag", internetGateway.Cloud.GetKubernetes().ID)

	err := module.Export()

	return err
}