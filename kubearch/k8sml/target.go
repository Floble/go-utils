package k8sml

type Target interface {
	GetTargetID() string
	GetVariableValue(variable string) interface{}
}