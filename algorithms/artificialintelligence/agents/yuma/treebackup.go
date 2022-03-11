package yuma

import (
	"fmt"
	"math"
	"math/rand"
	"time"
	"os"
	"gonum.org/v1/gonum/mat"
	"github.com/jmcvetta/randutil"
)

type TreeBackup struct {
	yuma *Yuma
	behaviorPolicy Policy
	targetPolicy Policy
	memory *mat.Dense
	episodes int
	alpha float64
	gamma float64
	n int
}

type Decision struct {
	state int
	action int
}

func NewTreeBackup(yuma *Yuma, behaviorPolicy Policy, targetPolicy Policy, episodes int, alpha float64, gamma float64, n int) *TreeBackup {
	tb := new(TreeBackup)
	tb.yuma = yuma
	tb.behaviorPolicy = behaviorPolicy
	tb.targetPolicy = targetPolicy
	tb.episodes = episodes
	tb.alpha = alpha
	tb.gamma = gamma
	tb.n = n

	return tb
}

func NewDecision(state int, action int) *Decision {
	decision := new(Decision)
	decision.state = state
	decision.action = action

	return decision
}

func (tb *TreeBackup) GetYuma() *Yuma {
	return tb.yuma
}

func (tb *TreeBackup) GetMemory() *mat.Dense {
	return tb.memory
}

func (tb *TreeBackup) GetBehaviorPolicy() Policy {
	return tb.behaviorPolicy
}

func (tb *TreeBackup) GetTargetPolicy() Policy {
	return tb.targetPolicy
}

func (tb *TreeBackup) GetEpisodes() int {
	return tb.episodes
}

func (tb *TreeBackup) GetAlpha() float64 {
	return tb.alpha
}

func (tb *TreeBackup) GetGamma() float64 {
	return tb.gamma
}

func (tb *TreeBackup) GetN() int {
	return tb.n
}

func (tb *TreeBackup) SetN(n int) {
	tb.n = n
}

func (decision *Decision) GetState() int {
	return decision.state
}

func (decision *Decision) GetAction() int {
	return decision.action
}

func (decision *Decision) SetState(state int) {
	decision.state = state
}

func (decision *Decision) SetAction(action int) {
	decision.action = action
}

func (tb *TreeBackup) Learn(target int) error {
	rand.Seed(time.Now().Unix())

	model := tb.GetYuma().GetModel()
	behaviorPolicy := tb.GetBehaviorPolicy()
	targetPolicy := tb.GetTargetPolicy()

	var err error
	var success bool
	var reward float64
	var successor int
	
	exportResults := fmt.Sprintln(tb.GetYuma().GetSubprocesses()) + "\n"
	if err = tb.log(exportResults, "results_" + string(target) + ".txt"); err != nil {
		return err
	}
	exportResults = ""

	// Initialize Q(s, a) arbitrarily, for all s, a
	tb.initializeQ(model, target)
	tb.initializeMemory()
	// Initialize pi to be greedy with respect to Q
	behaviorPolicy.DerivePolicy(model)
	targetPolicy.DerivePolicy(model)
	// All store and access operations can take their index mod n + 1
	decisions := make(map[int]*Decision, 0)
	rewards := make(map[int]float64, 0)

	// Loop for each episode
	for i := 0; i < tb.GetEpisodes(); i++ {
		exportResults += "++++++++++++++++++++++++++\n"
		exportResults += fmt.Sprintf("Episode: %d\n", i + 1)
		exportResults += "++++++++++++++++++++++++++\n\n"
		if err = tb.log(exportResults, "results_" + string(target) + ".txt"); err != nil {
			return err
		}
		exportResults = ""

		// Initialize and store S0 != terminal
		if tb.GetYuma().GetEnvironment().GetInstance() != nil {
			tb.GetYuma().GetEnvironment().DeleteInstance()
		}
		state, path := tb.initializeState()
		// Choose an action A0 arbitrarily as a function of S0; Store A0
		c, _ := randutil.WeightedChoice(behaviorPolicy.GetSuggestions()[state])
		action := c.Item.(int)
		decisions[0] = NewDecision(state, action)

		// T = infinity
		terminal := math.MaxInt64
		var tau int

		// Loop for t = 0, 1, 2, ...
		t := 0

		// Loop until tau = T - 1
		for tau < terminal - 1 {
			// Tau is the time whose estimate is being updated
			tau = t + 1 - tb.GetN()
			if t < terminal {
				f := mat.Formatted(model, mat.Prefix("        "), mat.Squeeze())
				exportResults += fmt.Sprintf("\nModel = %v\n\n\n", f)
				exportResults += fmt.Sprintf("T: %d\n", t)
				exportResults += fmt.Sprintf("S_%d: %d\n", t, decisions[t].GetState())
				exportResults += fmt.Sprintf("A_%d: %d\n\n", t, decisions[t].GetAction())
				if err = tb.log(exportResults, "results_" + string(target) + ".txt"); err != nil {
					return err
				}
				exportResults = ""
				
				// Take action At; observe and store the next reward and state as Rt+1, St+1
				if tb.GetMemory().At(decisions[t].GetState(), decisions[t].GetAction()) == 0.0 {
					err, success, reward, successor = tb.GetYuma().GetEnvironment().TakeAction(decisions[t].GetState(), decisions[t].GetAction(), path, success)
					if err != nil {
						return err
					} else if success {
						path = append(path, tb.GetYuma().GetConfigurations()[decisions[t].GetAction()])
					}
					decisions[t + 1] = NewDecision(successor, 0)
					rewards[t + 1] = reward
					tb.GetMemory().Set(decisions[t].GetState(), decisions[t].GetAction(), rewards[t + 1])
					if rewards[t + 1] == -10.0 {
						tau = t
					}
				} else {
					if tb.GetMemory().At(decisions[t].GetState(), decisions[t].GetAction()) == -1.0 {
						success = false
						reward = tb.GetMemory().At(decisions[t].GetState(), decisions[t].GetAction())
						rewards[t + 1] = reward
						successor = decisions[t].GetState() | decisions[t].GetAction()
						decisions[t + 1] = NewDecision(successor, 0)
						path = append(path, tb.GetYuma().GetConfigurations()[decisions[t].GetAction()])
					} else if tb.GetMemory().At(decisions[t].GetState(), decisions[t].GetAction()) == -10.0 {
						success = false
						reward = tb.GetMemory().At(decisions[t].GetState(), decisions[t].GetAction())
						rewards[t + 1] = reward
						successor = state
						decisions[t + 1] = NewDecision(successor, 0)
						tau = t
					}
				}

				exportResults += "R_" + fmt.Sprintf("%d", t + 1) + ": " + fmt.Sprintf("%f", reward) + "\n\n"
				if err = tb.log(exportResults, "results_" + string(target) + ".txt"); err != nil {
					return err
				}
				exportResults = ""

				// If St+1 is terminal
				if (decisions[t + 1].GetState() & target > 0) || (decisions[t].GetState() == decisions[t + 1].GetState()) {
					terminal = t + 1
				} else {
					// Choose an action At+1 arbitrarily as a function of St+1; Store At+1
					c, _ := randutil.WeightedChoice(behaviorPolicy.GetSuggestions()[decisions[t + 1].GetState()])
					action := c.Item.(int)
					decisions[t + 1].SetAction(action)
				}
			}
			var totalReturn float64
			if tau >= 0 {
				exportResults += fmt.Sprintf("T: %d\n", t)
				exportResults += fmt.Sprintf("Terminal: %d\n", terminal)
				exportResults += fmt.Sprintf("\nTau: %d - State: %d - Action: %d\n\n", tau, decisions[tau].GetState(), decisions[tau].GetAction())
				if err = tb.log(exportResults, "results_" + string(target) + ".txt"); err != nil {
					return err
				}
				exportResults = ""

				if t + 1 >= terminal {
					// G = RT
					totalReturn = rewards[terminal]
					
					exportResults += fmt.Sprintf("G = %f\n\n", totalReturn)
					if err = tb.log(exportResults, "results_" + string(target) + ".txt"); err != nil {
						return err
					}
					exportResults = ""
				} else {
					// G = Rt+1 + gamma * SUM_over_a(pi(a|St+1) * Q(St+1, a))
					var totalActionValue float64

					exportResults += fmt.Sprintf("G = %f + %f * (0 ", rewards[t + 1], tb.GetGamma())
					if err = tb.log(exportResults, "results_" + string(target) + ".txt"); err != nil {
						return err
					}
					exportResults = ""

					for _, action := range tb.GetYuma().Actions(decisions[t + 1].GetState()) {
						weight := float64(targetPolicy.GetWeight(decisions[t + 1].GetState(), action)) / 10.0
						totalActionValue += float64(weight) * model.At(decisions[t + 1].GetState(), action)

						exportResults += fmt.Sprintf("+ %f * %f", weight, model.At(decisions[t + 1].GetState(), action))
						if err = tb.log(exportResults, "results_" + string(target) + ".txt"); err != nil {
							return err
						}
						exportResults = ""
					}
					totalReturn = rewards[t + 1] + tb.GetGamma() * totalActionValue

					exportResults += fmt.Sprintf(") = %f\n\n", totalReturn)
					if err = tb.log(exportResults, "results_" + string(target) + ".txt"); err != nil {
						return err
					}
					exportResults = ""
				}
				// Loop for k = min(t, T - 1) down through tau + 1
				k := int(math.Min(float64(t), float64(terminal - 1)))
				for i := k; i >= tau + 1; i-- {
					exportResults += fmt.Sprintf("S_%d: %d\n", k, decisions[k].GetState())
					if err = tb.log(exportResults, "results_" + string(target) + ".txt"); err != nil {
						return err
					}
					exportResults = ""

					// G = Rk + gamma * SUM_over_a!=Ak(pi(a|Sk) * Q(Sk, a) + gamma * pi(Ak|Sk) * G)
					var totalActionValue float64

					exportResults += fmt.Sprintf("G = %f + %f * (0", rewards[k], tb.GetGamma())
					if err = tb.log(exportResults, "results_" + string(target) + ".txt"); err != nil {
						return err
					}
					exportResults = ""

					for _, action := range tb.GetYuma().Actions(decisions[k].GetState()) {
						if action == decisions[k].GetAction() {
							continue
						}
						weight := float64(targetPolicy.GetWeight(decisions[k].GetState(), action)) / 10.0
						totalActionValue += float64(weight) * model.At(decisions[k].GetState(), action)

						exportResults += fmt.Sprintf(" + %f * %f", weight, model.At(decisions[k].GetState(), action))
						if err = tb.log(exportResults, "results_" + string(target) + ".txt"); err != nil {
							return err
						}
						exportResults = ""
					}

					exportResults += fmt.Sprintf(") + %f * %f * %f = ", tb.GetGamma(), float64(targetPolicy.GetWeight(decisions[k].GetState(), decisions[k].GetAction())) / 10.0, totalReturn)
					if err = tb.log(exportResults, "results_" + string(target) + ".txt"); err != nil {
						return err
					}
					exportResults = ""

					totalReturn = rewards[k] + tb.GetGamma() * totalActionValue + tb.GetGamma() * (float64(targetPolicy.GetWeight(decisions[k].GetState(), decisions[k].GetAction())) / 10.0) * totalReturn

					exportResults += fmt.Sprintf("%f\n\n", totalReturn)
					if err = tb.log(exportResults, "results_" + string(target) + ".txt"); err != nil {
						return err
					}
					exportResults = ""
				}
				// Q(Stau, Atau) = Q(Stau, Atau) + alpha * (G - Q(Stau, Atau))
				exportResults += fmt.Sprintf("Q_Update = %f + %f * (%f - %f) = ", model.At(decisions[tau].GetState(), decisions[tau].GetAction()), tb.GetAlpha(), totalReturn, model.At(decisions[tau].GetState(), decisions[tau].GetAction()))
				if err = tb.log(exportResults, "results_" + string(target) + ".txt"); err != nil {
					return err
				}
				exportResults = ""

				qUpdate := model.At(decisions[tau].GetState(), decisions[tau].GetAction()) + tb.GetAlpha() * (totalReturn - model.At(decisions[tau].GetState(), decisions[tau].GetAction()))

				exportResults += fmt.Sprintf("%f\n\n", qUpdate)
				if err = tb.log(exportResults, "results_" + string(target) + ".txt"); err != nil {
					return err
				}
				exportResults = ""

				model.Set(decisions[tau].GetState(), decisions[tau].GetAction(), qUpdate)
				// If pi is being learned, then ensure that pi(.|Stau) is greedy wrt Q
				behaviorPolicy.DerivePolicy(model)
				targetPolicy.DerivePolicy(model)
			}
			t += 1
		}
	}

	exportResults += "++++++++++++++++++++++++++\n"
	exportResults += fmt.Sprintf("Completed Learning\n")
	exportResults += "++++++++++++++++++++++++++\n\n"
	if err = tb.log(exportResults, "results_" + string(target) + ".txt"); err != nil {
		return err
	}

	f := mat.Formatted(tb.GetMemory(), mat.Prefix("         "), mat.Squeeze())
	exportMemory := fmt.Sprintf("\nMemory = %v\n\n\n", f)
	if err = tb.log(exportMemory, "memory.txt"); err != nil {
		return err
	}

	return nil
}

func (tb *TreeBackup) Solve(target int) []string {
	state := tb.GetYuma().GetStartState()
	solution := make([]string, 0)
	model	 := tb.GetYuma().GetModel()

	for state & target == 0 {
		action := tb.ArgMaxAction(model, state)
		subprocess := tb.GetYuma().GetConfigurations()[action]
		solution = append(solution, subprocess)
		state = state | action
	}

	return solution
}

func (tb *TreeBackup) initializeQ(q *mat.Dense, target int) {
	for i := 0; i < int(math.Exp2(float64(len(tb.GetYuma().GetSubprocesses())))); i++ {
		for j := 0; j < int(math.Exp2(float64(len(tb.GetYuma().GetSubprocesses()) - 1)) + 1); j++ {
			if i & target != 0 || j == 0 {
				q.Set(i, j, 0.0)
			} else {
				rand.Seed(time.Now().UnixNano())
				q.Set(i, j, rand.Float64() * -1.0)
			}
		}
	}
}

func (tb *TreeBackup) initializeState() (int, []string) {
	return tb.GetYuma().GetStartState(), make([]string, 0)
}

func (tb *TreeBackup) initializeMemory() {
	tb.memory = mat.NewDense(int(math.Exp2(float64(len(tb.GetYuma().GetSubprocesses())))), int(math.Exp2(float64(len(tb.GetYuma().GetSubprocesses()) - 1)) + 1), nil)
}

func (tb *TreeBackup) ArgMaxAction(q *mat.Dense, state int) int {
	maxQ := math.MaxFloat64 * -1.0
	maxAction := 0
	for _, action := range tb.GetYuma().Actions(state) {
		if q.At(state, action) > maxQ {
			maxQ = q.At(state, action)
			maxAction = action
		}
	}

	return maxAction
}

func (tb *TreeBackup) log(exportResults string, filePath string) error {
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