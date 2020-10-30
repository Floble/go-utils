package k8sml

type Cloud interface {
	GetID() string
	GetVariableValue(variable string) interface{}
	GetCloudProvider() CloudProvider
	GetIPv4Cidr() []*IPv4Cidr
	GetInternetGateway() *InternetGateway
	GetKubernetes() *Kubernetes
	AddRuntimeVariable(key, value string)
	GetRuntimeVariables() map[string]string
	ExportModule() error
}