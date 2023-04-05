package yuma

import (
	"math"
	"github.com/jmcvetta/randutil"
	"gonum.org/v1/gonum/mat"
)

type NaraEpsilonGreedyPolicy struct {
	rationalThinking RationalThinking
	suggestions map[int][]randutil.Choice
	epsilon float64
}

func NewNaraEpsilonGreedyPolicy(epsilon float64) *NaraEpsilonGreedyPolicy {
	negp := new(NaraEpsilonGreedyPolicy)
	negp.suggestions  = make(map[int][]randutil.Choice, 0)
	negp.epsilon = epsilon

	return negp
}

func (negp *NaraEpsilonGreedyPolicy) GetRationalThinking() RationalThinking {
	return negp.rationalThinking
}

func (negp *NaraEpsilonGreedyPolicy) GetSuggestions() map[int][]randutil.Choice {
	return negp.suggestions
}

func (negp *NaraEpsilonGreedyPolicy) GetWeight(state int, action int) int {
	choices := negp.GetSuggestions()[state]
	for _, c := range choices {
		if c.Item == action {
			return c.Weight
		}
	}

	return 0
}

func (negp *NaraEpsilonGreedyPolicy) SetWeight(state int, action int, weight int) {
	choices := negp.GetSuggestions()[state]
	for _, c := range choices {
		if c.Item == action {
			c.Weight = weight
		}
	}
}

func (negp *NaraEpsilonGreedyPolicy) GetEpsilon() float64 {
	return negp.epsilon
}

func (negp *NaraEpsilonGreedyPolicy) SetRationalThinking(rationalThinking RationalThinking) {
	negp.rationalThinking = rationalThinking
}

func (negp *NaraEpsilonGreedyPolicy) DerivePolicy(q *mat.Dense, updates *mat.Dense) {
	for i := 0; i < int(math.Exp2(float64(len(negp.GetRationalThinking().GetYuma().GetSubprocesses())))); i++ {
		actions := make([]int, 0)
		for _, action := range negp.GetRationalThinking().GetYuma().Actions(i) {
			if (q.At(i, action) <= (float64(len(negp.GetRationalThinking().GetYuma().GetSubprocesses())) + 1.0) * -1.0) && (updates.At(i, action) >= 1) {
			} else {
				actions = append(actions, action)
			}
		}
		if len(actions) == 0 {
			continue
		}

		maxAction := negp.GetRationalThinking().ArgMaxAction(q, i, actions)
		choices := make([]randutil.Choice, 0)
		for _, action := range negp.GetRationalThinking().GetYuma().Actions(i) {
			c := randutil.Choice{}
			c.Item = action
			if (q.At(i, action) <= (float64(len(negp.GetRationalThinking().GetYuma().GetSubprocesses())) + 1.0) * -1.0) && (updates.At(i, action) >= 1) {
				c.Weight = 0.0
			} else {
				if action == maxAction {
					c.Weight = int((1.0 - negp.GetEpsilon()) * 10.0)
				} else {
					c.Weight = int(negp.GetEpsilon() * 10.0)
				}
			}
			choices = append(choices, c)
		}
		negp.suggestions[i] = choices
	}
}

func (negp *NaraEpsilonGreedyPolicy) getActionWeight(state int, action int) int {
	actions := negp.GetSuggestions()[state]
	for _, a := range actions {
		if a.Item.(int) == action {
			return a.Weight
		}
	}

	return -1
}