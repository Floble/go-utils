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

type DoubleQLearning struct {
	yuma *Yuma
	policy Policy
	memory *mat.Dense
	episodes int
	alpha float64
	gamma float64
}

func NewDoubleQLearning(yuma *Yuma, policy Policy, episodes int, alpha float64, gamma float64) *DoubleQLearning {
	dql := new(DoubleQLearning)
	dql.yuma = yuma
	dql.policy = policy
	dql.episodes = episodes
	dql.alpha = alpha
	dql.gamma = gamma

	return dql
}

func (dql *DoubleQLearning) GetYuma() *Yuma {
	return dql.yuma
}

func (dql *DoubleQLearning) GetMemory() *mat.Dense {
	return dql.memory
}

func (dql *DoubleQLearning) GetPolicy() Policy {
	return dql.policy
}

func (dql *DoubleQLearning) GetEpisodes() int {
	return dql.episodes
}

func (dql *DoubleQLearning) GetAlpha() float64 {
	return dql.alpha
}

func (dql *DoubleQLearning) GetGamma() float64 {
	return dql.gamma
}

func (dql *DoubleQLearning) Learn(target int) error {
	q1 := mat.NewDense(int(math.Exp2(float64(len(dql.GetYuma().GetSubprocesses())))), int(math.Exp2(float64(len(dql.GetYuma().GetSubprocesses()) - 1)) + 1), nil)
	q2 := mat.NewDense(int(math.Exp2(float64(len(dql.GetYuma().GetSubprocesses())))), int(math.Exp2(float64(len(dql.GetYuma().GetSubprocesses()) - 1)) + 1), nil)
	model := dql.GetYuma().GetModel(target)
	policy := dql.GetPolicy()
	
	var err error
	var success bool
	var reward float64
	var successor int

	exportResults := fmt.Sprintln(dql.GetYuma().GetSubprocesses()) + "\n"
	if err = dql.log(exportResults, "results.txt"); err != nil {
		return err
	}
	exportResults = ""

	choices := []randutil.Choice{randutil.Choice{5, 0}, randutil.Choice{5, 1}}

	// Initialize Q1(s, a) and Q2(s, a)
	dql.initializeQ(q1, target)
	dql.initializeQ(q2, target)
	dql.initializeMemory()
	// Repeat (for each episode)
	for i := 0; i < dql.GetEpisodes(); i++ {
		exportResults += "++++++++++++++++++++++++++\n"
		exportResults += fmt.Sprintf("Episode: %d\n", i + 1)
		exportResults += "++++++++++++++++++++++++++\n\n"
		if err = dql.log(exportResults, "results.txt"); err != nil {
			return err
		}
		exportResults = ""

		// Initialize S
		if dql.GetYuma().GetEnvironment().GetInstance(target) != nil {
			dql.GetYuma().GetEnvironment().DeleteInstance(target)
		}
		state, path := dql.initializeState()
		// Repeat (for each step of episode until S is target state)
		for state & target == 0 {
			f := mat.Formatted(q1, mat.Prefix("     "), mat.Squeeze())
			exportResults += fmt.Sprintf("\nQ1 = %v\n\n\n", f)
			if err = dql.log(exportResults, "results.txt"); err != nil {
				return err
			}
			exportResults = ""
			f = mat.Formatted(q2, mat.Prefix("     "), mat.Squeeze())
			exportResults += fmt.Sprintf("\nQ2 = %v\n\n\n", f)
			if err = dql.log(exportResults, "results.txt"); err != nil {
				return err
			}
			exportResults = ""

			// Choose A from S using the policy epsilon-greedy in Q1 + Q2
			model.Add(q1, q2)
			policy.DerivePolicy(model, nil)
			c, _ := randutil.WeightedChoice(policy.GetSuggestions()[state])
			action := c.Item.(int)

			f = mat.Formatted(model, mat.Prefix("        "), mat.Squeeze())
			exportResults += fmt.Sprintf("\nModel = %v\n\n\n", f)
			if err = dql.log(exportResults, "results.txt"); err != nil {
				return err
			}
			exportResults = ""

			exportResults += fmt.Sprintf("State: %d\n", state)
			exportResults += "Path: "
			exportResults += fmt.Sprintln(path)
			exportResults += "Policy: "
			exportResults += fmt.Sprintln(dql.GetPolicy().GetSuggestions())
			exportResults += fmt.Sprintf("Action: %d\n", action)
			if err = dql.log(exportResults, "results.txt"); err != nil {
				return err
			}
			exportResults = ""

			// Take action A, observe R, S'
			if dql.GetMemory().At(state, action) == 0.0 {
				err, success, reward, successor = dql.GetYuma().GetEnvironment().TakeAction(target, state, action, path, success)
				dql.GetMemory().Set(state, action, reward)
				if err != nil {
					return err
				} else if success {
					path = append(path, dql.GetYuma().GetConfigurations()[action])
				}
			} else {
				if dql.GetMemory().At(state, action) == -1.0 {
					success = false
					reward = dql.GetMemory().At(state, action)
					successor = state | action
					path = append(path, dql.GetYuma().GetConfigurations()[action])
				} else if dql.GetMemory().At(state, action) == -10.0 {
					success = false
					reward = dql.GetMemory().At(state, action)
					successor = state
				}
			}

			// With 0.5 probability
			c, _ = randutil.WeightedChoice(choices)
			var maxAction int
			var qUpdate, qSuccessor float64
			if c.Item.(int) == 1 {
				exportResults += "Q Selected: " + fmt.Sprintf("Q1\n")
				if err = dql.log(exportResults, "results.txt"); err != nil {
					return err
				}
				exportResults = ""

				// Q1(S, A) = Q1(S, A) + alpha * (R + gamma * Q2(S', argmax_a Q1(S', a)) - Q1(S, A))
				maxAction = dql.ArgMaxAction(q1, successor)
				qSuccessor = q2.At(successor, maxAction)
				qUpdate = q1.At(state, action) + dql.GetAlpha() * (reward + (dql.GetGamma() * qSuccessor - q1.At(state, action)))
				q1.Set(state, action, qUpdate)
			} else {
				exportResults += "Q Selected: " + fmt.Sprintf("Q2\n")
				if err = dql.log(exportResults, "results.txt"); err != nil {
					return err
				}
				exportResults = ""

				// Q2(S, A) = Q2(S, A) + alpha * (R + gamma * Q1(S', argmax_a Q2(S', a)) - Q2(S, A))
				maxAction = dql.ArgMaxAction(q2, successor)
				qSuccessor = q1.At(successor, maxAction)
				qUpdate = q2.At(state, action) + dql.GetAlpha() * (reward + (dql.GetGamma() * qSuccessor - q2.At(state, action)))
				q2.Set(state, action, qUpdate)
			}

			// S = S'
			state = successor

			exportResults += "Reward: " + fmt.Sprintf("%f", reward) + "\n"
			exportResults += "Q_Successor: " + fmt.Sprintf("%f", qSuccessor) + "\n"
			exportResults += "Q_Update: " + fmt.Sprintf("%f", qUpdate) + "\n"
			exportResults += "Successor: " + fmt.Sprintf("%d\n\n", successor)
			if err = dql.log(exportResults, "results.txt"); err != nil {
				return err
			}
			exportResults = ""
		}
	}

	exportResults += "++++++++++++++++++++++++++\n"
	exportResults += fmt.Sprintf("Completed Learning\n")
	exportResults += "++++++++++++++++++++++++++\n\n"
	if err = dql.log(exportResults, "results.txt"); err != nil {
		return err
	}

	return nil
}

func (dql *DoubleQLearning) Solve(target int) []string {
	state := dql.GetYuma().GetStartState()
	solution := make([]string, 0)
	model := dql.GetYuma().GetModel(target)

	for state & target == 0 {
		action := dql.ArgMaxAction(model, state)
		subprocess := dql.GetYuma().GetConfigurations()[action]
		solution = append(solution, subprocess)
		state = state | action
	}

	return solution
}

func (dql *DoubleQLearning) initializeQ(q *mat.Dense, target int) {
	for i := 0; i < int(math.Exp2(float64(len(dql.GetYuma().GetSubprocesses())))); i++ {
		for j := 0; j < int(math.Exp2(float64(len(dql.GetYuma().GetSubprocesses()) - 1)) + 1); j++ {
			if i & target != 0 || j == 0 {
				q.Set(i, j, 0.0)
			} else {
				rand.Seed(time.Now().UnixNano())
				q.Set(i, j, rand.Float64() * -1.0)
			}
		}
	}
}

func (dql *DoubleQLearning) initializeState() (int, []string) {
	return dql.GetYuma().GetStartState(), make([]string, 0)
}

func (dql *DoubleQLearning) initializeMemory() {
	dql.memory = mat.NewDense(int(math.Exp2(float64(len(dql.GetYuma().GetSubprocesses())))), int(math.Exp2(float64(len(dql.GetYuma().GetSubprocesses()) - 1)) + 1), nil)
}

func (dql *DoubleQLearning) ArgMaxAction(q *mat.Dense, state int) int {
	maxQ := math.MaxFloat64 * -1.0
	maxAction := 0
	for _, action := range dql.GetYuma().Actions(state) {
		if q.At(state, action) > maxQ {
			maxQ = q.At(state, action)
			maxAction = action
		}
	}

	return maxAction
}

func (dql *DoubleQLearning) log(exportResults string, filePath string) error {
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