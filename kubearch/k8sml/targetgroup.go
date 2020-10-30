package k8sml

import (
	"gopkg.in/yaml.v3"
	"reflect"
	"strings"
	terraform "KubeArch/kubearch/proletarian/terraform"
)

type TargetGroup struct {
	ID string
	Protocol string
	Port string
	RuntimeVariables map[string]string
	Target *Role
	LoadBalancer *NetworkLoadBalancer
	VirtualFirewall VirtualFirewall
}

type tmpTargetGroup struct {
	ID string `yaml:"id"`
	Protocol string `yaml:"protocol"`
	Port string `yaml:"port"`
	LoadBalancer string `yaml:"loadbalancer"`
	Target map[string]yaml.Node `yaml:",inline"`
	SecurityGroup *SecurityGroup
}

func (tg *TargetGroup) GetID() string {
	return tg.ID
}

func (tg *TargetGroup) GetVariableValue(variable string) interface{} {
	e := reflect.ValueOf(tg).Elem()
	var value interface{}

	for i := 0; i < e.NumField(); i++ {
		if strings.ToLower(e.Type().Field(i).Name) == strings.ToLower(variable) {
			value = e.Field(i).Interface()
			return value
		}
	}

	return nil
}

func (tg *TargetGroup) GetRuntimeVariables() map[string]string {
    return tg.RuntimeVariables
}

func (tg *TargetGroup) AddRuntimeVariable(key, value string) {
    tg.RuntimeVariables[key] = value
}

func (tg *TargetGroup) ExportModule() error {
	e := reflect.ValueOf(tg).Elem()

	provider := strings.Split(tg.VirtualFirewall.GetSubnet().Kubernetes.Cloud.GetCloudProvider().GetType(), "*k8sml.")[1]
	cloud := tg.VirtualFirewall.GetSubnet().Kubernetes.Cloud
	cloudType := strings.Split(reflect.TypeOf(cloud).String(), "*k8sml.")[1]

	module := terraform.NewModule(tg.ID, provider, strings.Split(reflect.TypeOf(tg).String(), "*k8sml.")[1])
	for i := 0; i < e.NumField(); i++ {
		key := e.Type().Field(i).Name
		value := e.Field(i).Interface()
		module.AddVariable(key, value)
	}
	module.AddVariable(cloudType, tg.VirtualFirewall.GetSubnet().Kubernetes.Cloud.GetRuntimeVariables()["id"])

	if err := module.Export(); err != nil {
		return err
	}

	for _, target := range tg.Target.VirtualMachines {
		tgType := strings.Split(reflect.TypeOf(tg).String(), "*k8sml.")[1]

		module := terraform.NewModule(target.GetID(), provider, "targetgroupattachment")
		module.AddVariable(tgType, "${module." + tg.ID + "." + strings.ToLower(tgType) + "}")
		module.AddVariable("target", target.GetRuntimeVariables()["private_ip"])
		module.AddVariable("port", tg.Port)

		if err := module.Export(); err != nil {
			return err
		}
	}

	return nil
}

func (tg *TargetGroup) UnmarshalYAML(value *yaml.Node) error {
	var tmpTg tmpTargetGroup
	
    if err := value.Decode(&tmpTg); err != nil {
        return err
	}

	tg.ID = tmpTg.ID
	tg.Protocol = tmpTg.Protocol
	tg.Port = tmpTg.Port

	for tag, node := range tmpTg.Target {
		role := &Role{}
		role.ID = tag
		role.TargetGroup = tg
		
		if err := node.Decode(role); err != nil {
			return nil
		}

		tg.Target = role
	}

	tg.LoadBalancer = &NetworkLoadBalancer{tmpTg.LoadBalancer, "", "", "", nil, nil, nil}

	return nil
}