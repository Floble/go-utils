package k8sml

import (
	"gopkg.in/yaml.v3"
	"strings"
	"reflect"
	terraform "KubeArch/kubearch/proletarian/terraform"
)

type Ingress struct {
	ID string
	FromPort string
	ToPort string
	Protocol string
	Cidr []string
	RuntimeVariables map[string]string
	VirtualFirewall VirtualFirewall
}

type tmpIngress struct {
	ID string `yaml:"id"`
	FromPort string `yaml:"from_port"`
	ToPort string `yaml:"to_port"`
	Protocol string `yaml:"protocol"`
	Cidr []string `yaml:"cidr"`
	SecurityGroup string `yaml:"security_group"`
}

func (ingress *Ingress) GetID() string {
	return ingress.ID
}

func (ingress *Ingress) GetVariableValue(variable string) interface{} {
	e := reflect.ValueOf(ingress).Elem()
	var value interface{}

	for i := 0; i < e.NumField(); i++ {
		if strings.ToLower(e.Type().Field(i).Name) == strings.ToLower(variable) {
			value = e.Field(i).Interface()
			return value
		}
	}

	return nil
}

func (ingress *Ingress) GetRuntimeVariables() map[string]string {
    return ingress.RuntimeVariables
}

func (ingress *Ingress) AddRuntimeVariable(key, value string) {
    ingress.RuntimeVariables[key] = value
}

func (ingress *Ingress) ExportModule() error {
	e := reflect.ValueOf(ingress).Elem()
	
	provider := strings.Split(ingress.VirtualFirewall.GetSubnet().Kubernetes.Cloud.GetCloudProvider().GetType(), "*k8sml.")[1]
	vfType := strings.Split(reflect.TypeOf(ingress.VirtualFirewall).String(), "*k8sml.")[1]

	module := terraform.NewModule(ingress.ID, provider, strings.Split(reflect.TypeOf(ingress).String(), "*k8sml.")[1])
	for i := 0; i < e.NumField(); i++ {
		key := e.Type().Field(i).Name
		value := e.Field(i).Interface()
		module.AddVariable(key, value)
	}
	module.AddVariable(vfType, "${module." + ingress.VirtualFirewall.GetID() + "." + strings.ToLower(vfType) + "}")

	err := module.Export()

	return err
}

func (ingress *Ingress) UnmarshalYAML(value *yaml.Node) error {
	var tmpIngress tmpIngress
	
    if err := value.Decode(&tmpIngress); err != nil {
        return err
	}

	ingress.ID = tmpIngress.ID
	ingress.FromPort = tmpIngress.FromPort
	ingress.ToPort = tmpIngress.ToPort
	ingress.Protocol = tmpIngress.Protocol
	ingress.Cidr = tmpIngress.Cidr

	return nil
}