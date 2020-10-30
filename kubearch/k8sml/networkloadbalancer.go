package k8sml

import (
	"gopkg.in/yaml.v3"
	"reflect"
	"strings"
	terraform "KubeArch/kubearch/proletarian/terraform"
)

type NetworkLoadBalancer struct {
	ID string
	Protocol string
	IP string
	Port string
	RuntimeVariables map[string]string
	TargetGroup *TargetGroup
	Subnet *Subnet
}

type tmpNetworkLoadBalancer struct {
	ID string `yaml:"id"`
	Protocol string `yaml:"protocol"`
	IP string `yaml:"ip"`
	Port string `yaml:"port"`
}

func (nlb *NetworkLoadBalancer) GetID() string {
	return nlb.ID
}

func (nlb *NetworkLoadBalancer) GetVariableValue(variable string) interface{} {
	e := reflect.ValueOf(nlb).Elem()
	var value interface{}

	for i := 0; i < e.NumField(); i++ {
		if strings.ToLower(e.Type().Field(i).Name) == strings.ToLower(variable) {
			value = e.Field(i).Interface()
			return value
		}
	}

	if value == nil {
		if val, ok := nlb.RuntimeVariables[strings.ToLower(variable)]; ok {
			value = val
			return value
		}
	}

	return nil
}

func (nlb *NetworkLoadBalancer) GetRuntimeVariables() map[string]string {
    return nlb.RuntimeVariables
}

func (nlb *NetworkLoadBalancer) AddRuntimeVariable(key, value string) {
    nlb.RuntimeVariables[key] = value
}

func (nlb *NetworkLoadBalancer) ExportModule() error {
	e := reflect.ValueOf(nlb).Elem()

	provider := strings.Split(nlb.Subnet.Kubernetes.Cloud.GetCloudProvider().GetType(), "*k8sml.")[1]
	subnet := nlb.Subnet
	nlbType := strings.Split(reflect.TypeOf(nlb).String(), "*k8sml.")[1]
	subnetType := strings.Split(reflect.TypeOf(subnet).String(), "*k8sml.")[1]
	tgType := strings.Split(reflect.TypeOf(nlb.TargetGroup).String(), "*k8sml.")[1]

	module := terraform.NewModule(nlb.ID, provider, strings.Split(reflect.TypeOf(nlb).String(), "*k8sml.")[1])
	for i := 0; i < e.NumField(); i++ {
		key := e.Type().Field(i).Name
		value := e.Field(i).Interface()
		module.AddVariable(key, value)
	}
	switch nlb.Subnet.Public {
	case true:
		module.AddVariable("private", "false")
	case false:
		module.AddVariable("private", "true")
	}
	module.AddVariable(subnetType, subnet.GetRuntimeVariables()["id"])
	module.AddVariable(tgType, "${module." + nlb.TargetGroup.ID + "." + strings.ToLower(tgType) + "}")

	output := terraform.NewOutput(nlb.ID + "_" + strings.ToLower(nlbType))

    for key, _ := range nlb.RuntimeVariables {
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

func (nlb *NetworkLoadBalancer) UnmarshalYAML(value *yaml.Node) error {
	var tmpNlb tmpNetworkLoadBalancer
	
    if err := value.Decode(&tmpNlb); err != nil {
        return err
	}

	nlb.RuntimeVariables = make(map[string]string, 0)
	nlb.ID = tmpNlb.ID
	nlb.Protocol = tmpNlb.Protocol
	nlb.IP = tmpNlb.IP
	nlb.Port = tmpNlb.Port

	nlb.AddRuntimeVariable("dns_name", "")

	return nil
}