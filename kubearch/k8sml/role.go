package k8sml

import (
  "gopkg.in/yaml.v3"
  "errors"
  "reflect"
  "strings"
  ansible "KubeArch/kubearch/proletarian/ansible"
)

type Role struct {
  ID string
  VirtualMachines []VirtualMachine
  Software []string
  VirtualFirewall VirtualFirewall
  TargetGroup *TargetGroup
}

type tmpRole struct {
  VirtualMachines map[string][]yaml.Node `yaml:",inline"`
  Software []string `yaml:"Software"`
}

func (role *Role) GetID() string {
  return role.ID
}

func (role *Role) GetVariableValue(variable string) interface{} {
	e := reflect.ValueOf(role).Elem()
	var value interface{}

	for i := 0; i < e.NumField(); i++ {
		if strings.ToLower(e.Type().Field(i).Name) == strings.ToLower(variable) {
			value = e.Field(i).Interface()
			return value
		}
	}

	return nil
}

func (role *Role) GetVirtualMachineByID(vmID string) (error, *VirtualMachine) {
  for _, vm := range role.VirtualMachines {
    if vm.GetID() == vmID {
      return nil, &vm
    }
  }

	err := errors.New("Failed to identify the virtual machine with ID \"" + vmID + "\"")
	return err, nil
}

func (role *Role) ExportInventory() error {
  group := ansible.NewInventory(role.ID)

  for _, vm := range role.VirtualMachines {
    group.AddVirtualMachine(vm.GetID())
    vm.ExportInventory(group)
  }

  if err := group.Export(); err != nil {
    return err
  }

  if _, ok := role.VirtualMachines[0].GetRuntimeVariables()["public_ip"]; !ok {
    vars := ansible.NewInventory(role.ID + ":vars")

    jumpHost := role.VirtualMachines[0].GetJumpHosts()[0]
    vars.AddVirtualMachine("ansible_ssh_common_args='")
    vars.AddVariable("ansible_ssh_common_args='", "-o StrictHostKeyChecking", "no -o ProxyCommand=\"ssh -i " + jumpHost.GetKey().Path + " -W %h:%p " + jumpHost.GetImage().User + "@" + strings.TrimSuffix(jumpHost.GetRuntimeVariables()["public_ip"], "\n") + "\"'")

    if err := vars.Export(); err != nil {
      return err
    }
  }

  return nil
}

func (role *Role) UnmarshalYAML(value *yaml.Node) error {
  var tmpRole tmpRole
	
  if err := value.Decode(&tmpRole); err != nil {
      return err
  }

  role.Software = make([]string, 0)

  switch role.ID {
  case "Etcd":
    role.Software = append(role.Software, "docker", "kubeadm", "kubelet", "etcd")
    role.Software = append(role.Software, tmpRole.Software...)
  case "Controlplane":
    role.Software = append(role.Software, "docker", "kubeadm", "kubelet", "controlplane", "kubectl")
    if role.TargetGroup != nil {
      role.Software = append(role.Software, role.TargetGroup.VirtualFirewall.GetSubnet().Kubernetes.ContainerNetworkInterface.ID)
      switch strings.Split(reflect.TypeOf(role.TargetGroup.VirtualFirewall.GetSubnet().Kubernetes.Cloud.GetCloudProvider()).String(), "*k8sml.")[1] {
      case "AmazonWebServices":
        role.Software = append(role.Software, "aws_provider")
      }
    } else {
      role.Software = append(role.Software, role.VirtualFirewall.GetSubnet().Kubernetes.ContainerNetworkInterface.ID)
      switch strings.Split(reflect.TypeOf(role.VirtualFirewall.GetSubnet().Kubernetes.Cloud.GetCloudProvider().GetID()).String(), "*k8sml.")[1] {
      case "AmazonWebServices":
        role.Software = append(role.Software, "aws_provider")
      }
    }
    role.Software = append(role.Software, tmpRole.Software...)
  case "Worker":
    role.Software = append(role.Software, "docker", "kubeadm", "kubelet", "worker")
    role.Software = append(role.Software, tmpRole.Software...)
  default:
    role.Software = tmpRole.Software
  }

  vms := make([]VirtualMachine, 0)

  for tag, nodes := range tmpRole.VirtualMachines {
    for _, node := range nodes {
      switch tag {
      case "EC2Instance":
        instance := &EC2Instance{}
        instance.Role = role

        if err := node.Decode(instance); err != nil {
          return err
        }

        vms = append(vms, instance)
      }
    }
  }

  role.VirtualMachines = vms
  
  return nil
}