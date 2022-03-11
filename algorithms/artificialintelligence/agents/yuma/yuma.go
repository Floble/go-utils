package yuma

import (
	"io/ioutil"
	"math"
	"sync"
	"gonum.org/v1/gonum/mat"
)

type Yuma struct {
	subprocesses map[string]int
	configurations map[int]string
	environment Environment
	model *mat.Dense
	rationalThinking RationalThinking
}

func NewYuma() *Yuma {
	yuma := new(Yuma)
	yuma.subprocesses = make(map[string]int, 0)
	yuma.configurations = make(map[int]string, 0)

	return yuma
}

func (yuma *Yuma) GetSubprocesses() map[string]int {
	return yuma.subprocesses
}

func (yuma *Yuma) GetConfigurations() map[int]string {
	return yuma.configurations
}

func (yuma *Yuma) GetEnvironment() Environment {
	return yuma.environment
}

func (yuma *Yuma) GetModel() *mat.Dense {
	return yuma.model
}

func (yuma *Yuma) SetModel(model *mat.Dense) {
	yuma.model = model
}

func (yuma *Yuma) GetRationalThinking() RationalThinking {
	return yuma.rationalThinking
}

func (yuma *Yuma) identifySubprocesses(path string) error {
	subprocesses, err := ioutil.ReadDir(path)
    if err != nil {
        return err
    }

	for i := 0; i < len(subprocesses); i++ {
		config := yuma.GetStartState()
		config |= int(math.Exp2(float64(i)))
		yuma.GetSubprocesses()[subprocesses[i].Name()] = config
		yuma.configurations[config] = subprocesses[i].Name()
	}
	
	return nil
}

func (yuma *Yuma) GetStartState() int {
	return 0
}

func (yuma *Yuma) SetEnvironment(environment Environment) error {
	yuma.environment = environment

	if err := yuma.identifySubprocesses(yuma.environment.GetExecutor().GetRepository()); err != nil {
		return err
	}
	yuma.model = mat.NewDense(int(math.Exp2(float64(len(yuma.subprocesses)))), int(math.Exp2(float64(len(yuma.subprocesses) - 1)) + 1), nil)

	return nil
}

func (yuma *Yuma) SetRationalThinking(rationalThinking RationalThinking) {
	yuma.rationalThinking = rationalThinking
}

func (yuma *Yuma) IsTerminal(state int) bool {
	if state == int(math.Exp2(float64(len(yuma.GetSubprocesses())))) - 1 {
		return true
	} else {
		return false
	}
}

func (yuma *Yuma) Actions(state int) []int {
	actions := make([]int, 0)

	for _, action := range yuma.GetSubprocesses() {
		if state & action != 0 {
			continue
		}

		actions = append(actions, action)
	}

	return actions
}

func (yuma *Yuma) LearnDependencies() <-chan error {
	var wg sync.WaitGroup

	//target := int(math.Exp2(float64(len(yuma.GetSubprocesses())))) - 1
	
	errs := make(chan error, 1)
	for i := 0; i < len(yuma.GetSubprocesses()); i++ {
		wg.Add(1)
		i := i
		go func() {
			defer wg.Done()
            if err := yuma.GetRationalThinking().Learn(i); err != nil {
				errs <- err
			}
        }()
	}

	wg.Wait()

	return errs
}

func (yuma *Yuma) DetermineMinimalExecutionOrder() error {
	target := 2

	mEO := yuma.GetRationalThinking().Solve(target)
	yuma.GetEnvironment().GetExecutor().SetPlaybook("playbook.yml")
	if err := yuma.GetEnvironment().GetExecutor().CreateExecutionOrder("", mEO, "install"); err != nil {
		return err
	}

	return nil
}