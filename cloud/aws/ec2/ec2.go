package ec2

import (
	"time"
	"math"
	"go-utils/algorithms/artificialintelligence/agents/yuma"
)

type EC2 struct {
	yuma *yuma.Yuma
	executor yuma.Executor
	instance yuma.Instance
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

func (ec2 *EC2) GetInstance() yuma.Instance {
	return ec2.instance
}

func (ec2 *EC2) SetInstance(instance yuma.Instance) {
	ec2.instance = instance
}

func (ec2 *EC2) GetOmega() int {
	return ec2.omega
}

func (ec2 *EC2) GetSigma() int {
	return ec2.sigma
}

func (ec2 *EC2) CreateInstance(waitingTime int) error {
	instance := NewEC2Instance()
	if err := instance.Create(); err != nil {
		return err
	}

	if err := ec2.executor.CreateEnvironmentDescription(instance.GetPublicIP()); err != nil {
		return err
	}

	time.Sleep(time.Duration(waitingTime) * time.Second)
	
	if err := instance.AddToKnownHosts(); err != nil {
		return err
	}

	ec2.SetInstance(instance)

	return nil
}

func (ec2 *EC2) DeleteInstance() error {
	if err := ec2.GetInstance().Delete(); err != nil {
		return err
	}
	if err := ec2.executor.RemoveEnvironmentDescription(); err != nil {
		return err
	}

	ec2.SetInstance(nil)

	return nil
}

func (ec2 *EC2) TakeAction(state int, action int, path []string, success bool) (error, bool, float64, int) {
	reward := math.MaxFloat64 * -1.0
	successor := -1

	for i := 0; i < ec2.GetSigma(); i++ {
		if ec2.GetInstance() == nil {
			if err := ec2.CreateInstance(ec2.GetOmega()); err != nil {
				return err, false, math.MaxFloat64 * -1.0, -1
			}
		}

		if !success && len(path) > 0 && !ec2.executor.Execute("", path, "install") {
			if err := ec2.DeleteInstance(); err != nil {
				return err, false, math.MaxFloat64 * -1.0, -1
			}
		}

		roles := make([]string, 0)
		roles = append(roles, ec2.GetYuma().GetConfigurations()[action])

		if ec2.executor.Execute("", roles, "install") {
			reward = -1.0
			success = true
			successor = state | action

			break
		} else {
			if err := ec2.DeleteInstance(); err != nil {
				return err, false, math.MaxFloat64 * -1.0, -1
			}
			reward = -10.0
			success = false
			successor = state
		}
	}

	return nil, success, reward, successor
}