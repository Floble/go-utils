package molecule

import (
	"bufio"
	"fmt"
	"go-utils/algorithms/artificialintelligence/agents/yuma"
	"go-utils/helper"
	"math"
	"os"
	"os/exec"
	"strconv"
	"time"
	"sync"
)

type Molecule struct {
	yuma *yuma.Yuma
	executor yuma.Executor
	instances sync.Map
	omega int
	sigma int
}

func NewMolecule(yuma *yuma.Yuma, executor yuma.Executor, omega int, sigma int) *Molecule {
	molecule := new(Molecule)
	molecule.yuma = yuma
	molecule.executor = executor
	molecule.omega = omega
	molecule.sigma = sigma

	return molecule
}

func (molecule *Molecule) GetYuma() *yuma.Yuma {
	return molecule.yuma
}

func (molecule *Molecule) GetExecutor() yuma.Executor {
	return molecule.executor
}

func (molecule *Molecule) GetInstances(target int) map[int][]yuma.Instance {
	if instances, ok := molecule.instances.Load(target); ok {
		if instances == nil {
			return nil
		} else {
			return instances.(map[int][]yuma.Instance)
		}
	} else {
		return nil
	}
}

func (molecule *Molecule) SetInstances(target int, instances map[int][]yuma.Instance) {
	molecule.instances.Store(target, instances)
}

func (molecule *Molecule) GetOmega() int {
	return molecule.omega
}

func (molecule *Molecule) GetSigma() int {
	return molecule.sigma
}

func (molecule *Molecule) Initialize() error {
	for i := 0; i < len(molecule.GetYuma().GetSubprocesses()); i++ {
		if err := helper.CopyDirectory("molecule_example", "molecule_" + strconv.Itoa(int(math.Exp2(float64(i))))); err != nil {
			return err
		}
	}

	return nil
}

func (molecule *Molecule) CleanUp() error {
	for i := 0; i < len(molecule.GetYuma().GetSubprocesses()); i++ {
		if err := os.RemoveAll("molecule_" + strconv.Itoa(int(math.Exp2(float64(i))))); err != nil {
			return err
		}
	}

	return nil
}

func (molecule *Molecule) CleanResults() error {
	for i := 0; i < len(molecule.GetYuma().GetSubprocesses()); i++ {
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

func (molecule *Molecule) CreateInstance(target int, action int, waitingTime int) error {
	cmd := exec.Command("molecule", "create")
	cmd.Dir = "molecule_" + strconv.Itoa(target)
	if err := cmd.Start(); err != nil {
		return err
	}
  
	if err := cmd.Wait(); err != nil {
		return err
	}

	time.Sleep(time.Duration(waitingTime) * time.Second)

	instances := molecule.GetInstances(target)
	var tmp []yuma.Instance
	if tmp, ok := instances[action]; ok {
		tmp = append(tmp, NewContainer(""))
	} else {
		tmp := make([]yuma.Instance, 0)
		tmp = append(tmp, NewContainer(""))
	}
	instances[action] = tmp
	molecule.SetInstances(target, instances)

	return nil
}

func (molecule *Molecule) DeleteInstance(target int, action int) error {
	cmd := exec.Command("molecule", "destroy")
	cmd.Dir = "molecule_" + strconv.Itoa(target)
	err := cmd.Start()
	if err != nil {
		return err
	}
  
	err = cmd.Wait()
	if err != nil {
		return err
	}

	instances := molecule.GetInstances(target)
	tmp := instances[action]
	tmp = tmp[:len(tmp) - 1]
	instances[action] = tmp
	molecule.SetInstances(target, instances)

	return nil
}

func (molecule *Molecule) DeleteAllInstances(target int) error {
	return nil
}

func (molecule *Molecule) TakeAction(target int, state int, action int, path []string, success bool) (error, bool, float64, int) {
	reward := math.MaxFloat64 * -1.0
	successor := -1

	for i := 0; i < molecule.GetSigma(); i++ {
		/* if molecule.GetInstance() == nil {
			if err := molecule.CreateInstance(target, 0); err != nil {
				return err, false, math.MaxFloat64 * -1.0, -1
			}
		} */

		/* if !success && len(path) > 0 && !molecule.test(target, "../../../", path, "create") {
			if err := molecule.DeleteInstance(target); err != nil {
				return err, false, math.MaxFloat64 * -1.0, -1
			}
			continue
		} */

		/* roles := make([]string, 0)
		roles = append(roles, molecule.GetYuma().GetConfigurations()[action]) */

		path = append(path, molecule.GetYuma().GetConfigurations()[action])

		/* if molecule.test(target, "../../../", roles, "create") {
			reward = -1.0
			success = true
			successor = state | action

			break
		} else {
			if err := molecule.DeleteInstance(target); err != nil {
				return err, false, math.MaxFloat64 * -1.0, -1
			}
			reward = -10.0
			success = false
			successor = state
		} */

		if molecule.test(target, "../../../", path, "create") {
			reward = -1.0
			success = true
			successor = state | action
			break
		} else {
			reward = (float64(len(molecule.GetYuma().GetSubprocesses())) + 1.0) * -1.0
			success = false
			successor = state
		}
	}

	return nil, success, reward, successor
}

func (molecule *Molecule) test(target int, pathPrefix string, roles []string, lifecycle string) bool {
	file, err := os.OpenFile("logs_" + strconv.Itoa(target) + ".txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("MOLECULE - Open File Error")
		return false
	}
	defer file.Close()

	if err := molecule.executor.CreateExecutionOrder(0, "molecule_" + strconv.Itoa(target) + "/molecule/default/converge.yml", pathPrefix, roles, lifecycle, "all"); err != nil {
		fmt.Println("MOLECULE - Create Execution Order Error")
		return false
	}
	
	cmd := exec.Command("molecule", "test", "--parallel")
	cmd.Dir = "molecule_" + strconv.Itoa(target)
	out, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("MOLECULE - Command Creation Error")
		molecule.executor.DeleteExecutionOrder("molecule_" + strconv.Itoa(target) + "/molecule/default/converge.yml")
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
		fmt.Println("MOLECULE - Command Start Error")
		fmt.Println(err)
		molecule.executor.DeleteExecutionOrder("molecule_" + strconv.Itoa(target) + "/molecule/default/converge.yml")
		return false
	}
	if err := cmd.Wait(); err != nil {
		molecule.executor.DeleteExecutionOrder("molecule_" + strconv.Itoa(target) + "/molecule/default/converge.yml")
		return false
	}

	molecule.executor.DeleteExecutionOrder("molecule_" + strconv.Itoa(target) + "/molecule/default/converge.yml")

	return true
}