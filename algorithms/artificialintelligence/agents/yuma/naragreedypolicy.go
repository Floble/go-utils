package yuma

import (
	"math"
	"github.com/jmcvetta/randutil"
	"gonum.org/v1/gonum/mat"
)

type NaraGreedyPolicy struct {
	rationalThinking RationalThinking
	suggestions map[int][]randutil.Choice
}

func NewNaraGreedyPolicy() *NaraGreedyPolicy {
	ngp := new(NaraGreedyPolicy)
	ngp.suggestions  = make(map[int][]randutil.Choice, 0)

	return ngp
}

func (ngp *NaraGreedyPolicy) GetRationalThinking() RationalThinking {
	return ngp.rationalThinking
}

func (ngp *NaraGreedyPolicy) GetSuggestions() map[int][]randutil.Choice {
	return ngp.suggestions
}

func (ngp *NaraGreedyPolicy) GetWeight(state int, action int) int {
	choices := ngp.GetSuggestions()[state]
	for _, c := range choices {
		if c.Item == action {
			return c.Weight
		}
	}

	return 0
}

func (ngp *NaraGreedyPolicy) SetWeight(state int, action int, weight int) {
	for i := 0; i < int(math.Exp2(float64(len(ngp.GetRationalThinking().GetYuma().GetSubprocesses())))); i++ {
		choices := make([]randutil.Choice, 0)
		for _, a := range ngp.GetRationalThinking().GetYuma().Actions(i) {
			c := randutil.Choice{}
			c.Item = a
			if i == state && a == action {
				c.Weight = 10
			} else {
				c.Weight = 0
			}
			choices = append(choices, c)
		}
		ngp.suggestions[i] = choices
	}
}

func (ngp *NaraGreedyPolicy) SetRationalThinking(rationalThinking RationalThinking) {
	ngp.rationalThinking = rationalThinking
}

func (ngp *NaraGreedyPolicy) DerivePolicy(q *mat.Dense, updates *mat.Dense) {
	for i := 0; i < int(math.Exp2(float64(len(ngp.GetRationalThinking().GetYuma().GetSubprocesses())))); i++ {
		actions := make([]int, 0)
		for _, action := range ngp.GetRationalThinking().GetYuma().Actions(i) {
			if (q.At(i, action) <= (float64(len(ngp.GetRationalThinking().GetYuma().GetSubprocesses())) + 1.0) * -1.0) && (updates.At(i, action) >= 1) {
			} else {
				actions = append(actions, action)
			}
		}
		if len(actions) == 0 {
			continue
		}

		maxAction := ngp.GetRationalThinking().ArgMaxAction(q, i, actions)
		choices := make([]randutil.Choice, 0)
		for _, action := range ngp.GetRationalThinking().GetYuma().Actions(i) {
			c := randutil.Choice{}
			c.Item = action
			if action == maxAction {
				if q.At(i, action) <= (float64(len(ngp.GetRationalThinking().GetYuma().GetSubprocesses())) + 1.0) * -1.0 {
					c.Weight = 0
				} else {
					c.Weight = 10
				}
			} else {
				c.Weight = 0
			}
			choices = append(choices, c)
		}
		ngp.suggestions[i] = choices
	}
}