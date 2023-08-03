package yuma

import (
	"math"
	"sync"
	"github.com/jmcvetta/randutil"
	"gonum.org/v1/gonum/mat"
)

type NaraExpansionPolicy struct {
	rationalThinking RationalThinking
	suggestions sync.Map
}

func NewNaraExpansionPolicy() *NaraExpansionPolicy {
	nep := new(NaraExpansionPolicy)

	return nep
}

func (nep *NaraExpansionPolicy) GetRationalThinking() RationalThinking {
	return nep.rationalThinking
}

func (nep *NaraExpansionPolicy) GetSuggestions(target int) []randutil.Choice {
	if suggestions, ok := nep.suggestions.Load(target); ok {
		return suggestions.([]randutil.Choice)
	} else {
		return nil
	}
}

func (nep *NaraExpansionPolicy) SetSuggestions(target int, suggestions []randutil.Choice) {
	nep.suggestions.Store(target, suggestions)
}

func (nep *NaraExpansionPolicy) GetWeight(state int, action int) int {
	choices := nep.GetSuggestions(state)
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
		nep.SetSuggestions(i, choices)
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
		nep.SetSuggestions(i, choices)
	}
}