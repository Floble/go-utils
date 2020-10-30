package k8sml

import (
	"gopkg.in/yaml.v3"
	 "reflect"
	 "strings"
	terraform "KubeArch/kubearch/proletarian/terraform" 
)

type VirtualPrivateCloud struct {
	ID, Cidr string
	RuntimeVariables map[string]string
	InternetGateway *InternetGateway
	IPv4Cidr []*IPv4Cidr
	CloudProvider CloudProvider
	Kubernetes *Kubernetes
}

type tmpVirtualPrivateCloud struct {
	ID string `yaml:"id"`
	Cidr string `yaml:"cidr"`
	InternetGateway yaml.Node `yaml:"InternetGateway"`
	IPv4Cidr []yaml.Node `yaml:"IPv4Cidr"`
	Kubernetes yaml.Node `yaml:"Kubernetes"`
}

func (vpc *VirtualPrivateCloud) GetID() string {
	return vpc.ID
}

func (vpc *VirtualPrivateCloud) GetVariableValue(variable string) interface{} {
	e := reflect.ValueOf(vpc).Elem()
	var value interface{}

	for i := 0; i < e.NumField(); i++ {
		if strings.ToLower(e.Type().Field(i).Name) == strings.ToLower(variable) {
			value = e.Field(i).Interface()
			return value
		}
	}

	return nil
}

func (vpc *VirtualPrivateCloud) GetCloudProvider() CloudProvider {
	return vpc.CloudProvider
}

func (vpc *VirtualPrivateCloud) GetRuntimeVariables() map[string]string {
    return vpc.RuntimeVariables
}

func (vpc *VirtualPrivateCloud) AddRuntimeVariable(key, value string) {
    vpc.RuntimeVariables[key] = value
}

func (vpc *VirtualPrivateCloud) ExportModule() error {
	vpcType := strings.Split(reflect.TypeOf(vpc).String(), "*k8sml.")[1]
	module := terraform.NewModule(vpc.ID, strings.Split(vpc.CloudProvider.GetType(), "*k8sml.")[1], strings.Split(reflect.TypeOf(vpc).String(), "*k8sml.")[1])

	e := reflect.ValueOf(vpc).Elem()
	
	for i := 0; i < e.NumField(); i++ {
		key := e.Type().Field(i).Name
		value := e.Field(i).Interface()
		module.AddVariable(key, value)
	}
	module.AddVariable("k8sTag", vpc.GetKubernetes().ID)

	output := terraform.NewOutput(vpc.ID + "_" + strings.ToLower(vpcType))

    for key, _ := range vpc.RuntimeVariables {
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

func (vpc *VirtualPrivateCloud) GetKubernetes() *Kubernetes {
	return vpc.Kubernetes
}

func (vpc *VirtualPrivateCloud) GetIPv4Cidr() []*IPv4Cidr {
	return vpc.IPv4Cidr
}

func (vpc *VirtualPrivateCloud) GetInternetGateway() *InternetGateway {
	return vpc.InternetGateway
}

func (vpc *VirtualPrivateCloud) UnmarshalYAML(value *yaml.Node) error {
	var tmpVirtualPrivateCloud tmpVirtualPrivateCloud
	
    if err := value.Decode(&tmpVirtualPrivateCloud); err != nil {
        return err
	}

	vpc.ID = tmpVirtualPrivateCloud.ID
	vpc.Cidr = tmpVirtualPrivateCloud.Cidr
	vpc.RuntimeVariables = make(map[string]string)

	internetGateway := &InternetGateway{}
	internetGateway.Cloud = vpc

	if err := tmpVirtualPrivateCloud.InternetGateway.Decode(&internetGateway); err != nil {
		return err
	}

	vpc.InternetGateway = internetGateway

	ipv4cidrs := make([]*IPv4Cidr, 0)

	for _, node := range tmpVirtualPrivateCloud.IPv4Cidr {
		ipv4cidr := &IPv4Cidr{}
		ipv4cidr.Cloud = vpc

		if err := node.Decode(ipv4cidr); err != nil {
			return err
		}
			
		ipv4cidrs = append(ipv4cidrs, ipv4cidr)
	}

	vpc.IPv4Cidr = ipv4cidrs

	k8s := &Kubernetes{}
	k8s.Cloud = vpc

	if err := tmpVirtualPrivateCloud.Kubernetes.Decode(&k8s); err != nil {
		return err
	}

	vpc.Kubernetes = k8s

	vpc.AddRuntimeVariable("id", "")
	
    return nil
}