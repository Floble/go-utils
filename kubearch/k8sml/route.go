package k8sml

import (
	"gopkg.in/yaml.v3"
	"reflect"
	"strings"
	terraform "KubeArch/kubearch/proletarian/terraform"
)

type Route struct {
	ID, Cidr string
	RuntimeVariables map[string]string
	Target Target
	RouteTable *RouteTable
}

type tmpRoute struct {
	ID string `yaml:"id"`
	Cidr string `yaml:"cidr"`
	InternetGateway string `yaml:"internetgateway"`
	NatGateway string `yaml:"natgateway"`
	EC2Instance string `yaml:"ec2instance"`
}

func (route *Route) GetID() string {
	return route.ID
}

func (route *Route) GetVariableValue(variable string) interface{} {
	e := reflect.ValueOf(route).Elem()
	var value interface{}

	for i := 0; i < e.NumField(); i++ {
		if strings.ToLower(e.Type().Field(i).Name) == strings.ToLower(variable) {
			value = e.Field(i).Interface()
			return value
		}
	}

	return nil
}

func (route *Route) GetRuntimeVariables() map[string]string {
    return route.RuntimeVariables
}

func (route *Route) AddRuntimeVariable(key, value string) {
    route.RuntimeVariables[key] = value
}

func (route *Route) ExportModule() error {
	e := reflect.ValueOf(route).Elem()

	provider := strings.Split(route.RouteTable.Subnet.Kubernetes.Cloud.GetCloudProvider().GetType(), "*k8sml.")[1]
	routeTableType := strings.Split(reflect.TypeOf(route.RouteTable).String(), "*k8sml.")[1]
	targetType := strings.Split(reflect.TypeOf(route.Target).String(), "*k8sml.")[1]
	var targetVariable string

	switch targetType {
	case "InternetGateway":
		targetVariable = "Gateway"
	case "NATGateway":
		targetVariable = "NatGateway"
	}

	module := terraform.NewModule(route.ID, provider, strings.Split(reflect.TypeOf(route).String(), "*k8sml.")[1])
	for i := 0; i < e.NumField(); i++ {
		key := e.Type().Field(i).Name
		value := e.Field(i).Interface()
		module.AddVariable(key, value)
	}

	module.AddVariable(routeTableType, "${module." + route.RouteTable.ID + "." + strings.ToLower(routeTableType) + "}")
	module.AddVariable(targetVariable, "${module." + route.Target.GetTargetID() + "." + strings.ToLower(targetType) + "}")

	err := module.Export()

	return err
}

func (route *Route) UnmarshalYAML(value *yaml.Node) error {
	var tmpRoute tmpRoute
	
    if err := value.Decode(&tmpRoute); err != nil {
        return err
	}

	route.ID = tmpRoute.ID
	route.Cidr = tmpRoute.Cidr
	
 	if tmpRoute.InternetGateway != "" {
		route.Target = &InternetGateway{tmpRoute.InternetGateway, nil, nil}
	} else if tmpRoute.NatGateway != "" {
		route.Target = &NATGateway{tmpRoute.NatGateway, nil, nil}
	}
	
    return nil
}