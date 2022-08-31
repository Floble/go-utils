# Introduction to YUMA
## Requirements
In order to execute the instantiation of YUMA you must install the following tools:

- `AWS CLI` - See the following link for installation instructions: [Installation Guide](https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html)
- `Ansible` - See the following link for installation instructions: [Installation Guide](https://docs.ansible.com/ansible/latest/installation_guide/intro_installation.html)
- `Go` - See the following link for installation instructions: [Installation Guide](https://go.dev/doc/install)
- `Git` - See the following link for installation instructions: [Installation Guide](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)

## Setup

The following steps must be performed to be able to execute the instantiation of YUMA:

- Clone this repository to your local machine
- Create a `<testing>` directory on your local machine and place all of your Ansible roles in it. E.g., you can use the Ansible rules from the following repository: [K8s example](https://github.com/Floble/ansible-utils/tree/main/k8s_example)
- Specify the `sigma` and `omega` parameters for the `BuildSearchTree` function. Please note the the parameters `state` and `depth` be defined as `0`
- Specify the `target` parameter for the `DetermineExecutionOrder_UniformCostSearch` function. This parameter represents the binary configuration of the Ansible role that must be composed from the rest of the Ansible roles. The binary configurations reflect the ordering of the Ansible roles in your local `<testing>` directory. The Ansible role placed first in the `<testing>` directory has the binary configuration of `1`. E.g., for the K8s example the binary configurations are as following:

    - configVM - 1
    - deployPod - 2
    - installDocker - 4
    - installKubernetes - 8
    - runKubernetes - 16
- Complile the `main.go` file to create the required binary file

## Execution

The following stepts must be performed to execut the instantiation of YUMA:

- Copy the binary file to the same directory which stores your `<testing>` directory
- Execute the binary file
- Check the `playbook.yml` and `result.txt` produced by the instantiation of YUMA in your local `<testing>` directory. The `playbook.yml` describes your composition plan as a Ansible playbook. The `result.txt` describes the resulting search tree and model of the agent

## Logical Architecture

The logical architecture of the instantiation of YUMA is described by the following diagram:

![logical_architecture](https://github.com/Floble/go-utils/blob/ucs/algorithms/artificialintelligence/search/logical_architecture.png)