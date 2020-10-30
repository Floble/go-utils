package k8sml

import (
	"gopkg.in/yaml.v3"
	"reflect"
	"strings"
	terraform "KubeArch/kubearch/proletarian/terraform"
)

type Egress struct {
	ID string
	FromPort string
	ToPort string
	Protocol string
	Cidr []string
	RuntimeVariables map[string]string
	VirtualFirewall VirtualFirewall
}

type tmpEgress struct {
	ID string `yaml:"id"`
	FromPort string `yaml:"from_port"`
	ToPort string `yaml:"to_port"`
	Protocol string `yaml:"protocol"`
	Cidr []string `yaml:"cidr"`
	SecurityGroup string `yaml:"security_group"`
}

func (egress *Egress) GetID() string {
	return egress.ID
}

func (egress *Egress) GetVariableValue(variable string) interface{} {
	e := reflect.ValueOf(egress).Elem()
	var value interface{}

	for i := 0; i < e.NumField(); i++ {
		if strings.ToLower(e.Type().Field(i).Name) == strings.ToLower(variable) {
			value = e.Field(i).Interface()
			return value
		}
	}

	return nil
}

func (egress *Egress) GetRuntimeVariables() map[string]string {
    return egress.RuntimeVariables
}

func (egress *Egress) AddRuntimeVariable(key, value string) {
    egress.RuntimeVariables[key] = value
}

func (egress *Egress) ExportModule() error {
	e := reflect.ValueOf(egress).Elem()

	provider := strings.Split(egress.VirtualFirewall.GetSubnet().Kubernetes.Cloud.GetCloudProvider().GetType(), "*k8sml.")[1]
	vfType := strings.Split(reflect.TypeOf(egress.VirtualFirewall).String(), "*k8sml.")[1]

	module := terraform.NewModule(egress.ID, provider, strings.Split(reflect.TypeOf(egress).String(), "*k8sml.")[1])
	for i := 0; i < e.NumField(); i++ {
		key := e.Type().Field(i).Name
		value := e.Field(i).Interface()
		module.AddVariable(key, value)
	}
	module.AddVariable(vfType, "${module." + egress.VirtualFirewall.GetID() + "." + strings.ToLower(vfType) + "}")

	err := module.Export()

	return err
}

func (egress *Egress) UnmarshalYAML(value *yaml.Node) error {
	var tmpEgress tmpEgress
	
    if err := value.Decode(&tmpEgress); err != nil {
        return err
	}

	egress.ID = tmpEgress.ID
	egress.FromPort = tmpEgress.FromPort
	egress.ToPort = tmpEgress.ToPort
	egress.Protocol = tmpEgress.Protocol
	egress.Cidr = tmpEgress.Cidr

	return nil
}