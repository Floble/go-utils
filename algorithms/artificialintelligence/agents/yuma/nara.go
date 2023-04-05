package yuma

import (
	"fmt"
	"os"
	"strconv"
	"math"
	"math/rand"
	"time"
	"gonum.org/v1/gonum/mat"
	"github.com/jmcvetta/randutil"
)

type Nara struct {
	yuma *Yuma
	treebackup *TreeBackup
	treePolicy Policy
	selectionPolicy Policy
	maxTime time.Duration
	maxHistory int
	alpha float64
	gamma float64
}

func NewNara(yuma *Yuma, treebackup *TreeBackup, treePolicy Policy, selectionPolicy Policy, maxTime time.Duration, maxHistory int, gamma float64, alpha float64) *Nara {
	nara := new(Nara)
	nara.yuma = yuma
	nara.treebackup = treebackup
	nara.treePolicy = treePolicy
	nara.selectionPolicy = selectionPolicy
	nara.maxTime = maxTime
	nara.maxHistory = maxHistory
	nara.gamma = gamma
	nara.alpha = alpha

	return nara
}

func (nara *Nara) GetYuma() *Yuma {
	return nara.yuma
}

func (nara *Nara) GetTreeBackup() *TreeBackup {
	return nara.treebackup
}

func (nara *Nara) GetTreePolicy() Policy {
	return nara.treePolicy
}

func (nara *Nara) GetSelectionPolicy() Policy {
	return nara.selectionPolicy
}

func (nara *Nara) GetMaxTime() time.Duration {
	return nara.maxTime
}

func (nara *Nara) GetGamma() float64 {
	return nara.gamma
}

func (nara *Nara) GetAlpha() float64 {
	return nara.alpha
}

func (nara *Nara) GetMaxHistory() int {
	return nara.maxHistory
}

func (nara *Nara) Learn(target int, model, updates, history, memory *mat.Dense) error {
	return nil
}

func (nara *Nara) Solve(target int, model, updates, history, memory *mat.Dense) []string {
	solution := make([]string, 0)
	tree := make([]int, 0)
	rand.Seed(time.Now().Unix())

	treePolicy := nara.GetTreePolicy()
	selectionPolicy := nara.GetSelectionPolicy()

	treePolicy.DerivePolicy(model, updates)
	selectionPolicy.DerivePolicy(model, updates)

	state := nara.GetYuma().GetStartState()
	for state & target <= 0 {
		exportResults := "++++++++++++++++++++++++++\n"
		exportResults += fmt.Sprintf("State: %d\n", state)
		exportResults += "++++++++++++++++++++++++++\n\n"
		if err := nara.log(exportResults, "results_" + strconv.Itoa(target) + ".txt"); err != nil {
			return make([]string, 0)
		}
		exportResults = ""

		start := time.Now()
		maxHistory := math.Min(float64(len(nara.GetYuma().Actions(state))), float64(nara.GetMaxHistory()))
		for time.Since(start) < nara.GetMaxTime() || float64(distinctSumOfRow(history, state)) < maxHistory {
			exportResults += fmt.Sprintf("Time Elapsed: %v\n", time.Since(start).Minutes())
			if err := nara.log(exportResults, "results_" + strconv.Itoa(target) + ".txt"); err != nil {
				return make([]string, 0)
			}
			exportResults = ""

			tree = nil
			tree = append(tree, state)
			tree = nara.selection(target, model, history, treePolicy, state, tree)
			if nara.GetYuma().IsTerminal(tree[len(tree) - 1]) || (tree[len(tree) - 1] & target > 0) {
				exportResults += fmt.Sprintf("Tree: %v\n\n", tree)
				if err := nara.log(exportResults, "results_" + strconv.Itoa(target) + ".txt"); err != nil {
					return make([]string, 0)
				}
				exportResults = ""

				break
			}

			exportResults += fmt.Sprintf("MaxHistory: %f\n", maxHistory)
			exportResults += fmt.Sprintf("Sum of Row: %d\n", distinctSumOfRow(history, state))
			exportResults += fmt.Sprintf("Tree: %v\n\n", tree)
			if err := nara.log(exportResults, "results_" + strconv.Itoa(target) + ".txt"); err != nil {
				return make([]string, 0)
			}
			exportResults = ""

			exportResults += "++++++++++++++++++++++++++\n"
			exportResults += fmt.Sprintf("Modal: %d\n", tree[len(tree) - 1])
			exportResults += "++++++++++++++++++++++++++\n\n"
			if err := nara.log(exportResults, "results_" + strconv.Itoa(target) + ".txt"); err != nil {
				return make([]string, 0)
			}
			exportResults = ""

			if err := nara.simulation(target, model, updates, history, memory, tree, solution); err != nil {
				return make([]string, 0)
			}

			nara.backup(model, model, tree)
			treePolicy.DerivePolicy(model, history)
			selectionPolicy.DerivePolicy(model, history)
		}

		c, _ := randutil.WeightedChoice(selectionPolicy.GetSuggestions()[state])
		action := c.Item.(int)
		solution = append(solution, nara.GetYuma().GetConfigurations()[action])
		state = state | action

		exportResults += "++++++++++++++++++++++++++\n"
		exportResults += fmt.Sprintf("Selected Action: %d\n", action)
		exportResults += "++++++++++++++++++++++++++\n\n"
		if err := nara.log(exportResults, "results_" + strconv.Itoa(target) + ".txt"); err != nil {
			return make([]string, 0)
		}
		exportResults = ""

		exportResults += "++++++++++++++++++++++++++\n"
		exportResults += fmt.Sprintf("Current Q of YUMA:\n")
		exportResults += "++++++++++++++++++++++++++\n\n"
		q := mat.Formatted(model, mat.Prefix("        "), mat.Squeeze())
		exportResults += fmt.Sprintf("\nModel = %v\n\n\n", q)
		exportResults += fmt.Sprintln()
		if err := nara.log(exportResults, "results_" + strconv.Itoa(target) + ".txt"); err != nil {
			return make([]string, 0)
		}
		exportResults = ""

		exportResults += "++++++++++++++++++++++++++\n"
		exportResults += fmt.Sprintf("Current History of YUMA:\n")
		exportResults += "++++++++++++++++++++++++++\n\n"
		h := mat.Formatted(history, mat.Prefix("          "), mat.Squeeze())
		exportResults += fmt.Sprintf("\nHistory = %v\n\n\n", h)
		exportResults += fmt.Sprintln()
		if err := nara.log(exportResults, "results_" + strconv.Itoa(target) + ".txt"); err != nil {
			return make([]string, 0)
		}
		exportResults = ""
	}

	nara.GetYuma().SetModel(target, model)
	nara.GetYuma().SetUpdates(target, updates)
	nara.GetYuma().SetHistory(target, history)
	nara.GetYuma().SetMemory(target, memory)

	return solution
}

func (nara *Nara) initializeQ(q *mat.Dense, target int) {
	for i := 0; i < int(math.Exp2(float64(len(nara.GetYuma().GetSubprocesses())))); i++ {
		for j := 0; j < int(math.Exp2(float64(len(nara.GetYuma().GetSubprocesses()) - 1)) + 1); j++ {
			if i & target != 0 || j == 0 {
				q.Set(i, j, 0.0)
			} else {
				q.Set(i, j, (float64(len(nara.GetYuma().GetSubprocesses())) + 1.0) * -1.0)
			}
		}
	}
}

func (nara *Nara) ArgMaxAction(q *mat.Dense, state int, actions []int) int {
	maxQ := math.MaxFloat64 * -1.0
	tmp := rand.Intn(len(actions) - 0) + 0
	maxAction := actions[tmp]
	if q.At(state, maxAction) > maxQ {
		maxQ = q.At(state, maxAction)
	}
	
	for _, action := range actions {
		if q.At(state, action) > maxQ {
			maxQ = q.At(state, action)
			maxAction = action
		}
	}

	return maxAction
}

func (nara *Nara) selection(target int, model *mat.Dense, history *mat.Dense, treePolicy Policy, state int, tree []int) []int {
	if c, _ := randutil.WeightedChoice(treePolicy.GetSuggestions()[state]); c.Item != nil {
		action := c.Item.(int)
		maxHistory := math.Min(float64(len(nara.GetYuma().Actions(state))), float64(nara.GetMaxHistory()))
		for (!nara.GetYuma().IsTerminal(state)) && (state & target <= 0) && (float64(distinctSumOfRow(history, state)) >= maxHistory) && (state != tree[len(tree) - 1]) {
			tree = append(tree, state)
			state = state | action
			if !nara.GetYuma().IsTerminal(state) {
				c, _ := randutil.WeightedChoice(treePolicy.GetSuggestions()[state])
				action = c.Item.(int)
			}
		}
		//tree = append(tree, state)
	}

	return tree
}

func (nara *Nara) simulation(target int, model, updates, history, memory *mat.Dense, tree []int, path []string) error {
	leaf := tree[len(tree) - 1]
	nara.treebackup.SetRoot(leaf)
	if len(tree) > 1 {
		tmpPath := make([]string, 0)
		for j := len(tree) - 1; j >= 1; j-- {
			i := j - 1
			action := tree[j] ^ tree [i]
			tmpElement := nara.GetYuma().configurations[action]
			tmpPath = append(tmpPath, tmpElement)
		}
		path = append(path, tmpPath...)
		nara.treebackup.SetPath(path)
	} else {
		nara.treebackup.SetPath(path)
	}

	if err := nara.treebackup.Learn(target, model, updates, history, memory); err != nil {
		return err
	}
	action := nara.GetTreeBackup().ArgMaxAction(model, leaf, nara.GetYuma().Actions(leaf))
	model.Set(leaf, action, model.At(leaf, action))

	return nil
}

func (nara *Nara) backup(model *mat.Dense, tmpModel *mat.Dense, tree []int) {
	for i := len(tree) - 1; i >= 1; i-- {
		j := i - 1
		action := tree[j] ^ tree[i]
		reward := tmpModel.At(tree[j], action)
		totalReturn := 0.0
		for a := range nara.GetYuma().Actions(tree[i]) {
			totalReturn += float64(nara.GetTreePolicy().GetWeight(tree[i], a)) / 10.0 * tmpModel.At(tree[i], a)
		}
		totalReturn = reward + nara.GetGamma() * totalReturn
		qUpdate := tmpModel.At(tree[j], action) + nara.GetAlpha() * (totalReturn - tmpModel.At(tree[j], action))
		model.Set(tree[j], action, qUpdate)
	}
}

func (nara *Nara) log(exportResults string, filePath string) error {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	if _, err := file.WriteString(exportResults); err != nil {
		return err
	}

	defer file.Close()

	return nil
}

func distinctSumOfRow (history *mat.Dense, r int) int {
	_, c := history.Dims()
	sum := 0
	for i := 0; i < c; i++ {
		if history.At(r, i) >= 1 {
			sum++
		}
	}

	return sum
}

func sumOfRow (history *mat.Dense, r int) float64 {
	_, c := history.Dims()
	sum := 0.0
	for i := 0; i < c; i++ {
		sum += history.At(r, i)
	}

	return sum
}