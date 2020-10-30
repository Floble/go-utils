package k8sml

import (
	"reflect"
	"strings"
)

type Image struct {
	ID string `yaml:"id"`
	User string `yaml:"user"`
	VirtualMachine VirtualMachine
}

func (image *Image) GetID() string {
	return image.ID
}

func (image *Image) GetVariableValue(variable string) interface{} {
	e := reflect.ValueOf(image).Elem()
	var value interface{}

	for i := 0; i < e.NumField(); i++ {
		if strings.ToLower(e.Type().Field(i).Name) == strings.ToLower(variable) {
			value = e.Field(i).Interface()
			return value
		}
	}

	return nil
}