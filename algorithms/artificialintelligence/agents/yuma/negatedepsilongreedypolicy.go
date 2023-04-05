package yuma

import (
	"math"
	"github.com/jmcvetta/randutil"
	"gonum.org/v1/gonum/mat"
)

type NegatedEpsilonGreedyPolicy struct {
	rationalThinking RationalThinking
	suggestions map[int][]randutil.Choice
	epsilon float64
}

func NewNegatedEpsilonGreedyPolicy(epsilon float64) *NegatedEpsilonGreedyPolicy {
	negp := new(NegatedEpsilonGreedyPolicy)
	negp.suggestions  = make(map[int][]randutil.Choice, 0)
	negp.epsilon = epsilon

	return negp
}

func (negp *NegatedEpsilonGreedyPolicy) GetRationalThinking() RationalThinking {
	return negp.rationalThinking
}

func (negp *NegatedEpsilonGreedyPolicy) GetSuggestions() map[int][]randutil.Choice {
	return negp.suggestions
}

func (negp *NegatedEpsilonGreedyPolicy) GetWeight(state int, action int) int {
	choices := negp.GetSuggestions()[state]
	for _, c := range choices {
		if c.Item == action {
			return c.Weight
		}
	}

	return 0
}

func (negp *NegatedEpsilonGreedyPolicy) SetWeight(state int, action int, weight int) {
	choices := negp.GetSuggestions()[state]
	for _, c := range choices {
		if c.Item == action {
			c.Weight = weight
		}
	}
}

func (negp *NegatedEpsilonGreedyPolicy) GetEpsilon() float64 {
	return negp.epsilon
}

func (negp *NegatedEpsilonGreedyPolicy) SetRationalThinking(rationalThinking RationalThinking) {
	negp.rationalThinking = rationalThinking
}

func (negp *NegatedEpsilonGreedyPolicy) DerivePolicy(q *mat.Dense, updates *mat.Dense) {
	for i := 0; i < int(math.Exp2(float64(len(negp.GetRationalThinking().GetYuma().GetSubprocesses())))); i++ {
		maxAction := negp.GetRationalThinking().ArgMaxAction(q, i, negp.GetRationalThinking().GetYuma().Actions(i))
		choices := make([]randutil.Choice, 0)
		for _, action := range negp.GetRationalThinking().GetYuma().Actions(i) {
			c := randutil.Choice{}
			c.Item = action
			if action == maxAction {
				c.Weight = int(negp.GetEpsilon() * 10.0)
			} else {
				if q.At(i, action) == q.At(i, maxAction) {
					c.Weight = 10
				} else {
					c.Weight = int((1.0 - negp.GetEpsilon()) * 10.0)
				}
			}
			choices = append(choices, c)
		}
		negp.suggestions[i] = choices
	}
}