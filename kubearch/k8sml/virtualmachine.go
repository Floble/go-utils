package k8sml

import (
    ansible "KubeArch/kubearch/proletarian/ansible"
)

type VirtualMachine interface {
    SetJumpHosts(jumpHosts []VirtualMachine)
    AddRuntimeVariable(key, value string)
    GetID() string
    GetVariableValue(variable string) interface{}
    GetVirtualMachineRole() *Role
    GetJumpHosts() []VirtualMachine
    GetImage() *Image
    GetKey() *Key
    GetIAMRole() *IAMRole
    GetType() string
    GetRuntimeVariables() map[string]string
    ExportModule() error
    ExportInventory(inventory *ansible.Inventory)
}