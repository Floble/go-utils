package k8sml

import (
	ansible "KubeArch/kubearch/proletarian/ansible"
	terraform "KubeArch/kubearch/proletarian/terraform"
	"reflect"
	"strings"

	"gopkg.in/yaml.v3"
)

type EC2Instance struct {
	ID, Type         string
	Image            *Image
	RuntimeVariables map[string]string
	Role             *Role
	IAMRole          *IAMRole
	Key              *Key
	JumpHosts        []VirtualMachine
}

type tmpEC2Instance struct {
	ID    string    `yaml:"id"`
	Type  string    `yaml:"type"`
	Image yaml.Node `yaml:"Image"`
	Key   yaml.Node `yaml:"Key"`
}

func (instance *EC2Instance) SetJumpHosts(jumpHosts []VirtualMachine) {
	instance.JumpHosts = jumpHosts
}

func (instance *EC2Instance) GetID() string {
	return instance.ID
}

func (instance *EC2Instance) GetVariableValue(variable string) interface{} {
	e := reflect.ValueOf(instance).Elem()
	var value interface{}

	for i := 0; i < e.NumField(); i++ {
		if strings.ToLower(e.Type().Field(i).Name) == strings.ToLower(variable) {
			value = e.Field(i).Interface()
			return value
		}
	}

	return nil
}

func (instance *EC2Instance) GetTargetID() string {
	return instance.ID
}

func (instance *EC2Instance) GetVirtualMachineRole() *Role {
	return instance.Role
}

func (instance *EC2Instance) GetRuntimeVariables() map[string]string {
	return instance.RuntimeVariables
}

func (instance *EC2Instance) GetJumpHosts() []VirtualMachine {
	return instance.JumpHosts
}

func (instance *EC2Instance) GetImage() *Image {
	return instance.Image
}

func (instance *EC2Instance) GetKey() *Key {
	return instance.Key
}

func (instance *EC2Instance) GetIAMRole() *IAMRole {
	return instance.IAMRole
}

func (instance *EC2Instance) GetType() string {
	return instance.Type
}

func (instance *EC2Instance) AddRuntimeVariable(key, value string) {
	instance.RuntimeVariables[key] = value
}

func (instance *EC2Instance) ExportModule() error {
	e := reflect.ValueOf(instance).Elem()

	var provider string
	if instance.Role.TargetGroup != nil {
		provider = strings.Split(instance.Role.TargetGroup.VirtualFirewall.GetSubnet().Kubernetes.Cloud.GetCloudProvider().GetType(), "*k8sml.")[1]
	} else {
		provider = strings.Split(instance.Role.VirtualFirewall.GetSubnet().Kubernetes.Cloud.GetCloudProvider().GetType(), "*k8sml.")[1]
	}

	var subnet *Subnet
	if instance.Role.TargetGroup != nil {
		subnet = instance.Role.TargetGroup.VirtualFirewall.GetSubnet()
	} else {
		subnet = instance.Role.VirtualFirewall.GetSubnet()
	}
	instanceType := strings.Split(reflect.TypeOf(instance).String(), "*k8sml.")[1]
	subnetType := strings.Split(reflect.TypeOf(subnet).String(), "*k8sml.")[1]

	var scgType string
	if instance.Role.TargetGroup != nil {
		scgType = strings.Split(reflect.TypeOf(instance.Role.TargetGroup.VirtualFirewall).String(), "*k8sml.")[1]
	} else {
		scgType = strings.Split(reflect.TypeOf(instance.Role.VirtualFirewall).String(), "*k8sml.")[1]
	}

	module := terraform.NewModule(instance.ID, provider, strings.Split(reflect.TypeOf(instance).String(), "*k8sml.")[1])
	for i := 0; i < e.NumField(); i++ {
		key := e.Type().Field(i).Name
		value := e.Field(i).Interface()
		module.AddVariable(key, value)
	}
	module.AddVariable("ami", instance.Image.ID)
	module.AddVariable("key", instance.Key.ID)
	module.AddVariable(subnetType, "${module."+subnet.ID+"."+strings.ToLower(subnetType)+"}")
	if instance.Role.TargetGroup != nil {
		module.AddVariable(scgType, "${module."+instance.Role.TargetGroup.VirtualFirewall.GetID()+"."+strings.ToLower(scgType)+"}")
		module.AddVariable("k8stag", instance.Role.TargetGroup.VirtualFirewall.GetSubnet().Kubernetes.ID)
	} else {
		module.AddVariable(scgType, "${module."+instance.Role.VirtualFirewall.GetID()+"."+strings.ToLower(scgType)+"}")
		module.AddVariable("k8stag", instance.Role.VirtualFirewall.GetSubnet().Kubernetes.ID)
	}

	if instance.IAMRole != nil {
		roleType := strings.Split(reflect.TypeOf(instance.IAMRole).String(), "*k8sml.")[1]
		module.AddVariable(roleType, "${module."+instance.IAMRole.ID+"."+strings.ToLower(roleType)+"}")
	}

	output := terraform.NewOutput(instance.ID + "_" + strings.ToLower(instanceType))

	for key, _ := range instance.RuntimeVariables {
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

func (instance *EC2Instance) ExportInventory(inventory *ansible.Inventory) {
	inventory.AddVariable(instance.ID, "user", instance.Image.User)

	if instance.Role.TargetGroup != nil {
		if !instance.Role.TargetGroup.VirtualFirewall.GetSubnet().Public {
			inventory.AddVariable(instance.ID, "public_ip", instance.RuntimeVariables["private_ip"])
		}
	} else {
		if !instance.Role.VirtualFirewall.GetSubnet().Public {
			inventory.AddVariable(instance.ID, "public_ip", instance.RuntimeVariables["private_ip"])
		}
	}

	for variable, value := range instance.RuntimeVariables {
		if variable == "id" {
			variable = "provider_id"
		}
		inventory.AddVariable(instance.ID, variable, value)
	}
}

func (instance *EC2Instance) UnmarshalYAML(value *yaml.Node) error {
	var tmpInstance tmpEC2Instance

	if err := value.Decode(&tmpInstance); err != nil {
		return err
	}

	instance.RuntimeVariables = make(map[string]string, 0)
	instance.ID = tmpInstance.ID
	instance.Type = tmpInstance.Type

	key := &Key{}
	key.VirtualMachines = append(key.VirtualMachines, instance)

	if err := tmpInstance.Key.Decode(key); err != nil {
		return err
	}

	instance.Key = key

	image := &Image{}
	image.VirtualMachine = instance

	if err := tmpInstance.Image.Decode(image); err != nil {
		return err
	}

	instance.Image = image

	instance.AddRuntimeVariable("private_ip", "")
	instance.AddRuntimeVariable("private_dns", "")
	instance.AddRuntimeVariable("id", "")
	if instance.Role.TargetGroup != nil {
		if instance.Role.TargetGroup.VirtualFirewall.GetSubnet().Public {
			instance.AddRuntimeVariable("public_ip", "")
			instance.AddRuntimeVariable("public_dns", "")
		}
	} else {
		if instance.Role.VirtualFirewall.GetSubnet().Public {
			instance.AddRuntimeVariable("public_ip", "")
			instance.AddRuntimeVariable("public_dns", "")
		}
	}

	switch instance.Role.ID {
	case "Controlplane":
		instance.IAMRole = NewIAMRole("aws-cloudprovider-master")
	default:
		instance.IAMRole = NewIAMRole("aws-cloudprovider-worker")
	}

	return nil
}
