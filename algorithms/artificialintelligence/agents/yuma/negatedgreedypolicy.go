package yuma

import (
	"math"
	"github.com/jmcvetta/randutil"
	"gonum.org/v1/gonum/mat"
)

type NegatedGreedyPolicy struct {
	rationalThinking RationalThinking
	suggestions map[int][]randutil.Choice
}

func NewNegatedGreedyPolicy() *NegatedGreedyPolicy {
	ngp := new(NegatedGreedyPolicy)
	ngp.suggestions  = make(map[int][]randutil.Choice, 0)

	return ngp
}

func (ngp *NegatedGreedyPolicy) GetRationalThinking() RationalThinking {
	return ngp.rationalThinking
}

func (ngp *NegatedGreedyPolicy) GetSuggestions() map[int][]randutil.Choice {
	return ngp.suggestions
}

func (ngp *NegatedGreedyPolicy) GetWeight(state int, action int) int {
	choices := ngp.GetSuggestions()[state]
	for _, c := range choices {
		if c.Item == action {
			return c.Weight
		}
	}

	return 0
}

func (ngp *NegatedGreedyPolicy) SetWeight(state int, action int, weight int) {
	choices := ngp.GetSuggestions()[state]
	for _, c := range choices {
		if c.Item == action {
			c.Weight = weight
		}
	}
}

func (ngp *NegatedGreedyPolicy) SetRationalThinking(rationalThinking RationalThinking) {
	ngp.rationalThinking = rationalThinking
}

func (ngp *NegatedGreedyPolicy) DerivePolicy(q *mat.Dense, updates *mat.Dense) {
	for i := 0; i < int(math.Exp2(float64(len(ngp.GetRationalThinking().GetYuma().GetSubprocesses())))); i++ {
		maxAction := ngp.GetRationalThinking().ArgMaxAction(q, i, ngp.GetRationalThinking().GetYuma().Actions(i))
		choices := make([]randutil.Choice, 0)
		for _, action := range ngp.GetRationalThinking().GetYuma().Actions(i) {
			c := randutil.Choice{}
			c.Item = action
			if action == maxAction {
				c.Weight = 0
			} else {
				c.Weight = 10
			}
			choices = append(choices, c)
		}
		ngp.suggestions[i] = choices
	}
}