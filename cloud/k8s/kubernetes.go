package k8s

import (
	"go-utils/algorithms/artificialintelligence/agents/yuma"
	"go-utils/helper"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"
)

type Kubernetes struct {
	yuma *yuma.Yuma
	executor yuma.Executor
	//instance yuma.Instance
	instances sync.Map
	omega int
	sigma int
}

func NewKubernetes(yuma *yuma.Yuma, executor yuma.Executor, omega int, sigma int) *Kubernetes {
	ec2 := new(Kubernetes)
	ec2.yuma = yuma
	ec2.executor = executor
	ec2.omega = omega
	ec2.sigma = sigma

	return ec2
}

func (k8s *Kubernetes) GetYuma() *yuma.Yuma {
	return k8s.yuma
}

func (k8s *Kubernetes) GetExecutor() yuma.Executor {
	return k8s.executor
}

func (k8s *Kubernetes) GetInstances(target int) map[int][]yuma.Instance {
	if instances, ok := k8s.instances.Load(target); ok {
		if instances == nil {
			return nil
		} else {
			return instances.(map[int][]yuma.Instance)
		}
	} else {
		return nil
	}
}

func (k8s *Kubernetes) SetInstances(target int, instances map[int][]yuma.Instance) {
	k8s.instances.Store(target, instances)
}

func (k8s *Kubernetes) GetOmega() int {
	return k8s.omega
}

func (k8s *Kubernetes) GetSigma() int {
	return k8s.sigma
}

func (k8s *Kubernetes) Initialize() error {
	return nil
}

func (k8s *Kubernetes) CleanUp() error {
	return nil
}

func (k8s *Kubernetes) CleanResults() error {
	for i := 0; i < len(k8s.GetYuma().GetSubprocesses()); i++ {
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

func (k8s *Kubernetes) CreateInstance(target int, action int, waitingTime int) error {	
	instance := NewPod()
	err, inputs := k8s.GetExecutor().DetermineInputs(k8s.GetYuma().GetConfigurations()[action])
	if err == nil {
		for i, input := range inputs {
			if input == "port" {
				_, values := k8s.GetExecutor().DetermineValues(k8s.GetYuma().GetConfigurations()[action])
				port, _ := strconv.Atoi(strings.TrimSpace(values[i]))
				instance.SetPort(port)
			}
		}
	} else {
		return err
	}

	var tmp []yuma.Instance
	instances := k8s.GetInstances(target)
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
	k8s.SetInstances(target, instances)

	return nil
}

func (k8s *Kubernetes) DeleteInstance(target int, action int) error {
	instances := k8s.GetInstances(target)
	instances[action] = instances[action][:len(instances[action]) - 1]
	k8s.SetInstances(target, instances)

	return nil
}

func (k8s *Kubernetes) DeleteAllInstances(target int) error {
	instances := k8s.GetInstances(target)
	for action := range instances {
		for _, instance := range instances[action] {
			if err := instance.Delete(); err != nil {
				return err
			}
			roles := make([]string, 0)
			roles = append(roles, k8s.GetYuma().GetConfigurations()[action])
			k8s.executor.Execute(target, "", roles, "remove", "localhost")
		}
		instances[action] = make([]yuma.Instance, 0)
	}

	return nil
}

func (k8s *Kubernetes) TakeAction(target int, state int, action int, path []string, success bool) (error, bool, float64, int) {
	reward := math.MaxFloat64 * -1.0
	successor := -1

	switch k8s.GetYuma().GetMode() {
	case 1:
		preExecute := true
		for _, preA := range path {
			if action == k8s.GetYuma().GetSubprocesses()[preA] {
				preExecute = false
			}
		}

		if len(path) > 0 && preExecute {
			instances := k8s.GetInstances(target)
			for l := 0; l < len(path); l++ {
				_, ok := instances[k8s.GetYuma().GetSubprocesses()[path[l]]]
				if ok {
					for len(instances[k8s.GetYuma().GetSubprocesses()[path[l]]]) < k8s.GetYuma().GetQuantities()[k8s.GetYuma().GetSubprocesses()[path[l]]] {
						if err, _, _, _ := k8s.TakeAction(target, state, k8s.GetYuma().GetSubprocesses()[path[l]], path, true); err != nil {
							return err, false, math.MaxFloat64 * -1.0, -1
						}
						instances = k8s.GetInstances(target)
					}
				} else {
					for len(instances[k8s.GetYuma().GetSubprocesses()[path[l]]]) < k8s.GetYuma().GetQuantities()[k8s.GetYuma().GetSubprocesses()[path[l]]] {
						if err, _, _, _ := k8s.TakeAction(target, state, k8s.GetYuma().GetSubprocesses()[path[l]], path, true); err != nil {
							return err, false, math.MaxFloat64 * -1.0, -1
						}
						instances = k8s.GetInstances(target)
					}
				}
			}
		}

		for i := 0; i < k8s.GetSigma(); i++ {
			err, inputs := k8s.GetExecutor().DetermineInputs(k8s.GetYuma().GetConfigurations()[action])
			execute := true
			if err == nil {
				for _, input := range inputs {
					_, ok := k8s.GetYuma().GetSubprocesses()[input]
					if ok {
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
				err, inputs := k8s.GetExecutor().DetermineInputs(k8s.GetYuma().GetConfigurations()[action])
				if err == nil {
					for _, input := range inputs {
						port := false
						if strings.Contains(input, "_") {
							if strings.Split(input, "_")[1] == "port" {
								port = true
							}
							input = strings.Split(input, "_")[0]
						}
						_, ok := k8s.GetYuma().GetSubprocesses()[input]
						if ok && port {
							instances := k8s.GetInstances(target)[k8s.GetYuma().GetSubprocesses()[input]]
							if err := helper.ReplaceStringInFile(k8s.GetExecutor().GetRepository() + k8s.GetYuma().GetConfigurations()[action] + "/defaults/main.yml", "{{ " + input + "_port }}", strconv.Itoa(instances[0].GetPort())); err != nil {
								return err, false, math.MaxFloat64 * -1.0, -1
							}
						} else if ok && !port {
							if err := helper.ReplaceStringInFile(k8s.GetExecutor().GetRepository() + k8s.GetYuma().GetConfigurations()[action] + "/defaults/main.yml", "{{ " + input + " }}", input); err != nil {
								return err, false, math.MaxFloat64 * -1.0, -1
							}
						}
					}
				}

				for j := 0; j < k8s.GetYuma().GetQuantities()[action]; j++ {
					if err := k8s.CreateInstance(target, action, k8s.GetOmega()); err != nil {
						return err, false, math.MaxFloat64 * -1.0, -1
					}
				}

				roles := make([]string, 0)
				roles = append(roles, k8s.GetYuma().GetConfigurations()[action])

				if k8s.executor.Execute(target, "", roles, "create", "localhost") {
					//time.Sleep(1 * time.Minute)

					reward = -1.0
					success = true
					successor = state | action
				} else {
					for i := 0; i < k8s.GetYuma().GetQuantities()[action]; i++ {
						if err := k8s.DeleteInstance(target, action); err != nil {
							return err, false, math.MaxFloat64 * -1.0, -1
						}
					}

					k8s.executor.Execute(target, "", roles, "remove", "localhost")
					//time.Sleep(1 * time.Minute)

					reward = (float64(len(k8s.GetYuma().GetSubprocesses())) + 1.0) * -1.0
					success = false
					successor = state
				}
				if success {
					break
				}
			} else {
				reward = (float64(len(k8s.GetYuma().GetSubprocesses())) + 1.0) * -1.0
				success = false
				successor = state

				break
			}
		}
	}

	return nil, success, reward, successor
}