package k8sml

import (
	"gopkg.in/yaml.v3"
	"reflect"
	"strings"
	terraform "KubeArch/kubearch/proletarian/terraform"
)

type SecurityGroup struct {
	ID string
	RuntimeVariables map[string]string
	Subnet *Subnet
	Ingress []*Ingress
	Egress []*Egress
	Roles []*Role
	TargetGroups []*TargetGroup
}

type tmpSecurityGroup struct {
	ID string `yaml:"id"`
	Ingress []yaml.Node `yaml:"Ingress"`
	Egress []yaml.Node `yaml:"Egress"`
	Roles map[string]yaml.Node `yaml:",inline"`
	TargetGroups []yaml.Node `yaml:"TargetGroup"`
}

func (sg *SecurityGroup) GetSubnet() *Subnet {
	return sg.Subnet
}

func (sg *SecurityGroup) GetID() string {
	return sg.ID
}

func (sg *SecurityGroup) GetVariableValue(variable string) interface{} {
	e := reflect.ValueOf(sg).Elem()
	var value interface{}

	for i := 0; i < e.NumField(); i++ {
		if strings.ToLower(e.Type().Field(i).Name) == strings.ToLower(variable) {
			value = e.Field(i).Interface()
			return value
		}
	}

	return nil
}

func (sg *SecurityGroup) GetIngress() []*Ingress {
	return sg.Ingress
}

func (sg *SecurityGroup) GetEgress() []*Egress {
	return sg.Egress
}

func (sg *SecurityGroup) GetTargetGroups() []*TargetGroup {
	return sg.TargetGroups
}

func (sg *SecurityGroup) GetRoles() []*Role {
	return sg.Roles
}

func (sg *SecurityGroup) GetRuntimeVariables() map[string]string {
    return sg.RuntimeVariables
}

func (sg *SecurityGroup) AddRuntimeVariable(key, value string) {
    sg.RuntimeVariables[key] = value
}

func (sg *SecurityGroup) ExportModule() error {
	e := reflect.ValueOf(sg).Elem()

	cloud := sg.Subnet.Kubernetes.Cloud
	provider := strings.Split(cloud.GetCloudProvider().GetType(), "*k8sml.")[1]
	cloudType := strings.Split(reflect.TypeOf(cloud).String(), "*k8sml.")[1]

	module := terraform.NewModule(sg.ID, provider, strings.Split(reflect.TypeOf(sg).String(), "*k8sml.")[1])
	for i := 0; i < e.NumField(); i++ {
		key := e.Type().Field(i).Name
		value := e.Field(i).Interface()
		module.AddVariable(key, value)
	}
	module.AddVariable(cloudType, "${module." + cloud.GetID() + "." + strings.ToLower(cloudType) + "}")
	module.AddVariable("k8stag", sg.Subnet.Kubernetes.ID)

	err := module.Export()

	return err
}

 func (sg *SecurityGroup) UnmarshalYAML(value *yaml.Node) error {
	var tmpSg tmpSecurityGroup
	
    if err := value.Decode(&tmpSg); err != nil {
        return err
	}

	sg.ID = tmpSg.ID

	ingress := make([]*Ingress, 0)

	for _, node := range tmpSg.Ingress {
		rule := &Ingress{}
		rule.VirtualFirewall = sg

		if err := node.Decode(&rule); err != nil {
			return err
		}

		ingress = append(ingress, rule)
	}

	sg.Ingress = ingress

	egress := make([]*Egress, 0)

	for _, node := range tmpSg.Egress {
		rule := &Egress{}
		rule.VirtualFirewall = sg

		if err := node.Decode(&rule); err != nil {
			return err
		}

		egress = append(egress, rule)
	}

	sg.Egress = egress

	roles := make([]*Role, 0)

	for tag, node := range tmpSg.Roles {
		role := &Role{}
		role.ID = tag
		role.VirtualFirewall = sg
		
		if err := node.Decode(role); err != nil {
			return nil
		}

		roles = append(roles, role)
	}

	sg.Roles = roles

	targetGroups := make([]*TargetGroup, 0)

	for _, node := range tmpSg.TargetGroups {
		targetGroup := &TargetGroup{}
		targetGroup.VirtualFirewall = sg

		if err := node.Decode(targetGroup); err != nil {
			return err
		}

		targetGroups = append(targetGroups, targetGroup)
	}

	sg.TargetGroups = targetGroups

	return nil
}