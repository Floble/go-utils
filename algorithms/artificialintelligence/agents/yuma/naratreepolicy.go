package yuma

import (
	"math"
	"math/rand"
	"sync"
	"github.com/jmcvetta/randutil"
	"gonum.org/v1/gonum/mat"
)

type NaraTreePolicy struct {
	rationalThinking RationalThinking
	suggestions sync.Map
	epsilon float64
}

func NewNaraTreePolicy(epsilon float64) *NaraTreePolicy {
	ntp := new(NaraTreePolicy)
	ntp.epsilon = epsilon

	return ntp
}

func (ntp *NaraTreePolicy) GetRationalThinking() RationalThinking {
	return ntp.rationalThinking
}

func (ntp *NaraTreePolicy) GetSuggestions(target int) []randutil.Choice {
	if suggestions, ok := ntp.suggestions.Load(target); ok {
		return suggestions.([]randutil.Choice)
	} else {
		return nil
	}
}

func (ntp *NaraTreePolicy) SetSuggestions(target int, suggestions []randutil.Choice) {
	ntp.suggestions.Store(target, suggestions)
}

func (ntp *NaraTreePolicy) GetWeight(state int, action int) int {
	choices := ntp.GetSuggestions(state)
	for _, c := range choices {
		if c.Item == action {
			return c.Weight
		}
	}

	return 0
}

func (ntp *NaraTreePolicy) GetEpsilon() float64 {
	return ntp.epsilon
}

func (ntp *NaraTreePolicy) SetWeight(state int, action int, weight int) {
	for i := 0; i < int(math.Exp2(float64(len(ntp.GetRationalThinking().GetYuma().GetSubprocesses())))); i++ {
		choices := make([]randutil.Choice, 0)
		for _, a := range ntp.GetRationalThinking().GetYuma().Actions(i) {
			c := randutil.Choice{}
			c.Item = a
			if i == state && a == action {
				c.Weight = 10
			} else {
				c.Weight = 0
			}
			choices = append(choices, c)
		}
		ntp.SetSuggestions(i, choices)
	}
}

func (ntp *NaraTreePolicy) SetRationalThinking(rationalThinking RationalThinking) {
	ntp.rationalThinking = rationalThinking
}

func (ntp *NaraTreePolicy) DerivePolicy(q *mat.Dense, updates *mat.Dense) {
	for i := 0; i < int(math.Exp2(float64(len(ntp.GetRationalThinking().GetYuma().GetSubprocesses())))); i++ {
		actions := make([]int, 0)
		for _, action := range ntp.GetRationalThinking().GetYuma().Actions(i) {
			if (q.At(i, action) <= (float64(len(ntp.GetRationalThinking().GetYuma().GetSubprocesses())) + 1.0) * -1.0) && (updates.At(i, action) >= 1) {
			} else if updates.At(i, action) >= 1 {
				actions = append(actions, action)
			}
		}
		if len(actions) == 0 {
			continue
		}

		maxAction := argMaxAction(q, updates, i, actions)
		choices := make([]randutil.Choice, 0)
		for _, action := range ntp.GetRationalThinking().GetYuma().Actions(i) {
			c := randutil.Choice{}
			c.Item = action
			if action == maxAction {
				if q.At(i, action) <= (float64(len(ntp.GetRationalThinking().GetYuma().GetSubprocesses())) + 1.0) * -1.0 {
					c.Weight = 0
				} else {
					c.Weight = int((1.0 - ntp.GetEpsilon()) * 10.0)
					choices = append(choices, c)
				}
			} else {
				c.Weight = 0
			}
		}
		tmp := randutil.Choice{}
		tmp.Item = -1
		tmp.Weight = int(ntp.GetEpsilon() * 10.0)
		choices = append(choices, tmp)
		ntp.SetSuggestions(i, choices)
	}
}

/* func argMaxAction(q, updates *mat.Dense, state int, actions []int) int {
	maxQ := math.MaxFloat64 * -1.0
	tmp := rand.Intn(len(actions) - 0) + 0
	maxAction := actions[tmp]
	if q.At(state, maxAction) * updates.At(state, maxAction) > maxQ {
		maxQ = q.At(state, maxAction) * updates.At(state, maxAction)
	}
	
	for _, action := range actions {
		if q.At(state, action) * updates.At(state, action) > maxQ {
			maxQ = q.At(state, action) * updates.At(state, action)
			maxAction = action
		}
	}

	return maxAction
} */

/* func (ntp *NaraTreePolicy) DerivePolicy(q *mat.Dense, updates *mat.Dense) {
	for i := 0; i < int(math.Exp2(float64(len(ntp.GetRationalThinking().GetYuma().GetSubprocesses())))); i++ {
		actions := make([]int, 0)
		for _, action := range ntp.GetRationalThinking().GetYuma().Actions(i) {
			if (q.At(i, action) <= (float64(len(ntp.GetRationalThinking().GetYuma().GetSubprocesses())) + 1.0) * -1.0) && (updates.At(i, action) >= 1) {
			} else if updates.At(i, action) >= 1 {
				actions = append(actions, action)
			}
		}
		if len(actions) == 0 {
			continue
		}

		maxAction := argMaxAction(q, updates, i, actions)
		choices := make([]randutil.Choice, 0)
		for _, action := range ntp.GetRationalThinking().GetYuma().Actions(i) {
			c := randutil.Choice{}
			c.Item = action
			if action == maxAction {
				if q.At(i, action) <= (float64(len(ntp.GetRationalThinking().GetYuma().GetSubprocesses())) + 1.0) * -1.0 {
					c.Weight = 0
				} else {
					c.Weight = 10
				}
			} else {
				c.Weight = 0
			}
			choices = append(choices, c)
		}
		ntp.SetSuggestions(i, choices)
	}
} */

func argMaxAction(q, updates *mat.Dense, state int, actions []int) int {
	maxQ := math.MaxFloat64 * -1.0
	tmp := rand.Intn(len(actions) - 0) + 0
	maxAction := actions[tmp]
	//c := math.Sqrt(2)
	c := 2.0

	if (q.At(state, maxAction) + c * math.Sqrt(math.Log(sumOfRow(updates, state)) / updates.At(state, maxAction))) > maxQ {
		maxQ = q.At(state, maxAction) + c * math.Sqrt(math.Log(sumOfRow(updates, state)) / updates.At(state, maxAction))
	}
	
	for _, action := range actions {
		if (q.At(state, maxAction) + c * math.Sqrt(math.Log(sumOfRow(updates, state)) / updates.At(state, maxAction))) > maxQ {
			maxQ = q.At(state, maxAction) + c * math.Sqrt(math.Log(sumOfRow(updates, state)) / updates.At(state, maxAction))
			maxAction = action
		}
	}

	return maxAction
}

func sumOfRow (matrix *mat.Dense, r int) float64 {
	_, c := matrix.Dims()
	sum := 0.0
	for i := 0; i < c; i++ {
		sum += matrix.At(r, i)
	}

	return sum
}