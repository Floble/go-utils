package k8sml

import (
	"gopkg.in/yaml.v3"
	"reflect"
	"strings"
)

type Kubernetes struct {
	ID, Domain, CidrPod, CidrService, Port string
	ContainerNetworkInterface *ContainerNetworkInterface
	Subnet []*Subnet
	Cloud Cloud
}

type tmpKubernetes struct {
	ID string `yaml:"id"`
	Domain string `yaml:"domain"`
	CidrPod string `yaml:"cidr_pod"`
	CidrService string `yaml:"cidr_service"`
	Port string `yaml:"port"`
	ContainerNetworkInterface yaml.Node `yaml:"ContainerNetworkInterface"`
	Subnet []yaml.Node `yaml:"Subnet"`
}

func (k8s *Kubernetes) GetID() string {
	return k8s.ID
}

func (k8s *Kubernetes) GetVariableValue(variable string) interface{} {
	e := reflect.ValueOf(k8s).Elem()
	var value interface{}

	for i := 0; i < e.NumField(); i++ {
		if strings.ToLower(e.Type().Field(i).Name) == strings.ToLower(variable) {
			value = e.Field(i).Interface()
			return value
		}
	}

	return nil
}

func (k8s *Kubernetes) UnmarshalYAML(value *yaml.Node) error {
	var tmpKubernetes tmpKubernetes

	if err := value.Decode(&tmpKubernetes); err != nil {
		return err
	}

	k8s.ID = tmpKubernetes.ID
	k8s.Port = tmpKubernetes.Port
	k8s.Domain = tmpKubernetes.Domain
	k8s.CidrPod = tmpKubernetes.CidrPod
	k8s.CidrService = tmpKubernetes.CidrService
	
	cni := &ContainerNetworkInterface{}
	cni.Kubernetes = k8s

	if err := tmpKubernetes.ContainerNetworkInterface.Decode(&cni); err != nil {
		return err
	}

	k8s.ContainerNetworkInterface = cni

	var subnets = make([]*Subnet, 0)

 	for _, node := range tmpKubernetes.Subnet {
		subnet := &Subnet{}
		subnet.Kubernetes = k8s

		if err := node.Decode(&subnet); err != nil {
			return err
		}

		subnets = append(subnets, subnet)
	}

	k8s.Subnet = subnets

	return nil
}