package k8sml

type VirtualFirewall interface {
	GetID() string
	GetVariableValue(variable string) interface{}
	GetIngress() []*Ingress
	GetEgress() []*Egress
	GetTargetGroups() []*TargetGroup
	GetRoles() []*Role
	GetSubnet() *Subnet
	AddRuntimeVariable(key, value string)
	GetRuntimeVariables() map[string]string
	ExportModule() error
}