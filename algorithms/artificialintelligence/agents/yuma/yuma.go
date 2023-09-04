package yuma

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
	"gonum.org/v1/gonum/mat"
)

type Yuma struct {
	mode int
	subprocesses map[string]int
	configurations map[int]string
	quantities map[int]int
	solutions map[string][]string
	environment Environment
	//models map[int]*mat.Dense
	models sync.Map
	memory sync.Map
	timestamps sync.Map
	history sync.Map
	updates sync.Map
	rationalThinking RationalThinking
}

func NewYuma() *Yuma {
	yuma := new(Yuma)
	yuma.mode = 0

	yuma.quantities = make(map[int]int, 0)
	yuma.subprocesses = make(map[string]int, 0)
	yuma.configurations = make(map[int]string, 0)
	yuma.solutions = make(map[string][]string, 0)

	return yuma
}

func (yuma *Yuma) GetMode() int {
	return yuma.mode
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

func (yuma *Yuma) GetTimestamps(target int) *mat.Dense {
	if timestamps, ok := yuma.timestamps.Load(target); ok {
		return timestamps.(*mat.Dense)
	} else {
		return nil
	}
}

func (yuma *Yuma) SetTimestamps(target int, timestamps *mat.Dense) {
	yuma.timestamps.Store(target, timestamps)
}

func (yuma *Yuma) GetMemory(target int) *mat.Dense {
	if memory, ok := yuma.memory.Load(target); ok {
		return memory.(*mat.Dense)
	} else {
		return nil
	}
}

func (yuma *Yuma) SetMemory(target int, memory *mat.Dense) {
	yuma.memory.Store(target, memory)
}

func (yuma *Yuma) GetHistory(target int) *mat.Dense {
	if history, ok := yuma.history.Load(target); ok {
		return history.(*mat.Dense)
	} else {
		return nil
	}
}

func (yuma *Yuma) SetHistory(target int, history *mat.Dense) {
	yuma.history.Store(target, history)
}

func (yuma *Yuma) GetUpdates(target int) *mat.Dense {
	if updates, ok := yuma.updates.Load(target); ok {
		return updates.(*mat.Dense)
	} else {
		return nil
	}
}

func (yuma *Yuma) SetUpdates(target int, updates *mat.Dense) {
	yuma.updates.Store(target, updates)
}

func (yuma *Yuma) GetSolutions() map[string][]string {
	return yuma.solutions
}

func (yuma *Yuma) GetSolution(target string) []string {
	if _, ok := yuma.solutions[target]; ok {
		return yuma.solutions[target]
	} else {
		return nil
	}
}

func (yuma *Yuma) SetSolution(target string, solution []string) {
	yuma.solutions[target] = solution
}

func (yuma *Yuma) GetQuantities() map[int]int {
	return yuma.quantities
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

	for i := 0; i < len(yuma.GetSubprocesses()); i++ {
		yuma.GetQuantities()[int(math.Exp2(float64(i)))] = 1
	}
	if _, err := os.Stat("roles.cfg"); err == nil {
		readFile, err := os.Open("roles.cfg")
		if err != nil {
			return err
		} else {
			fileScanner := bufio.NewScanner(readFile)
    		fileScanner.Split(bufio.ScanLines)
			fileLines := make([]string, 0)
			for fileScanner.Scan() {
				fileLines = append(fileLines, fileScanner.Text())
			}
			if err := yuma.setDistinctMode(fileLines); err != nil {
				return err
			}
		}
		readFile.Close()
	}
	
	return nil
}

func (yuma *Yuma) setDistinctMode(fileLines []string) error {
	yuma.mode = 1

	for _, line := range fileLines {
		splits := strings.Split(line, ":")
		quantity, err := strconv.Atoi(strings.TrimSpace(splits[1]))
		if err != nil {
			return err
		}

		yuma.GetQuantities()[yuma.GetSubprocesses()[splits[0]]] = quantity
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
		yuma.updates.Store(int(math.Exp2(float64(i))), mat.NewDense(int(math.Exp2(float64(len(yuma.subprocesses)))), int(math.Exp2(float64(len(yuma.subprocesses) - 1)) + 1), nil))
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

func (yuma *Yuma) LearnDependenciesAsynchronously() <-chan error {
	var wg sync.WaitGroup
	errs := make(chan error, 1)
	
	if err := yuma.GetEnvironment().Initialize(); err != nil {
		fmt.Println("DIRECTORY COPY ERROR: " + err.Error())
		errs <- err
	}

	for sim := 0; sim < 21; sim++ {
		fmt.Printf("SIMULATION: %d\n", sim)

		wg.Add(len(yuma.GetSubprocesses()))
		for i := 0; i < len(yuma.GetSubprocesses()); i++ {
			target := int(math.Exp2(float64(i)))
			model := yuma.GetModel(target)
			timestamps := yuma.GetTimestamps(target)
			updates := yuma.GetUpdates(target)
			history := yuma.GetHistory(target)
			memory := yuma.GetMemory(target)

			if sim == 0 {
				// Initialize Q(s, a) arbitrarily, for all s, a
				for j := 0; j < int(math.Exp2(float64(len(yuma.GetSubprocesses())))); j++ {
					for k := 0; k < int(math.Exp2(float64(len(yuma.GetSubprocesses()) - 1)) + 1); k++ {
						if j & target != 0 || k == 0 {
							model.Set(j, k, 0.0)
						} else {
							model.Set(j, k, (float64(len(yuma.GetSubprocesses())) + 1.0) * -1.0)
						}
					}
				}
				updates = mat.NewDense(int(math.Exp2(float64(len(yuma.GetSubprocesses())))), int(math.Exp2(float64(len(yuma.GetSubprocesses()) - 1)) + 1), nil)
				history = mat.NewDense(int(math.Exp2(float64(len(yuma.GetSubprocesses())))), int(math.Exp2(float64(len(yuma.GetSubprocesses()) - 1)) + 1), nil)
				memory = mat.NewDense(int(math.Exp2(float64(len(yuma.GetSubprocesses())))), int(math.Exp2(float64(len(yuma.GetSubprocesses()) - 1)) + 1), nil)
				timestamps = mat.NewDense(int(math.Exp2(float64(len(yuma.GetSubprocesses())))), int(math.Exp2(float64(len(yuma.GetSubprocesses()) - 1)) + 1), nil)
			}

			go func() {
				defer wg.Done()

				behaviorPolicy := NewNaraEpsilonGreedyPolicy(0.3)
				targetPolicy := NewNaraGreedyPolicy()
				treeBackup := NewTreeBackup(yuma, behaviorPolicy, targetPolicy, 1, 0.5, 1, 24 * time.Hour, 0)
				treeBackup.SetN((len(yuma.GetSubprocesses())))

				treePolicy := NewNaraTreePolicy(0.3)
				expansionPolicy := NewNaraExpansionPolicy()
				selectionPolicy := NewNaraGreedyPolicy()
				rationalThinking := NewNara(yuma, treeBackup, treePolicy, expansionPolicy, selectionPolicy, 1 * time.Minute, 1, 0.8, 0.5, 1, 24 * time.Hour)
				//rationalThinking := NewNara(yuma, treeBackup, treePolicy, selectionPolicy, memory, nil, 5 * time.Minute, 1, 1, 0.5)
				yuma.SetRationalThinking(rationalThinking)

				behaviorPolicy.SetRationalThinking(treeBackup)
				targetPolicy.SetRationalThinking(treeBackup)
				treePolicy.SetRationalThinking(rationalThinking)
				expansionPolicy.SetRationalThinking(rationalThinking)
				selectionPolicy.SetRationalThinking(rationalThinking)

				//yuma.GetRationalThinking().Learn(int(target))
				eo := yuma.GetRationalThinking().Solve(int(target), model, updates, history, memory, timestamps)
				if err := yuma.ExportExecutionOrder(eo, float64(target)); err != nil {
					fmt.Println("EXPORT EXECUTION ORDER ERROR: " + err.Error())
					errs <- err
				}
			}()
		}
		wg.Wait()
	}

	if err := yuma.GetEnvironment().CleanUp(); err != nil {
		fmt.Println("DIRECTORY DELETE ERROR: " + err.Error())
		errs <- err
	}

	return errs
}

func (yuma *Yuma) LearnDependenciesSequentielly() error {
	if err := yuma.GetEnvironment().Initialize(); err != nil {
		fmt.Println("DIRECTORY COPY ERROR: " + err.Error())
		return err
	}

	history := mat.NewDense(int(math.Exp2(float64(len(yuma.GetSubprocesses())))), int(math.Exp2(float64(len(yuma.GetSubprocesses()) - 1)) + 1), nil)
	memory := mat.NewDense(int(math.Exp2(float64(len(yuma.GetSubprocesses())))), int(math.Exp2(float64(len(yuma.GetSubprocesses()) - 1)) + 1), nil)
	timestamps := mat.NewDense(int(math.Exp2(float64(len(yuma.GetSubprocesses())))), int(math.Exp2(float64(len(yuma.GetSubprocesses()) - 1)) + 1), nil)

	sim := 0
	i := 0
	for {
		if (sim > 0) && (sim % 5 == 0) {
			if err := yuma.GetEnvironment().CleanResults(); err != nil {
				fmt.Println("RESULTS DELETE ERROR: " + err.Error())
				return err
			}
			sim = 1
		}
		fmt.Printf("SIMULATION: %d\n", i)

		for i := 0; i < len(yuma.GetSubprocesses()); i++ {
			target := int(math.Exp2(float64(i)))
			model := yuma.GetModel(target)
			//timestamps := yuma.GetTimestamps(target)
			updates := yuma.GetUpdates(target)
			//history := yuma.GetHistory(target)
			//memory := yuma.GetMemory(target)

			if sim == 0 {
				// Initialize Q(s, a) arbitrarily, for all s, a
				for j := 0; j < int(math.Exp2(float64(len(yuma.GetSubprocesses())))); j++ {
					for k := 0; k < int(math.Exp2(float64(len(yuma.GetSubprocesses()) - 1)) + 1); k++ {
						if j & target != 0 || k == 0 {
							model.Set(j, k, 0.0)
						} else {
							model.Set(j, k, (float64(len(yuma.GetSubprocesses())) + 1.0) * -1.0)
						}
					}
				}
				updates = mat.NewDense(int(math.Exp2(float64(len(yuma.GetSubprocesses())))), int(math.Exp2(float64(len(yuma.GetSubprocesses()) - 1)) + 1), nil)
				//history = mat.NewDense(int(math.Exp2(float64(len(yuma.GetSubprocesses())))), int(math.Exp2(float64(len(yuma.GetSubprocesses()) - 1)) + 1), nil)
				//memory = mat.NewDense(int(math.Exp2(float64(len(yuma.GetSubprocesses())))), int(math.Exp2(float64(len(yuma.GetSubprocesses()) - 1)) + 1), nil)
				//timestamps = mat.NewDense(int(math.Exp2(float64(len(yuma.GetSubprocesses())))), int(math.Exp2(float64(len(yuma.GetSubprocesses()) - 1)) + 1), nil)
			}

			//behaviorPolicy := NewNegatedEpsilonGreedyPolicy(0.1)
			//behaviorPolicy := NewHistoricPolicy()
			behaviorPolicy := NewNaraEpsilonGreedyPolicy(0.3)
			targetPolicy := NewNaraGreedyPolicy()
			treeBackup := NewTreeBackup(yuma, behaviorPolicy, targetPolicy, 1, 0.5, 1, 24 * time.Hour, 0)
			treeBackup.SetN((len(yuma.GetSubprocesses())))

			treePolicy := NewNaraTreePolicy(0.3)
			expansionPolicy := NewNaraExpansionPolicy()
			selectionPolicy := NewNaraGreedyPolicy()
			rationalThinking := NewNara(yuma, treeBackup, treePolicy, expansionPolicy, selectionPolicy, 1 * time.Minute, 1, 0.8, 0.5, 1, 24 * time.Hour)
			//rationalThinking := NewNara(yuma, treeBackup, treePolicy, selectionPolicy, memory, nil, 5 * time.Minute, 1, 1, 0.5)
			yuma.SetRationalThinking(rationalThinking)

			behaviorPolicy.SetRationalThinking(treeBackup)
			targetPolicy.SetRationalThinking(treeBackup)
			treePolicy.SetRationalThinking(rationalThinking)
			expansionPolicy.SetRationalThinking(rationalThinking)
			selectionPolicy.SetRationalThinking(rationalThinking)

			//yuma.GetRationalThinking().Learn(int(target))
			eo := yuma.GetRationalThinking().Solve(int(target), model, updates, history, memory, timestamps)
			yuma.SetSolution(yuma.GetConfigurations()[target], eo)
			if err := yuma.ExportExecutionOrder(eo, float64(target)); err != nil {
				fmt.Println("EXPORT EXECUTION ORDER ERROR: " + err.Error())
				return err
			}
			if err := yuma.ExportServiceTree(target); err != nil {
				fmt.Println("EXPORT SERVICE TREE ERROR: " + err.Error())
				return err
			}
		}
		sim++
		i++
	}

	if err := yuma.GetEnvironment().CleanUp(); err != nil {
		fmt.Println("DIRECTORY DELETE ERROR: " + err.Error())
		return err
	}

	return nil
}

/* func (yuma *Yuma) DetermineMinimalExecutionOrder() error {
	for i := 0; i < len(yuma.GetSubprocesses()); i++ {
		target := math.Exp2(float64(i))
		pathExecutionOrder := "playbook_" + strconv.Itoa(int(math.Exp2(float64(i)))) + ".yml"
		
		mEO := yuma.GetRationalThinking().Solve(int(target))
		if err := yuma.GetEnvironment().GetExecutor().CreateExecutionOrder(0, pathExecutionOrder, "", mEO, "create", "localhost"); err != nil {
			return err
		}
	}

	return nil
} */

func (yuma *Yuma) ExportExecutionOrder(eo []string, target float64) error {
	pathExecutionOrder := "playbook_" + strconv.Itoa(int(target)) + ".yml"
	if err := yuma.GetEnvironment().GetExecutor().CreateExecutionOrder(0, pathExecutionOrder, "", eo, "create", "localhost"); err != nil {
		return err
	}

	return nil
}

func (yuma *Yuma) ExportServiceTree(target int) error {
	root := NewNode(yuma.GetConfigurations()[target])
	servicetree := DetermineServiceTree(root, yuma.GetConfigurations()[target], yuma.GetSolutions())
	export := PrintServiceTree(servicetree, "", "")

	file, err := os.OpenFile("servicetree_" + strconv.Itoa(target), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	defer file.Close()

	if _, err := file.WriteString(export); err != nil {
		return err
	}

	return nil
}