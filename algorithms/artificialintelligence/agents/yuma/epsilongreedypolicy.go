package yuma

import (
	"math"
	"github.com/jmcvetta/randutil"
	"gonum.org/v1/gonum/mat"
)

type EpsilonGreedyPolicy struct {
	rationalThinking RationalThinking
	suggestions map[int][]randutil.Choice
	epsilon float64
}

func NewEpsilonGreedyPolicy(epsilon float64) *EpsilonGreedyPolicy {
	egp := new(EpsilonGreedyPolicy)
	egp.suggestions  = make(map[int][]randutil.Choice, 0)
	egp.epsilon = epsilon

	return egp
}

func (egp *EpsilonGreedyPolicy) GetRationalThinking() RationalThinking {
	return egp.rationalThinking
}

func (egp *EpsilonGreedyPolicy) GetSuggestions() map[int][]randutil.Choice {
	return egp.suggestions
}

func (egp *EpsilonGreedyPolicy) GetEpsilon() float64 {
	return egp.epsilon
}

func (egp *EpsilonGreedyPolicy) SetRationalThinking(rationalThinking RationalThinking) {
	egp.rationalThinking = rationalThinking
}

func (egp *EpsilonGreedyPolicy) DerivePolicy(q *mat.Dense) {
	for i := 0; i < int(math.Exp2(float64(len(egp.GetRationalThinking().GetYuma().GetSubprocesses())))); i++ {
		maxAction := egp.GetRationalThinking().ArgMaxAction(q, i)
		choices := make([]randutil.Choice, 0)
		for _, action := range egp.GetRationalThinking().GetYuma().Actions(i) {
			c := randutil.Choice{}
			c.Item = action
			if action == maxAction {
				c.Weight = int((1.0 - egp.GetEpsilon()) * 10.0)
			} else {
				c.Weight = int(egp.GetEpsilon() * 10.0)
			}
			choices = append(choices, c)
		}
		egp.suggestions[i] = choices
	}
}