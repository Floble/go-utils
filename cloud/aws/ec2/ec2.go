package ec2

import (
	"fmt"
	"go-utils/algorithms/artificialintelligence/agents/yuma"
	"math"
	"sync"
	"time"
	"os"
	"strconv"
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

func (ec2 *EC2) GetInstance(target int) yuma.Instance {
	if instance, ok := ec2.instances.Load(target); ok {
		if instance == nil {
			return nil
		} else {
			return instance.(yuma.Instance)
		}
	} else {
		return nil
	}
}

func (ec2 *EC2) SetInstance(target int, instance yuma.Instance) {
	ec2.instances.Store(target, instance)
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

func (ec2 *EC2) CreateInstance(target int, waitingTime int) error {
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

	if err := ec2.executor.CreateEnvironmentDescription(target, instance.GetPublicIP()); err != nil {
		fmt.Println("EC2 ERROR: CREATE ENVIRONMENT FROM DESCRIPTION")
		return err
	}

	time.Sleep(time.Duration(waitingTime) * time.Second)
	
	added := false
	for added == false {
		if err := instance.AddToKnownHosts(); err != nil {
			fmt.Println("EC2 ERROR: ADD TO KNOWN HOSTS")
			//return err
			fmt.Println(err)
			time.Sleep(time.Duration(waitingTime) * time.Second)
		} else {
			added = true
		}
	}

	ec2.SetInstance(target, instance)

	return nil
}

func (ec2 *EC2) DeleteInstance(target int) error {
	deleted := false
	for deleted == false {
		if err := ec2.GetInstance(target).Delete(); err != nil {
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

	ec2.SetInstance(target, nil)

	return nil
}

func (ec2 *EC2) TakeAction(target int, state int, action int, path []string, success bool) (error, bool, float64, int) {
	reward := math.MaxFloat64 * -1.0
	successor := -1

	for i := 0; i < ec2.GetSigma(); i++ {
		if ec2.GetInstance(target) == nil {
			if err := ec2.CreateInstance(target, ec2.GetOmega()); err != nil {
				return err, false, math.MaxFloat64 * -1.0, -1
			}
		}

		if !success && len(path) > 0 && !ec2.executor.Execute(target, "", path, "create", "all") {
			if err := ec2.DeleteInstance(target); err != nil {
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
			if err := ec2.DeleteInstance(target); err != nil {
				return err, false, math.MaxFloat64 * -1.0, -1
			}
			reward = (float64(len(ec2.GetYuma().GetSubprocesses())) + 1.0) * -1.0
			success = false
			successor = state
		}
	}

	return nil, success, reward, successor
}