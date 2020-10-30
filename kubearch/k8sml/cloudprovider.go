package k8sml

type CloudProvider interface {
	GetID() string
	GetVariableValue(variable string) interface{}
	GetCloud() []Cloud
	GetType() string
	GetPolicy() []*IAMPolicy
	AddRuntimeVariable(key, value string)
	GetRuntimeVariables() map[string]string
	ExportModule() error
}