package yuma

import (
	"math"
	"github.com/jmcvetta/randutil"
	"gonum.org/v1/gonum/mat"
)

type NaraExpansionPolicy struct {
	rationalThinking RationalThinking
	suggestions map[int][]randutil.Choice
}

func NewNaraExpansionPolicy() *NaraExpansionPolicy {
	nep := new(NaraExpansionPolicy)
	nep.suggestions  = make(map[int][]randutil.Choice, 0)

	return nep
}

func (nep *NaraExpansionPolicy) GetRationalThinking() RationalThinking {
	return nep.rationalThinking
}

func (nep *NaraExpansionPolicy) GetSuggestions() map[int][]randutil.Choice {
	return nep.suggestions
}

func (nep *NaraExpansionPolicy) GetWeight(state int, action int) int {
	choices := nep.GetSuggestions()[state]
	for _, c := range choices {
		if c.Item == action {
			return c.Weight
		}
	}

	return 0
}

func (nep *NaraExpansionPolicy) SetWeight(state int, action int, weight int) {
	for i := 0; i < int(math.Exp2(float64(len(nep.GetRationalThinking().GetYuma().GetSubprocesses())))); i++ {
		choices := make([]randutil.Choice, 0)
		for _, a := range nep.GetRationalThinking().GetYuma().Actions(i) {
			c := randutil.Choice{}
			c.Item = a
			if i == state && a == action {
				c.Weight = 10
			} else {
				c.Weight = 0
			}
			choices = append(choices, c)
		}
		nep.suggestions[i] = choices
	}
}

func (nep *NaraExpansionPolicy) SetRationalThinking(rationalThinking RationalThinking) {
	nep.rationalThinking = rationalThinking
}

func (nep *NaraExpansionPolicy) DerivePolicy(q *mat.Dense, updates *mat.Dense) {
	for i := 0; i < int(math.Exp2(float64(len(nep.GetRationalThinking().GetYuma().GetSubprocesses())))); i++ {
		choices := make([]randutil.Choice, 0)
		for _, action := range nep.GetRationalThinking().GetYuma().Actions(i) {
			c := randutil.Choice{}
			c.Item = action
			if updates.At(i, action) == 0.0 {
				c.Weight = 10
			} else {
				c.Weight = 0
			}
			choices = append(choices, c)
		}
		nep.suggestions[i] = choices
	}
}