package molecule

import (
	"math"
	"os"
	"os/exec"
	"bufio"
	"time"
	"go-utils/algorithms/artificialintelligence/agents/yuma"
)

type Molecule struct {
	yuma *yuma.Yuma
	executor yuma.Executor
	instance yuma.Instance
	omega int
	sigma int
	path string
}

func NewMolecule(yuma *yuma.Yuma, executor yuma.Executor, omega int, sigma int, path string) *Molecule {
	molecule := new(Molecule)
	molecule.yuma = yuma
	molecule.executor = executor
	molecule.omega = omega
	molecule.sigma = sigma
	molecule.path = path

	return molecule
}

func (molecule *Molecule) GetYuma() *yuma.Yuma {
	return molecule.yuma
}

func (molecule *Molecule) GetExecutor() yuma.Executor {
	return molecule.executor
}

func (molecule *Molecule) GetInstance() yuma.Instance {
	return molecule.instance
}

func (molecule *Molecule) SetInstance(instance yuma.Instance) {
	molecule.instance = instance
}

func (molecule *Molecule) GetOmega() int {
	return molecule.omega
}

func (molecule *Molecule) GetSigma() int {
	return molecule.sigma
}

func (molecule *Molecule) GetPath() string {
	return molecule.path
}

func (molecule *Molecule) CreateInstance(waitingTime int) error {
	cmd := exec.Command("molecule", "create")
	cmd.Dir = molecule.path
	if err := cmd.Start(); err != nil {
		return err
	}
  
	if err := cmd.Wait(); err != nil {
		return err
	}

	time.Sleep(time.Duration(waitingTime) * time.Second)

	molecule.SetInstance(NewContainer(""))

	return nil
}

func (molecule *Molecule) DeleteInstance() error {
	cmd := exec.Command("molecule", "destroy")
	cmd.Dir = molecule.path
	err := cmd.Start()
	if err != nil {
		return err
	}
  
	err = cmd.Wait()
	if err != nil {
		return err
	}

	molecule.SetInstance(nil)

	return nil
}

func (molecule *Molecule) TakeAction(state int, action int, path []string, success bool) (error, bool, float64, int) {
	reward := math.MaxFloat64 * -1.0
	successor := -1

	for i := 0; i < molecule.GetSigma(); i++ {
		if molecule.GetInstance() == nil {
			if err := molecule.CreateInstance(0); err != nil {
				return err, false, math.MaxFloat64 * -1.0, -1
			}
		}

		if !success && len(path) > 0 && !molecule.converge("../../../", path, "install") {
			if err := molecule.DeleteInstance(); err != nil {
				return err, false, math.MaxFloat64 * -1.0, -1
			}
			continue
		}

		roles := make([]string, 0)
		roles = append(roles, molecule.GetYuma().GetConfigurations()[action])

		if molecule.converge("../../../", roles, "install") {
			reward = -1.0
			success = true
			successor = state | action

			break
		} else {
			if err := molecule.DeleteInstance(); err != nil {
				return err, false, math.MaxFloat64 * -1.0, -1
			}
			reward = -10.0
			success = false
			successor = state
		}
	}

	return nil, success, reward, successor
}

func (molecule *Molecule) converge(pathPrefix string, roles []string, lifecycle string) bool {
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return false
	}
	defer file.Close()

	if err := molecule.executor.CreateExecutionOrder(pathPrefix, roles, lifecycle); err != nil {
		return false
	}
	
	cmd := exec.Command("molecule", "converge")
	cmd.Dir = molecule.path
	out, err := cmd.StdoutPipe()
	if err != nil {
		molecule.executor.DeleteExecutionOrder()
		return false
	}
  
	scanner := bufio.NewScanner(out)
	go func() {
		for scanner.Scan() {
			file.WriteString(scanner.Text())
		}
	}()
	file.WriteString("\n\n")

	if err := cmd.Start(); err != nil {
		molecule.executor.DeleteExecutionOrder()
		return false
	}
	if err := cmd.Wait(); err != nil {
		molecule.executor.DeleteExecutionOrder()
		return false
	}

	molecule.executor.DeleteExecutionOrder()

	return true
}