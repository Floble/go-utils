package k8sml

type Infrastructure interface {
	GetID() string
	GetVariableValue(variable string) interface{}
	ExportModule() error
    AddRuntimeVariable(key, value string)
    GetRuntimeVariables() map[string]string
}