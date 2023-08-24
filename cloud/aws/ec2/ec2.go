package ec2

import (
	"fmt"
	"go-utils/algorithms/artificialintelligence/agents/yuma"
	"go-utils/helper"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type EC2 struct {
	yuma *yuma.Yuma
	executor yuma.Executor
	//instance yuma.Instance
	instances sync.Map
	omega int
	sigma int
}

func NewEC2(yuma *yuma.Yuma, executor yuma.Executor, omega int, sigma int) *EC2 {
	ec2 := new(EC2)
	ec2.yuma = yuma
	ec2.executor = executor
	ec2.omega = omega
	ec2.sigma = sigma

	return ec2
}

func (ec2 *EC2) GetYuma() *yuma.Yuma {
	return ec2.yuma
}

func (ec2 *EC2) GetExecutor() yuma.Executor {
	return ec2.executor
}

func (ec2 *EC2) GetInstances(target int) map[int][]yuma.Instance {
	if instances, ok := ec2.instances.Load(target); ok {
		if instances == nil {
			return nil
		} else {
			return instances.(map[int][]yuma.Instance)
		}
	} else {
		return nil
	}
}

func (ec2 *EC2) SetInstances(target int, instances map[int][]yuma.Instance) {
	ec2.instances.Store(target, instances)
}

func (ec2 *EC2) GetOmega() int {
	return ec2.omega
}

func (ec2 *EC2) GetSigma() int {
	return ec2.sigma
}

func (ec2 *EC2) Initialize() error {
	return nil
}

func (ec2 *EC2) CleanUp() error {
	return nil
}

func (ec2 *EC2) CleanResults() error {
	for i := 0; i < len(ec2.GetYuma().GetSubprocesses()); i++ {
		if _, err := os.Stat("logs_" + strconv.Itoa(int(math.Exp2(float64(i)))) + ".txt"); err == nil {
			if err := os.Remove("logs_" + strconv.Itoa(int(math.Exp2(float64(i)))) + ".txt"); err != nil {
				return err
			}
		}
		if err := os.Remove("memory_" + strconv.Itoa(int(math.Exp2(float64(i)))) + ".txt"); err != nil {
			return err
		}
		if err := os.Remove("playbook_" + strconv.Itoa(int(math.Exp2(float64(i)))) + ".yml"); err != nil {
			return err
		}
		if err := os.Remove("results_" + strconv.Itoa(int(math.Exp2(float64(i)))) + ".txt"); err != nil {
			return err
		}
	}

	return nil
}

func (ec2 *EC2) CreateInstance(target int, action int, waitingTime int) error {	
	instance := NewEC2Instance()

	created := false
	for created == false {
		if err := instance.Create(); err != nil {
			fmt.Println("EC2 ERROR: CREATE INSTANCE")
			//return err
			fmt.Println(err)
		} else {
			created = true
		}
	}

	/* if err := ec2.executor.CreateEnvironmentDescription(target, instance.GetPublicIP()); err != nil {
		fmt.Println("EC2 ERROR: CREATE ENVIRONMENT FROM DESCRIPTION")
		return err
	} */

	time.Sleep(time.Duration(waitingTime) * time.Second)
	
	added := false
	for added == false {
		if err := instance.AddToKnownHosts(); err != nil {
			fmt.Println("EC2 ERROR: ADD TO KNOWN HOSTS")
			//return err
			fmt.Println(err)
		} else {
			added = true
		}
	}

	var tmp []yuma.Instance
	instances := ec2.GetInstances(target)
	if instances == nil {
		instances = make(map[int][]yuma.Instance, 0)
		tmp = make([]yuma.Instance, 0)
		tmp = append(tmp, instance)
		instances[action] = tmp
	} else {
		if tmp, ok := instances[action]; ok {
			instances[action] = append(instances[action], instance)
		} else {
			tmp = make([]yuma.Instance, 0)
			tmp = append(tmp, instance)
			instances[action] = tmp
		}
	}
	ec2.SetInstances(target, instances)

	return nil
}

func (ec2 *EC2) DeleteInstance(target int, action int) error {
	instances := ec2.GetInstances(target)
	instance := instances[action][len(instances[action]) - 1]
	
	deleted := false
	for deleted == false {
		if err := instance.Delete(); err != nil {
			fmt.Println("EC2 ERROR: DELETE INSTANCE")
			//return err
			fmt.Println(err)
		} else {
			deleted = true
		}
	}

	if err := ec2.executor.RemoveEnvironmentDescription(target); err != nil {
		fmt.Println("EC2 ERROR: REMOVE ENVIRONMENT DESCRIPTION")
		return err
	}

	instances[action] = instances[action][:len(instances[action]) - 1]
	ec2.SetInstances(target, instances)

	return nil
}

func (ec2 *EC2) DeleteAllInstances(target int) error {
	instances := ec2.GetInstances(target)
	for action := range instances {
		for _, instance := range instances[action] {
			if err := instance.Delete(); err != nil {
				return err
			}
		}
		instances[action] = make([]yuma.Instance, 0)
	}

	return nil
}

func (ec2 *EC2) TakeAction(target int, state int, action int, path []string, success bool) (error, bool, float64, int) {
	reward := math.MaxFloat64 * -1.0
	successor := -1

	switch ec2.GetYuma().GetMode() {
	case 0:
		for i := 0; i < ec2.GetSigma(); i++ {
			instances := ec2.GetInstances(target)
			tmp, ok := instances[0]
			if !ok || (len(tmp) < 1) {
				if err := ec2.CreateInstance(target, 0, ec2.GetOmega()); err != nil {
					return err, false, math.MaxFloat64 * -1.0, -1
				}
			}
			
			instances = ec2.GetInstances(target)
			tmp = instances[0]
			publicIps := make([]string, 0)
			for i := len(tmp) - 1; i > (len(tmp) - 1) - ec2.GetYuma().GetQuantities()[action]; i-- {
				publicIps = append(publicIps, tmp[i].GetPublicIP())
			}
			if err := ec2.executor.CreateEnvironmentDescription(target, publicIps); err != nil {
				return err, false, math.MaxFloat64 * -1.0, -1
			}

			if !success && len(path) > 0 && !ec2.executor.Execute(target, "", path, "create", "all") {
				if err := ec2.DeleteInstance(target, 0); err != nil {
					return err, false, math.MaxFloat64 * -1.0, -1
				}
				continue
			}

			roles := make([]string, 0)
			roles = append(roles, ec2.GetYuma().GetConfigurations()[action])

			if ec2.executor.Execute(target, "", roles, "create", "all") {
				reward = -1.0
				success = true
				successor = state | action

				break
			} else {
				if err := ec2.DeleteInstance(target, 0); err != nil {
					return err, false, math.MaxFloat64 * -1.0, -1
				}
				reward = (float64(len(ec2.GetYuma().GetSubprocesses())) + 1.0) * -1.0
				success = false
				successor = state
			}
		}
	case 1:
		preExecute := true
		for _, preA := range path {
			if action == ec2.GetYuma().GetSubprocesses()[preA] {
				preExecute = false
			}
		}

		if len(path) > 0 && preExecute {
			instances := ec2.GetInstances(target)
			for l := 0; l < len(path); l++ {
				_, ok := instances[ec2.GetYuma().GetSubprocesses()[path[l]]]
				if ok {
					for len(instances[ec2.GetYuma().GetSubprocesses()[path[l]]]) < ec2.GetYuma().GetQuantities()[ec2.GetYuma().GetSubprocesses()[path[l]]] {
						if err, _, _, _ := ec2.TakeAction(target, state, ec2.GetYuma().GetSubprocesses()[path[l]], path, true); err != nil {
							return err, false, math.MaxFloat64 * -1.0, -1
						}
						instances = ec2.GetInstances(target)
					}
				} else {
					for len(instances[ec2.GetYuma().GetSubprocesses()[path[l]]]) < ec2.GetYuma().GetQuantities()[ec2.GetYuma().GetSubprocesses()[path[l]]] {
						if err, _, _, _ := ec2.TakeAction(target, state, ec2.GetYuma().GetSubprocesses()[path[l]], path, true); err != nil {
							return err, false, math.MaxFloat64 * -1.0, -1
						}
						instances = ec2.GetInstances(target)
					}
				}
			}
		}

		for i := 0; i < ec2.GetSigma(); i++ {
			err, inputs := ec2.GetExecutor().DetermineInputs(ec2.GetYuma().GetConfigurations()[action])
			execute := true
			if err == nil {
				for _, input := range inputs {
					if strings.Contains(input, "_") {
						input = strings.Split(input, "_")[0]
					}
					if _, ok := ec2.GetYuma().GetSubprocesses()[input]; ok {
						executed := false
						for _, a := range path {
							if input == a {
								executed = true
							}
						}
						if !executed {
							execute = false
							break
						}
					}
				}
			}

			if execute {
				err, inputs := ec2.GetExecutor().DetermineInputs(ec2.GetYuma().GetConfigurations()[action])
				if err == nil {
					for _, input := range inputs {
						index := -1
						if strings.Contains(input, "_") {
							index, _ = strconv.Atoi(strings.Split(input, "_")[1])
							input = strings.Split(input, "_")[0]
						}
						if _, ok := ec2.GetYuma().GetSubprocesses()[input]; ok {
							instances := ec2.GetInstances(target)[ec2.GetYuma().GetSubprocesses()[input]]
							if err := helper.ReplaceStringInFile(ec2.GetExecutor().GetRepository() + ec2.GetYuma().GetConfigurations()[action] + "/defaults/main.yml", "{{ " + input + "_" + strconv.Itoa(index) + " }}", instances[index - 1].GetPrivateIP()); err != nil {
								return err, false, math.MaxFloat64 * -1.0, -1
							}
						}
					}
				}

				for j := 0; j < ec2.GetYuma().GetQuantities()[action]; j++ {
					if err := ec2.CreateInstance(target, action, ec2.GetOmega()); err != nil {
						return err, false, math.MaxFloat64 * -1.0, -1
					}
				}
				instances := ec2.GetInstances(target)
				tmp := instances[action]
				publicIps := make([]string, 0)
				for k := len(tmp) - 1; k > (len(tmp) - 1) - ec2.GetYuma().GetQuantities()[action]; k-- {
					publicIps = append(publicIps, tmp[k].GetPublicIP())
				}
				if err := ec2.executor.CreateEnvironmentDescription(target, publicIps); err != nil {
					fmt.Println(err)
					return err, false, math.MaxFloat64 * -1.0, -1
				}

				roles := make([]string, 0)
				roles = append(roles, ec2.GetYuma().GetConfigurations()[action])

				if ec2.executor.Execute(target, "", roles, "create", "all") {
					reward = -1.0
					success = true
					successor = state | action
				} else {
					for i := 0; i < ec2.GetYuma().GetQuantities()[action]; i++ {
						if err := ec2.DeleteInstance(target, action); err != nil {
							return err, false, math.MaxFloat64 * -1.0, -1
						}
					}
					reward = (float64(len(ec2.GetYuma().GetSubprocesses())) + 1.0) * -1.0
					success = false
					successor = state
				}

				for _, input := range inputs {
					index := -1
					if strings.Contains(input, "_") {
						index, _ = strconv.Atoi(strings.Split(input, "_")[1])
						input = strings.Split(input, "_")[0]
					}
					if _, ok := ec2.GetYuma().GetSubprocesses()[input]; ok {
						instances := ec2.GetInstances(target)[ec2.GetYuma().GetSubprocesses()[input]]
						if err := helper.ReplaceStringInFile(ec2.GetExecutor().GetRepository() + ec2.GetYuma().GetConfigurations()[action] + "/defaults/main.yml", instances[index - 1].GetPrivateIP(), "{{ " + input + "_" + strconv.Itoa(index) + " }}"); err != nil {
							return err, false, math.MaxFloat64 * -1.0, -1
						}
					}
				}

				if success {
					break
				}
			} else {
				reward = (float64(len(ec2.GetYuma().GetSubprocesses())) + 1.0) * -1.0
				success = false
				successor = state

				break
			}
		}
	}

	return nil, success, reward, successor
}