package k8sml

import (
	"reflect"
	"strings"
)

type ContainerNetworkInterface struct {
	ID string `yaml:"id"`
	Kubernetes *Kubernetes
}

func (cni *ContainerNetworkInterface) GetID() string {
	return cni.ID
}

func (cni *ContainerNetworkInterface) GetVariableValue(variable string) interface{} {
	e := reflect.ValueOf(cni).Elem()
	var value interface{}

	for i := 0; i < e.NumField(); i++ {
		if strings.ToLower(e.Type().Field(i).Name) == strings.ToLower(variable) {
			value = e.Field(i).Interface()
			return value
		}
	}

	return nil
}