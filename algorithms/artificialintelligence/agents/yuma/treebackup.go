package yuma

import (
	"fmt"
	"math"
	"math/rand"
	"time"
	"os"
	"strconv"
	"gonum.org/v1/gonum/mat"
	"github.com/jmcvetta/randutil"
)

type TreeBackup struct {
	yuma *Yuma
	behaviorPolicy Policy
	targetPolicy Policy
	root int
	path []string
	tree []int
	episodes int
	alpha float64
	gamma float64
	theta time.Duration
	n int
}

type Decision struct {
	state int
	action int
}

func NewTreeBackup(yuma *Yuma, behaviorPolicy Policy, targetPolicy Policy, episodes int, alpha float64, gamma float64, theta time.Duration, n int) *TreeBackup {
	tb := new(TreeBackup)
	tb.yuma = yuma
	tb.behaviorPolicy = behaviorPolicy
	tb.targetPolicy = targetPolicy
	tb.episodes = episodes
	tb.alpha = alpha
	tb.gamma = gamma
	tb.theta = theta
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

func (tb *TreeBackup) GetBehaviorPolicy() Policy {
	return tb.behaviorPolicy
}

func (tb *TreeBackup) GetTargetPolicy() Policy {
	return tb.targetPolicy
}

func (tb *TreeBackup) GetRoot() int {
	return tb.root
}

func (tb *TreeBackup) SetRoot(root int) {
	tb.root = root
}

func (tb *TreeBackup) GetPath() []string {
	return tb.path
}

func (tb *TreeBackup) SetPath(path []string) {
	tb.path = path
}

func (tb *TreeBackup) GetTree() []int {
	return tb.tree
}

func (tb *TreeBackup) SetTree(tree []int) {
	tb.tree = tree
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

func (tb *TreeBackup) GetTheta() time.Duration {
	return tb.theta
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

func (tb *TreeBackup) Learn(target int, model, updates, history, memory *mat.Dense, timestamps *mat.Dense) error {
	rand.Seed(time.Now().Unix())

	var err error
	var success bool
	var reward float64
	var successor int
	
	exportResults := fmt.Sprintln(tb.GetYuma().GetSubprocesses()) + "\n"
	if err = tb.log(exportResults, "results_" + strconv.Itoa(target) + ".txt"); err != nil {
		return err
	}
	exportResults = ""

	behaviorPolicy := tb.GetBehaviorPolicy()
	targetPolicy := tb.GetTargetPolicy()

	// Initialize pi to be greedy with respect to Q
	behaviorPolicy.DerivePolicy(model, updates)
	//behaviorPolicy.DerivePolicy(history)
	targetPolicy.DerivePolicy(model, updates)
	// All store and access operations can take their index mod n + 1
	decisions := make(map[int]*Decision, 0)
	rewards := make(map[int]float64, 0)

	// Loop for each episode
	for i := 0; i < tb.GetEpisodes(); i++ {
		exportResults += "++++++++++++++++++++++++++\n"
		exportResults += fmt.Sprintf("Episode: %d\n", i + 1)
		exportResults += "++++++++++++++++++++++++++\n\n"
		if err = tb.log(exportResults, "results_" + strconv.Itoa(target) + ".txt"); err != nil {
			return err
		}
		exportResults = ""

		// Initialize and store S0 != terminal
		switch tb.GetYuma().GetMode() {
		case 0:
			if instances, ok := tb.GetYuma().GetEnvironment().GetInstances(target)[0]; ok {
				if len(instances) >= 1 {
					tb.GetYuma().GetEnvironment().DeleteInstance(target, 0)
				}
			}
		case 1:
			for i := 0; i < len(tb.GetYuma().GetSubprocesses()); i++ {
				if instances, ok := tb.GetYuma().GetEnvironment().GetInstances(target)[int(math.Exp2(float64(i)))]; ok {
					for len(instances) > 0 {
						tb.GetYuma().GetEnvironment().DeleteInstance(target, int(math.Exp2(float64(i))))
					}
				}
			}
		}

		state, path := tb.initializeState()
		// Choose an action A0 arbitrarily as a function of S0; Store A0
		if len(tb.GetTree()) > 0 {
			decisions[0] = NewDecision(state, tb.GetTree()[0])
		} else {
			c, _ := randutil.WeightedChoice(behaviorPolicy.GetSuggestions(state))
			action := c.Item.(int)
			decisions[0] = NewDecision(state, action)
		}

		// T = infinity
		terminal := math.MaxInt64
		rewind := false

		// Loop for t = 0, 1, 2, ...
		t := 0
		tau := t + 1 - tb.GetN()

		// Loop until tau = T - 1
		for tau < terminal - 1 {
			if t < terminal {
				if history.At(decisions[t].GetState(), target) == 0 {
					decisions[t].SetAction(target)
				}

				if history.At(decisions[t].GetState(), target) >= 1 && memory.At(decisions[t].GetState(), target) == -1.0 {
					decisions[t].SetAction(target)
				}

				f := mat.Formatted(history, mat.Prefix("          "), mat.Squeeze())
				exportResults += fmt.Sprintf("\nHistory = %v\n\n\n", f)
				f = mat.Formatted(updates, mat.Prefix("          "), mat.Squeeze())
				exportResults += fmt.Sprintf("\nUpdates = %v\n\n\n", f)
				f = mat.Formatted(model, mat.Prefix("        "), mat.Squeeze())
				exportResults += fmt.Sprintf("\nModel = %v\n\n\n", f)
				exportResults += fmt.Sprintf("T: %d\n", t)
				exportResults += fmt.Sprintf("S_%d: %d\n", t, decisions[t].GetState())
				exportResults += fmt.Sprintf("Action Weights: %v\n", behaviorPolicy.GetSuggestions(decisions[t].GetState()))
				exportResults += fmt.Sprintf("A_%d: %d\n", t, decisions[t].GetAction())
				exportResults += fmt.Sprintf("Path: ")
				exportResults += fmt.Sprintln(path)
				exportResults += fmt.Sprintln()
				if err = tb.log(exportResults, "results_" + strconv.Itoa(target) + ".txt"); err != nil {
					return err
				}
				exportResults = ""
				
				// Take action At; observe and store the next reward and state as Rt+1, St+1
				if (history.At(decisions[t].GetState(), decisions[t].GetAction()) == 0.0) || (time.Since(time.Unix(0, int64(timestamps.At(decisions[t].GetState(), decisions[t].GetAction())))) >= tb.GetTheta()) {
					err, success, reward, successor = tb.GetYuma().GetEnvironment().TakeAction(target, decisions[t].GetState(), decisions[t].GetAction(), path, success)
					if err != nil {
						return err
					} else if success {
						path = append(path, tb.GetYuma().GetConfigurations()[decisions[t].GetAction()])
					}
					decisions[t + 1] = NewDecision(successor, 0)
					rewards[t + 1] = reward
					memory.Set(decisions[t].GetState(), decisions[t].GetAction(), rewards[t + 1])
					timestamps.Set(decisions[t].GetState(), decisions[t].GetAction(), float64(time.Now().UnixNano()))
					if rewards[t + 1] == (float64(len(tb.GetYuma().GetSubprocesses())) + 1.0) * -1.0 {
						tau = t
					}
				} else {
					if memory.At(decisions[t].GetState(), decisions[t].GetAction()) == -1.0 {
						//time.Sleep(1 * time.Minute)
						//time.Sleep(10 * time.Second)
						success = false
						reward = memory.At(decisions[t].GetState(), decisions[t].GetAction())
						rewards[t + 1] = reward
						successor = decisions[t].GetState() | decisions[t].GetAction()
						decisions[t + 1] = NewDecision(successor, 0)
						path = append(path, tb.GetYuma().GetConfigurations()[decisions[t].GetAction()])
					} else if memory.At(decisions[t].GetState(), decisions[t].GetAction()) == (float64(len(tb.GetYuma().GetSubprocesses())) + 1.0) * -1.0 {
						//time.Sleep(1 * time.Minute)
						//time.Sleep(10 * time.Second)
						success = false
						reward = memory.At(decisions[t].GetState(), decisions[t].GetAction())
						rewards[t + 1] = reward
						successor = decisions[t].GetState()
						decisions[t + 1] = NewDecision(successor, 0)
						tau = t
					}
				}
				history.Set(decisions[t].GetState(), decisions[t].GetAction(), history.At(decisions[t].GetState(), decisions[t].GetAction()) + 1)

				exportResults += "R_" + fmt.Sprintf("%d", t + 1) + ": " + fmt.Sprintf("%f", reward) + "\n"
				if err = tb.log(exportResults, "results_" + strconv.Itoa(target) + ".txt"); err != nil {
					return err
				}
				exportResults = ""

				if decisions[t].GetState() == decisions[t + 1].GetState() {
					rewind = true
				}

				exportResults += fmt.Sprintf("S_%d: %d\n", t + 1, decisions[t + 1].GetState())
				if err = tb.log(exportResults, "results_" + strconv.Itoa(target) + ".txt"); err != nil {
					return err
				}
				exportResults = ""

				// If St+1 is terminal
				if (decisions[t + 1].GetState() & target > 0) || (tb.GetYuma().IsTerminal(decisions[t + 1].GetState())) {
					terminal = t + 1
				} else {
					// Choose an action At+1 arbitrarily as a function of St+1; Store At+1
					if len(tb.GetTree()) - 1 >= t + 1 {
						decisions[t + 1].SetAction(tb.GetTree()[t + 1])
					} else {
						if history.At(decisions[t + 1].GetState(), target) == 0 {
							decisions[t + 1].SetAction(target)
						} else {
							c, _ := randutil.WeightedChoice(behaviorPolicy.GetSuggestions(decisions[t + 1].GetState()))
							action := c.Item.(int)
							decisions[t + 1].SetAction(action)
						}
					}

					if reward == -1.0 {
						exportResults += fmt.Sprintf("Action Weights: %v\n", behaviorPolicy.GetSuggestions(decisions[t + 1].GetState()))
						exportResults += fmt.Sprintf("A_%d: %d\n", t + 1, decisions[t + 1].GetAction())
						if err = tb.log(exportResults, "results_" + strconv.Itoa(target) + ".txt"); err != nil {
							return err
						}
						exportResults = ""
					}
				}
			}
			var totalReturn float64

			// Tau is the time whose estimate is being updated
			tau = t + 1 - tb.GetN()
			if tau >= 0 {
				exportResults += fmt.Sprintf("\nN: %d\n", tb.GetN())
				exportResults += fmt.Sprintf("T: %d\n", t)
				exportResults += fmt.Sprintf("Terminal: %d\n", terminal)
				exportResults += fmt.Sprintf("\nTau: %d - State: %d - Action: %d\n\n", tau, decisions[tau].GetState(), decisions[tau].GetAction())
				if err = tb.log(exportResults, "results_" + strconv.Itoa(target) + ".txt"); err != nil {
					return err
				}
				exportResults = ""

				if (t + 1 >= terminal) || rewind {
					// G = RT
					if rewind {
						totalReturn = rewards[t + 1]
					} else {
						totalReturn = rewards[terminal]
					}
					
					exportResults += fmt.Sprintf("G = %f\n\n", totalReturn)
					if err = tb.log(exportResults, "results_" + strconv.Itoa(target) + ".txt"); err != nil {
						return err
					}
					exportResults = ""
				} else {
					// G = Rt+1 + gamma * SUM_over_a(pi(a|St+1) * Q(St+1, a))
					var totalActionValue float64

					exportResults += fmt.Sprintf("Action Weights: %v\n", behaviorPolicy.GetSuggestions(decisions[t + 1].GetState()))
					exportResults += fmt.Sprintf("G = %f + %f * (0 ", rewards[t + 1], tb.GetGamma())
					if err = tb.log(exportResults, "results_" + strconv.Itoa(target) + ".txt"); err != nil {
						return err
					}
					exportResults = ""

					for _, action := range tb.GetYuma().Actions(decisions[t + 1].GetState()) {
						weight := float64(targetPolicy.GetWeight(decisions[t + 1].GetState(), action)) / 10.0
						totalActionValue += float64(weight) * model.At(decisions[t + 1].GetState(), action)

						exportResults += fmt.Sprintf("+ %f * %f", weight, model.At(decisions[t + 1].GetState(), action))
						if err = tb.log(exportResults, "results_" + strconv.Itoa(target) + ".txt"); err != nil {
							return err
						}
						exportResults = ""
					}
					totalReturn = rewards[t + 1] + tb.GetGamma() * totalActionValue

					exportResults += fmt.Sprintf(") = %f\n\n", totalReturn)
					if err = tb.log(exportResults, "results_" + strconv.Itoa(target) + ".txt"); err != nil {
						return err
					}
					exportResults = ""
				}
				// Loop for k = min(t, T - 1) down through tau + 1
				k := int(math.Min(float64(t), float64(terminal - 1)))
				for i := k; i >= tau + 1; i-- {
					exportResults += fmt.Sprintf("S_%d: %d\n", i, decisions[i].GetState())
					if err = tb.log(exportResults, "results_" + strconv.Itoa(target) + ".txt"); err != nil {
						return err
					}
					exportResults = ""

					// G = Rk + gamma * SUM_over_a!=Ak(pi(a|Sk) * Q(Sk, a) + gamma * pi(Ak|Sk) * G)
					var totalActionValue float64
					if c, _ := randutil.WeightedChoice(targetPolicy.GetSuggestions(decisions[i].GetState())); c.Item == nil {
						targetPolicy.SetWeight(decisions[i].GetState(), decisions[i].GetAction(), 10)
					}

					exportResults += fmt.Sprintf("Target Weights: %v\n", targetPolicy.GetSuggestions(decisions[i].GetState()))
					exportResults += fmt.Sprintf("G = %f + %f * (0", rewards[i], tb.GetGamma())
					if err = tb.log(exportResults, "results_" + strconv.Itoa(target) + ".txt"); err != nil {
						return err
					}
					exportResults = ""

					returnCorrection := false
					for _, action := range tb.GetYuma().Actions(decisions[i].GetState()) {
						if action == decisions[i].GetAction() {
							continue
						}
						weight := float64(targetPolicy.GetWeight(decisions[i].GetState(), action)) / 10.0
						if (weight == 1.0) && (model.At(decisions[i].GetState(), action) < totalReturn) {
							weight = 0.0
							returnCorrection = true
						}
						totalActionValue += float64(weight) * model.At(decisions[i].GetState(), action)

						exportResults += fmt.Sprintf(" + %f * %f", weight, model.At(decisions[i].GetState(), action))
						if err = tb.log(exportResults, "results_" + strconv.Itoa(target) + ".txt"); err != nil {
							return err
						}
						exportResults = ""
					}

					exportResults += fmt.Sprintf(") + %f * %f * %f = ", tb.GetGamma(), float64(targetPolicy.GetWeight(decisions[i].GetState(), decisions[i].GetAction())) / 10.0, totalReturn)
					if err = tb.log(exportResults, "results_" + strconv.Itoa(target) + ".txt"); err != nil {
						return err
					}
					exportResults = ""

					if !returnCorrection {
						totalReturn = rewards[i] + tb.GetGamma() * totalActionValue + tb.GetGamma() * (float64(targetPolicy.GetWeight(decisions[i].GetState(), decisions[i].GetAction())) / 10.0) * totalReturn
					} else {
						totalReturn = rewards[i] + tb.GetGamma() * totalActionValue + tb.GetGamma() * 1.0 * totalReturn
					}

					exportResults += fmt.Sprintf("%f\n\n", totalReturn)
					if err = tb.log(exportResults, "results_" + strconv.Itoa(target) + ".txt"); err != nil {
						return err
					}
					exportResults = ""
				}
				// Q(Stau, Atau) = Q(Stau, Atau) + alpha * (G - Q(Stau, Atau))
				var qOld float64
				/* if totalReturn == (float64(len(tb.GetYuma().GetSubprocesses())) + 1.0) * -1.0 {
					qOld = totalReturn
				} else if history != nil {
					if history.At(decisions[tau].GetState(), decisions[tau].GetAction()) == 0 {
						qOld = 0.0
					} else {
						qOld = model.At(decisions[tau].GetState(), decisions[tau].GetAction())
					}
				} */
				if history != nil {
					if history.At(decisions[tau].GetState(), decisions[tau].GetAction()) == 0 {
						qOld = 0.0
					} else {
						qOld = model.At(decisions[tau].GetState(), decisions[tau].GetAction())
					}
				}

				exportResults += fmt.Sprintf("Q_Update = %f + %f * (%f - %f) = ", qOld, tb.GetAlpha(), totalReturn, qOld)
				if err = tb.log(exportResults, "results_" + strconv.Itoa(target) + ".txt"); err != nil {
					return err
				}
				exportResults = ""

				qUpdate := qOld + tb.GetAlpha() * (totalReturn - qOld)

				exportResults += fmt.Sprintf("%f\n\n", qUpdate)
				if err = tb.log(exportResults, "results_" + strconv.Itoa(target) + ".txt"); err != nil {
					return err
				}
				exportResults = ""

				model.Set(decisions[tau].GetState(), decisions[tau].GetAction(), qUpdate)
				updates.Set(decisions[tau].GetState(), decisions[tau].GetAction(), updates.At(decisions[tau].GetState(), decisions[tau].GetAction()) + 1.0)
				// If pi is being learned, then ensure that pi(.|Stau) is greedy wrt Q
				behaviorPolicy.DerivePolicy(model, updates)
				targetPolicy.DerivePolicy(model, updates)
			}
			if rewind {
				updates.Set(decisions[t].GetState(), decisions[t].GetAction(), updates.At(decisions[t].GetState(), decisions[t].GetAction()) + 1.0)
				behaviorPolicy.DerivePolicy(model, updates)
				targetPolicy.DerivePolicy(model, updates)

				c, _ := randutil.WeightedChoice(behaviorPolicy.GetSuggestions(decisions[t].GetState()))
				action := c.Item.(int)
				decisions[t].SetAction(action)

				exportResults += fmt.Sprintf("Action Weights: %v\n", behaviorPolicy.GetSuggestions(decisions[t].GetState()))
				exportResults += fmt.Sprintf("A_%d: %d\n", t, decisions[t].GetAction())
				if err = tb.log(exportResults, "results_" + strconv.Itoa(target) + ".txt"); err != nil {
					return err
				}
				exportResults = ""

				rewind = false
				t -= 1
			}
			t += 1
		}
	}

	exportResults += "++++++++++++++++++++++++++\n"
	exportResults += fmt.Sprintf("Completed Learning\n")
	exportResults += "++++++++++++++++++++++++++\n\n"
	if err = tb.log(exportResults, "results_" + strconv.Itoa(target) + ".txt"); err != nil {
		return err
	}

	f := mat.Formatted(memory, mat.Prefix("         "), mat.Squeeze())
	exportMemory := fmt.Sprintf("\nMemory = %v\n\n\n", f)
	if err = tb.log(exportMemory, "memory_" + strconv.Itoa(target) + ".txt"); err != nil {
		return err
	}

	return nil
}

func (tb *TreeBackup) Solve(target int, model, updates, history, memory *mat.Dense, timestamps *mat.Dense) []string {
	state := tb.GetYuma().GetStartState()
	solution := make([]string, 0)

	for state & target == 0 {
		action := tb.ArgMaxAction(model, state, tb.GetYuma().Actions(state))
		subprocess := tb.GetYuma().GetConfigurations()[action]
		solution = append(solution, subprocess)
		state = state | action
	}

	return solution
}

func (tb *TreeBackup) initializeState() (int, []string) {
	if tb.GetRoot() != -1 {
		return tb.GetRoot(), tb.GetPath()
	} else {
		return tb.GetYuma().GetStartState(), make([]string, 0)
	}
}

func (tb *TreeBackup) ArgMaxAction(q *mat.Dense, state int, actions []int) int {
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