package k8sml

type K8sML interface {
	GetID() string
	GetVariableValue(variable string) interface{}
}