package yuma

import (
	"math"
	"github.com/jmcvetta/randutil"
	"gonum.org/v1/gonum/mat"
)

type RandomPolicy struct {
	rationalThinking RationalThinking
	suggestions map[int][]randutil.Choice
}

func NewRandomPolicy() *RandomPolicy {
	rp := new(RandomPolicy)
	rp.suggestions  = make(map[int][]randutil.Choice, 0)

	return rp
}

func (rp *RandomPolicy) GetRationalThinking() RationalThinking {
	return rp.rationalThinking
}

func (rp *RandomPolicy) GetSuggestions() map[int][]randutil.Choice {
	return rp.suggestions
}

func (rp *RandomPolicy) GetWeight(state int, action int) int {
	choices := rp.GetSuggestions()[state]
	for _, c := range choices {
		if c.Item == action {
			return c.Weight
		}
	}

	return 0
}

func (rp *RandomPolicy) SetWeight(state int, action int, weight int) {
	choices := rp.GetSuggestions()[state]
	for _, c := range choices {
		if c.Item == action {
			c.Weight = weight
		}
	}
}

func (rp *RandomPolicy) SetRationalThinking(rationalThinking RationalThinking) {
	rp.rationalThinking = rationalThinking
}

func (rp *RandomPolicy) DerivePolicy(q *mat.Dense, updates *mat.Dense) {
	for i := 0; i < int(math.Exp2(float64(len(rp.GetRationalThinking().GetYuma().GetSubprocesses())))); i++ {
		choices := make([]randutil.Choice, 0)
		for _, action := range rp.GetRationalThinking().GetYuma().Actions(i) {
			c := randutil.Choice{}
			c.Item = action
			c.Weight = 5
			choices = append(choices, c)
		}
		rp.suggestions[i] = choices
	}
}