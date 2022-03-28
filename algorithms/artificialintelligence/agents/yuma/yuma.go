package yuma

import (
	"fmt"
	"io/ioutil"
	"math"
	"sync"
	"strconv"
	"gonum.org/v1/gonum/mat"
)

type Yuma struct {
	subprocesses map[string]int
	configurations map[int]string
	environment Environment
	//models map[int]*mat.Dense
	models sync.Map
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

func (yuma *Yuma) GetModel(target int) *mat.Dense {
	if model, ok := yuma.models.Load(target); ok {
		return model.(*mat.Dense)
	} else {
		return nil
	}
}

func (yuma *Yuma) SetModel(target int, model *mat.Dense) {
	yuma.models.Store(target, model)
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
	
	for i := 0; i < len(yuma.GetSubprocesses()); i++ {
		yuma.models.Store(int(math.Exp2(float64(i))), mat.NewDense(int(math.Exp2(float64(len(yuma.subprocesses)))), int(math.Exp2(float64(len(yuma.subprocesses) - 1)) + 1), nil))
	}

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
	wg.Add(len(yuma.GetSubprocesses()))
	errs := make(chan error, 1)

	if err := yuma.GetEnvironment().Initialize(); err != nil {
		fmt.Println("DIRECTORY COPY ERROR: " + err.Error())
		errs <- err
	}

	for i := 0; i < len(yuma.GetSubprocesses()); i++ {
		target := math.Exp2(float64(i))
		go func() {
			defer wg.Done()
			behaviorPolicy := NewEpsilonGreedyPolicy(0.1)
			targetPolicy := NewGreedyPolicy()
			rationalThinking := NewTreeBackup(yuma, behaviorPolicy, targetPolicy, 50000, 0.5, 1, 0)
			rationalThinking.SetN((len(yuma.GetSubprocesses()) / 2) + 1)
			yuma.SetRationalThinking(rationalThinking)
			behaviorPolicy.SetRationalThinking(rationalThinking)
			targetPolicy.SetRationalThinking(rationalThinking)
            yuma.GetRationalThinking().Learn(int(target))
        }()
	}
	wg.Wait()

	if err := yuma.GetEnvironment().CleanUp(); err != nil {
		fmt.Println("DIRECTORY DELETE ERROR: " + err.Error())
		errs <- err
	}

	return errs
}

func (yuma *Yuma) DetermineMinimalExecutionOrder() error {
	for i := 0; i < len(yuma.GetSubprocesses()); i++ {
		target := math.Exp2(float64(i))
		pathExecutionOrder := "playbook_" + strconv.Itoa(int(math.Exp2(float64(i)))) + ".yml"
		
		mEO := yuma.GetRationalThinking().Solve(int(target))
		if err := yuma.GetEnvironment().GetExecutor().CreateExecutionOrder(pathExecutionOrder, "", mEO, "install"); err != nil {
			return err
		}
	}

	return nil
}