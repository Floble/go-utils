package k8sml

import (
	"gopkg.in/yaml.v3"
 	"reflect"
	"strings"
	"errors"
	terraform "KubeArch/kubearch/proletarian/terraform"
)

type Subnet struct {
	ID string
	AvailabilityZone string
	Cidr string
	Public bool
	RuntimeVariables map[string]string
	Kubernetes *Kubernetes
	RouteTable []*RouteTable
	VirtualFirewall []VirtualFirewall
	NATGateway []*NATGateway
	NetworkLoadBalancer []*NetworkLoadBalancer
}

type tmpSubnet struct {
	ID string `yaml:"id"`
	AvailabilityZone string `yaml:"availability_zone"`
	Cidr string `yaml:"cidr"`
	Public bool `yaml:"public"`
	NATGateway []yaml.Node `yaml:"NatGateway"`
	RouteTable []yaml.Node `yaml:"RouteTable"`
	NetworkLoadBalancer []yaml.Node `yaml:"NetworkLoadBalancer"`
	VirtualFirewall map[string][]yaml.Node `yaml:",inline"`
}

func (subnet *Subnet) GetID() string {
	return subnet.ID
}

func (subnet *Subnet) GetVariableValue(variable string) interface{} {
	e := reflect.ValueOf(subnet).Elem()
	var value interface{}

	for i := 0; i < e.NumField(); i++ {
		if strings.ToLower(e.Type().Field(i).Name) == strings.ToLower(variable) {
			value = e.Field(i).Interface()
			return value
		}
	}

	return nil
}

func (subnet *Subnet) GetRuntimeVariables() map[string]string {
    return subnet.RuntimeVariables
}

func (subnet *Subnet) AddRuntimeVariable(key, value string) {
    subnet.RuntimeVariables[key] = value
}

func (subnet *Subnet) ExportModule() error {
	e := reflect.ValueOf(subnet).Elem()

	subnetType := strings.Split(reflect.TypeOf(subnet).String(), "*k8sml.")[1]
	cloud := subnet.Kubernetes.Cloud
	provider := strings.Split(cloud.GetCloudProvider().GetType(), "*k8sml.")[1]
	cloudType := strings.Split(reflect.TypeOf(cloud).String(), "*k8sml.")[1]
	ipv4cidr := subnet.Kubernetes.Cloud.GetIPv4Cidr()
	ipv4cidrType := strings.Split(reflect.TypeOf(ipv4cidr).String(), "*k8sml.")[1]

	module := terraform.NewModule(subnet.ID, provider, strings.Split(reflect.TypeOf(subnet).String(), "*k8sml.")[1])
	for i := 0; i < e.NumField(); i++ {
		key := e.Type().Field(i).Name
		value := e.Field(i).Interface()
		module.AddVariable(key, value)
	}
	module.AddVariable(cloudType, "${module." + cloud.GetID() + "." + strings.ToLower(cloudType) + "}")
	module.AddVariable("k8sTag", subnet.Kubernetes.ID)
	if subnet.Public {
		module.AddVariable("lbTag", "elb")
	} else {
		module.AddVariable("lbTag", "internal-elb")
	}

	associatedCidr, err := getAssociatedCidr(subnet, ipv4cidr)
	if err != nil {
		return err
	}

	if associatedCidr != "" {
		module.AddVariable(ipv4cidrType, "${module." + associatedCidr + "." + strings.ToLower(ipv4cidrType) + "}")
	}

	output := terraform.NewOutput(subnet.ID + "_" + strings.ToLower(subnetType))

    for key, _ := range subnet.RuntimeVariables {
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

func getAssociatedCidr(subnet *Subnet, ipv4cidr []*IPv4Cidr) (string, error) {
	for _, cidr := range ipv4cidr {
		check, err := cidr.Contains(subnet.Cidr)
		if err != nil {
			return "", err
		} else if check {
			return cidr.ID, nil
		}
	}

	return "", nil
}

func (subnet *Subnet) UnmarshalYAML(value *yaml.Node) error {
	var tmpSubnet tmpSubnet
	
	if err := value.Decode(&tmpSubnet); err != nil {
        return err
	}

	subnet.ID = tmpSubnet.ID
	subnet.AvailabilityZone = tmpSubnet.AvailabilityZone
	subnet.Cidr = tmpSubnet.Cidr
	subnet.Public = tmpSubnet.Public
	subnet.RuntimeVariables = make(map[string]string)

	natGateways := make([]*NATGateway, 0)
	natGateway := &NATGateway{}
	natGateway.Subnet = subnet

	for _, node := range tmpSubnet.NATGateway {
		if err := node.Decode(natGateway); err != nil {
			return err
		}

		natGateways = append(natGateways, natGateway)
	}

	subnet.NATGateway = natGateways

	routeTables := make([]*RouteTable, 0)

	for _, node := range tmpSubnet.RouteTable {
		routeTable := &RouteTable{}
		routeTable.Subnet = subnet

		if err := node.Decode(routeTable); err != nil {
			return err
		}

		routeTables = append(routeTables, routeTable)
	}
	
	subnet.RouteTable = routeTables

	nlbs := make([]*NetworkLoadBalancer, 0)

	for _, node := range tmpSubnet.NetworkLoadBalancer {
		nlb := &NetworkLoadBalancer{}
		nlb.Subnet = subnet

		if err := node.Decode(nlb); err != nil {
			return err
		}

		nlbs = append(nlbs, nlb)
	}

	subnet.NetworkLoadBalancer = nlbs

 	virtualFirewalls := make([]VirtualFirewall, 0)

	for tag, nodes := range tmpSubnet.VirtualFirewall {
		switch tag {
		case "SecurityGroup":
			for _, node := range nodes {
				securityGroup := &SecurityGroup{}
				securityGroup.Subnet = subnet

				if err := node.Decode(securityGroup); err != nil {
					return err
				}

				virtualFirewalls = append(virtualFirewalls, securityGroup)
			}
		default:
			return errors.New("Failed to interpret the virtual firewall of type: \"" + tag + "\"")
		}
	}

	subnet.VirtualFirewall = virtualFirewalls

	subnet.AddRuntimeVariable("id", "")
	
    return nil
}