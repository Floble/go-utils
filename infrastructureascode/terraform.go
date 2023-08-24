package infrastructureascode

import (
	"go-utils/algorithms/artificialintelligence/agents/yuma"
	//"go-utils/helper"
	"math"
	"os"
	//"fmt"
	"strconv"
	"strings"
)

type Terraform struct {
	yuma *yuma.Yuma
	executor yuma.Executor
	repository string
	omega int
	sigma int
}

func NewTerraform(yuma *yuma.Yuma, executor yuma.Executor, repository string, omega int, sigma int) *Terraform {
	terraform := new(Terraform)
	terraform.yuma = yuma
	terraform.executor = executor
	terraform.repository = repository
	terraform.omega = omega
	terraform.sigma = sigma

	return terraform
}

func (terraform *Terraform) GetYuma() *yuma.Yuma {
	return terraform.yuma
}

func (terraform *Terraform) GetExecutor() yuma.Executor {
	return terraform.executor
}

func (terraform *Terraform) GetInstances(target int) map[int][]yuma.Instance {
	return nil
}

func (terraform *Terraform) SetInstances(target int, instances map[int][]yuma.Instance) {
}

func (terraform *Terraform) GetRepository() string {
	return terraform.repository
}

func (terraform *Terraform) GetOmega() int {
	return terraform.omega
}

func (terraform *Terraform) GetSigma() int {
	return terraform.sigma
}

/* func (terraform *Terraform) Initialize() error {
	for i := 0; i < len(terraform.GetYuma().GetSubprocesses()); i++ {
		if err := helper.CopyDirectory(terraform.GetRepository(), terraform.GetRepository() + "_" + strconv.Itoa(int(math.Exp2(float64(i))))); err != nil {
			return err
		}

		if err := helper.CopyDirectory(terraform.GetExecutor().GetRepository(), strings.Split(terraform.GetExecutor().GetRepository(), "/")[0] + "_" + strconv.Itoa(int(math.Exp2(float64(i))))); err != nil {
			return err
		}
	}

	for k := 0; k < len(terraform.GetYuma().GetSubprocesses()); k++ {
		for i := 0; i < len(terraform.GetYuma().GetSubprocesses()); i++ {
			err, inputs := terraform.executor.DetermineInputs(terraform.GetYuma().GetConfigurations()[int(math.Exp2(float64(i)))])
			if err == nil {
				for j := 0; j < len(inputs); j++ {
					if err := helper.ReplaceStringInFile(strings.Split(terraform.GetExecutor().GetRepository(), "/")[0] + "_" + strconv.Itoa(int(math.Exp2(float64(k)))) + "/" + terraform.GetYuma().GetConfigurations()[int(math.Exp2(float64(i)))] + "/defaults/main.yml", inputs[j] + ":", inputs[j] + strconv.Itoa(int(math.Exp2(float64(k)))) + ":"); err != nil {
						fmt.Println(err)
						return err
					}
	
					if err := helper.ReplaceStringInFile(strings.Split(terraform.GetExecutor().GetRepository(), "/")[0] + "_" + strconv.Itoa(int(math.Exp2(float64(k)))) + "/" + terraform.GetYuma().GetConfigurations()[int(math.Exp2(float64(i)))] + "/tasks/main.yml", "{{ " + inputs[j] + " }}", "{{ " + inputs[j] + strconv.Itoa(int(math.Exp2(float64(k)))) + " }}"); err != nil {
						fmt.Println(err)
						return err
					}
				}
			}
	
			if err := helper.ReplaceStringInFile(strings.Split(terraform.GetExecutor().GetRepository(), "/")[0] + "_" + strconv.Itoa(int(math.Exp2(float64(k)))) + "/" + terraform.GetYuma().GetConfigurations()[int(math.Exp2(float64(i)))] + "/tasks/main.yml", "!module_path!", terraform.GetRepository() + "_" + strconv.Itoa(int(math.Exp2(float64(k))))); err != nil {
				fmt.Println(err)
				return err
			}

			err, outputs := terraform.executor.DetermineOutputs(terraform.GetYuma().GetConfigurations()[int(math.Exp2(float64(i)))])
			if err == nil && len(outputs) > 0 {
				for _, output := range outputs {
					if err := helper.ReplaceStringInFile(strings.Split(terraform.GetExecutor().GetRepository(), "/")[0] + "_" + strconv.Itoa(int(math.Exp2(float64(k)))) + "/" + terraform.GetYuma().GetConfigurations()[int(math.Exp2(float64(i)))] + "/tasks/main.yml", output + ":", output + strconv.Itoa(int(math.Exp2(float64(k)))) + ":"); err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
} */

func (terraform *Terraform) Initialize() error {
	return nil
}

func (terraform *Terraform) CleanUp() error {
	for i := 0; i < len(terraform.GetYuma().GetSubprocesses()); i++ {
		if err := os.RemoveAll(terraform.GetRepository() + "_" + strconv.Itoa(int(math.Exp2(float64(i))))); err != nil {
			return err
		}
	}

	for i := 0; i < len(terraform.GetYuma().GetSubprocesses()); i++ {
		if err := os.RemoveAll(strings.Split(terraform.GetExecutor().GetRepository(), "/")[0] + "_" + strconv.Itoa(int(math.Exp2(float64(i))))); err != nil {
			return err
		}
	}

	return nil
}

func (terraform *Terraform) CleanResults() error {
	for i := 0; i < len(terraform.GetYuma().GetSubprocesses()); i++ {
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

func (terraform *Terraform) CreateInstance(target int, action int, waitingTime int) error {
	return nil
}

func (terraform *Terraform) DeleteInstance(target int, action int) error {
	return nil
}

func (terraform *Terraform) DeleteAllInstances(target int) error {
	return nil
}

func (terraform *Terraform) TakeAction(target int, state int, action int, path []string, success bool) (error, bool, float64, int) {
	reward := math.MaxFloat64 * -1.0
	successor := -1

	for i := 0; i < terraform.GetSigma(); i++ {
		role := terraform.GetYuma().GetConfigurations()[action]
		var pathCopy []string
		pathCopy = append(pathCopy, path...)
		pathCopy = append(pathCopy, role)

		if terraform.executor.Execute(target, "", pathCopy, "create", "localhost") {
			terraform.executor.Execute(target, "", pathCopy, "remove", "localhost")
			reward = -1.0
			success = true
			successor = state | action

			break
		} else {
			terraform.executor.Execute(target, "", pathCopy, "remove", "localhost")
			reward = (float64(len(terraform.GetYuma().GetSubprocesses())) + 1.0) * -1.0
			success = false
			successor = state
		}
	}

	return nil, success, reward, successor
}