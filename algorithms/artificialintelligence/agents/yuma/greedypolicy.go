package yuma

import (
	"math"
	"github.com/jmcvetta/randutil"
	"gonum.org/v1/gonum/mat"
)

type GreedyPolicy struct {
	rationalThinking RationalThinking
	suggestions map[int][]randutil.Choice
}

func NewGreedyPolicy() *GreedyPolicy {
	gp := new(GreedyPolicy)
	gp.suggestions  = make(map[int][]randutil.Choice, 0)

	return gp
}

func (gp *GreedyPolicy) GetRationalThinking() RationalThinking {
	return gp.rationalThinking
}

func (gp *GreedyPolicy) GetSuggestions() map[int][]randutil.Choice {
	return gp.suggestions
}

func (gp *GreedyPolicy) GetWeight(state int, action int) int {
	choices := gp.GetSuggestions()[state]
	for _, c := range choices {
		if c.Item == action {
			return c.Weight
		}
	}

	return 0
}

func (gp *GreedyPolicy) SetRationalThinking(rationalThinking RationalThinking) {
	gp.rationalThinking = rationalThinking
}

func (gp *GreedyPolicy) DerivePolicy(q *mat.Dense) {
	for i := 0; i < int(math.Exp2(float64(len(gp.GetRationalThinking().GetYuma().GetSubprocesses())))); i++ {
		maxAction := gp.GetRationalThinking().ArgMaxAction(q, i)
		choices := make([]randutil.Choice, 0)
		for _, action := range gp.GetRationalThinking().GetYuma().Actions(i) {
			c := randutil.Choice{}
			c.Item = action
			if action == maxAction {
				c.Weight = 10
			} else {
				c.Weight = 0
			}
			choices = append(choices, c)
		}
		gp.suggestions[i] = choices
	}
}