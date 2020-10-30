package k8sml

import (
	"reflect"
	"strings"
)

type Key struct {
	ID string `yaml:"id"`
	Path string `yaml:"path"`
	VirtualMachines []VirtualMachine
}

func (key *Key) GetID() string {
	return key.ID
}

func (key *Key) GetVariableValue(variable string) interface{} {
	e := reflect.ValueOf(key).Elem()
	var value interface{}

	for i := 0; i < e.NumField(); i++ {
		if strings.ToLower(e.Type().Field(i).Name) == strings.ToLower(variable) {
			value = e.Field(i).Interface()
			return value
		}
	}

	return nil
}