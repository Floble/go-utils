package k8sml

import (
	"gopkg.in/yaml.v3"
	"reflect"
	"strings"
	terraform "KubeArch/kubearch/proletarian/terraform"
)

type RouteTable struct {
	ID string
	RuntimeVariables map[string]string
	Subnet *Subnet
	Route []*Route
}

type tmpRouteTable struct {
	ID string `yaml:"id"`
	Route []yaml.Node `yaml:"Route"`
}

func (routeTable *RouteTable) GetID() string {
	return routeTable.ID
}

func (routeTable *RouteTable) GetVariableValue(variable string) interface{} {
	e := reflect.ValueOf(routeTable).Elem()
	var value interface{}

	for i := 0; i < e.NumField(); i++ {
		if strings.ToLower(e.Type().Field(i).Name) == strings.ToLower(variable) {
			value = e.Field(i).Interface()
			return value
		}
	}

	return nil
}

func (routeTable *RouteTable) GetRuntimeVariables() map[string]string {
    return routeTable.RuntimeVariables
}

func (routeTable *RouteTable) AddRuntimeVariable(key, value string) {
    routeTable.RuntimeVariables[key] = value
}

func (routeTable *RouteTable) ExportModule() error {
	e := reflect.ValueOf(routeTable).Elem()

	cloud := routeTable.Subnet.Kubernetes.Cloud
	provider := strings.Split(cloud.GetCloudProvider().GetType(), "*k8sml.")[1]
	cloudType := strings.Split(reflect.TypeOf(cloud).String(), "*k8sml.")[1]
	subnet := routeTable.Subnet
	subnetType := strings.Split(reflect.TypeOf(subnet).String(), "*k8sml.")[1]

	module := terraform.NewModule(routeTable.ID, provider, strings.Split(reflect.TypeOf(routeTable).String(), "*k8sml.")[1])
	for i := 0; i < e.NumField(); i++ {
		key := e.Type().Field(i).Name
		value := e.Field(i).Interface()
		module.AddVariable(key, value)
	}
	module.AddVariable(cloudType, "${module." + cloud.GetID() + "." + strings.ToLower(cloudType) + "}")
	module.AddVariable(subnetType, "${module." + subnet.ID + "." + strings.ToLower(subnetType) + "}")
	module.AddVariable("k8sTag", routeTable.Subnet.Kubernetes.ID)

	err := module.Export()

	return err
}

func (routeTable *RouteTable) UnmarshalYAML(value *yaml.Node) error {
	var tmpRouteTable tmpRouteTable
	
	if err := value.Decode(&tmpRouteTable); err != nil {
        return err
	}

	routeTable.ID = tmpRouteTable.ID

	routes := make([]*Route, 0)
	
	for _, node := range tmpRouteTable.Route {
		route := &Route{}
		route.RouteTable = routeTable

		if err := node.Decode(route); err != nil {
			return err
		}

		routes = append(routes, route)
	}
	
	routeTable.Route = routes

	return nil
}