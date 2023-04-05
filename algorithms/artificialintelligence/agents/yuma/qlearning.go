package yuma

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"
	"github.com/jmcvetta/randutil"
	"gonum.org/v1/gonum/mat"
)

type QLearning struct {
	yuma *Yuma
	policy Policy
	memory *mat.Dense
	episodes int
	alpha float64
	gamma float64
}

func NewQLearning(yuma *Yuma, policy Policy, episodes int, alpha float64, gamma float64) *QLearning {
	ql := new(QLearning)
	ql.yuma = yuma
	ql.policy = policy
	ql.episodes = episodes
	ql.alpha = alpha
	ql.gamma = gamma

	return ql
}

func (ql *QLearning) GetYuma() *Yuma {
	return ql.yuma
}

func (ql *QLearning) GetMemory() *mat.Dense {
	return ql.memory
}

func (ql *QLearning) GetPolicy() Policy {
	return ql.policy
}

func (ql *QLearning) GetEpisodes() int {
	return ql.episodes
}

func (ql *QLearning) GetAlpha() float64 {
	return ql.alpha
}

func (ql *QLearning) GetGamma() float64 {
	return ql.gamma
}

func (ql *QLearning) Learn(target int) error {
	model := ql.GetYuma().GetModel(target)
	policy := ql.GetPolicy()

	var err error
	var success bool
	var reward float64
	var successor int
	
	exportResults := fmt.Sprintln(ql.GetYuma().GetSubprocesses()) + "\n"
	if err = ql.log(exportResults, "results.txt"); err != nil {
		return err
	}
	exportResults = ""

	// Initialize Q(s, a)
	ql.initializeQ(model, target)
	ql.initializeMemory()
	// Repeat (for each episode)
	for i := 0; i < ql.GetEpisodes(); i++ {
		exportResults += "++++++++++++++++++++++++++\n"
		exportResults += fmt.Sprintf("Episode: %d\n", i + 1)
		exportResults += "++++++++++++++++++++++++++\n\n"
		if err = ql.log(exportResults, "results.txt"); err != nil {
			return err
		}
		exportResults = ""

		// Initialize S
		if ql.GetYuma().GetEnvironment().GetInstance(target) != nil {
			ql.GetYuma().GetEnvironment().DeleteInstance(target)
		}
		state, path := ql.initializeState()
		// Repeat (for each step of episode until S is target state)
		for state & target == 0 {
			f := mat.Formatted(model, mat.Prefix("        "), mat.Squeeze())
			exportResults += fmt.Sprintf("\nModel = %v\n\n\n", f)
			if err = ql.log(exportResults, "results.txt"); err != nil {
				return err
			}
			exportResults = ""

			// Choose A from S using policy derived from Q (e.g., epsilon-greedy)
			policy.DerivePolicy(model, nil)
			c, _ := randutil.WeightedChoice(policy.GetSuggestions()[state])
			action := c.Item.(int)

			exportResults += fmt.Sprintf("State: %d\n", state)
			exportResults += "Path: "
			exportResults += fmt.Sprintln(path)
			exportResults += "Policy: "
			exportResults += fmt.Sprintln(ql.GetPolicy().GetSuggestions())
			exportResults += fmt.Sprintf("Action: %d\n", action)
			if err = ql.log(exportResults, "results.txt"); err != nil {
				return err
			}
			exportResults = ""

			// Take action A, observe R, S'
			if ql.GetMemory().At(state, action) == 0.0 {
				err, success, reward, successor = ql.GetYuma().GetEnvironment().TakeAction(target, state, action, path, success)
				ql.GetMemory().Set(state, action, reward)
				if err != nil {
					return err
				} else if success {
					path = append(path, ql.GetYuma().GetConfigurations()[action])
				}
			} else {
				if ql.GetMemory().At(state, action) == -1.0 {
					success = false
					reward = ql.GetMemory().At(state, action)
					successor = state | action
					path = append(path, ql.GetYuma().GetConfigurations()[action])
				} else if ql.GetMemory().At(state, action) == -10.0 {
					success = false
					reward = ql.GetMemory().At(state, action)
					successor = state
				}
			}

			maxAction := ql.ArgMaxAction(model, successor)
			qSuccessor := model.At(successor, maxAction)

			// Q(S, A) = Q(S, A) + alpha * [reward + gamma * max_Q(S', a) - Q(S, A)]	
			qUpdate := model.At(state, action) + ql.GetAlpha() * (reward + (ql.GetGamma() * qSuccessor - model.At(state, action)))
			model.Set(state, action, qUpdate)

			// S = S'
			state = successor

			exportResults += "Reward: " + fmt.Sprintf("%f", reward) + "\n"
			exportResults += "Q_Successor: " + fmt.Sprintf("%f", qSuccessor) + "\n"
			exportResults += "Q_Update: " + fmt.Sprintf("%f", qUpdate) + "\n"
			exportResults += "Successor: " + fmt.Sprintf("%d\n\n", successor)
			if err = ql.log(exportResults, "results.txt"); err != nil {
				return err
			}
			exportResults = ""
		}
	}

	exportResults += "++++++++++++++++++++++++++\n"
	exportResults += fmt.Sprintf("Completed Learning\n")
	exportResults += "++++++++++++++++++++++++++\n\n"
	if err = ql.log(exportResults, "results.txt"); err != nil {
		return err
	}

	return nil
}

func (ql *QLearning) Solve(target int) []string {
	state := ql.GetYuma().GetStartState()
	solution := make([]string, 0)
	model := ql.GetYuma().GetModel(target)

	for state & target == 0 {
		action := ql.ArgMaxAction(model, state)
		subprocess := ql.GetYuma().GetConfigurations()[action]
		solution = append(solution, subprocess)
		state = state | action
	}

	return solution
}

func (ql *QLearning) initializeQ(q *mat.Dense, target int) {
	for i := 0; i < int(math.Exp2(float64(len(ql.GetYuma().GetSubprocesses())))); i++ {
		for j := 0; j < int(math.Exp2(float64(len(ql.GetYuma().GetSubprocesses()) - 1)) + 1); j++ {
			if i & target != 0 || j == 0 {
				q.Set(i, j, 0.0)
			} else {
				rand.Seed(time.Now().UnixNano())
				q.Set(i, j, rand.Float64() * -1.0)
			}
		}
	}
}

func (ql *QLearning) initializeState() (int, []string) {
	return ql.GetYuma().GetStartState(), make([]string, 0)
}

func (ql *QLearning) initializeMemory() {
	ql.memory = mat.NewDense(int(math.Exp2(float64(len(ql.GetYuma().GetSubprocesses())))), int(math.Exp2(float64(len(ql.GetYuma().GetSubprocesses()) - 1)) + 1), nil)
}

func (ql *QLearning) ArgMaxAction(q *mat.Dense, state int) int {
	maxQ := math.MaxFloat64 * -1.0
	maxAction := 0
	for _, action := range ql.GetYuma().Actions(state) {
		if q.At(state, action) > maxQ {
			maxQ = q.At(state, action)
			maxAction = action
		}
	}

	return maxAction
}

func (ql *QLearning) log(exportResults string, filePath string) error {
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