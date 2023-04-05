package yuma

import (
	"math"
	"github.com/jmcvetta/randutil"
	"gonum.org/v1/gonum/mat"
)

type HistoricPolicy struct {
	rationalThinking RationalThinking
	suggestions map[int][]randutil.Choice
}

func NewHistoricPolicy() *HistoricPolicy {
	hp := new(HistoricPolicy)
	hp.suggestions  = make(map[int][]randutil.Choice, 0)

	return hp
}

func (hp *HistoricPolicy) GetRationalThinking() RationalThinking {
	return hp.rationalThinking
}

func (hp *HistoricPolicy) GetSuggestions() map[int][]randutil.Choice {
	return hp.suggestions
}

func (hp *HistoricPolicy) GetWeight(state int, action int) int {
	choices := hp.GetSuggestions()[state]
	for _, c := range choices {
		if c.Item == action {
			return c.Weight
		}
	}

	return 0
}

func (hp *HistoricPolicy) SetWeight(state int, action int, weight int) {
	choices := hp.GetSuggestions()[state]
	for _, c := range choices {
		if c.Item == action {
			c.Weight = weight
		}
	}
}

func (hgp *HistoricPolicy) SetRationalThinking(rationalThinking RationalThinking) {
	hgp.rationalThinking = rationalThinking
}

func (hp *HistoricPolicy) DerivePolicy(q *mat.Dense, history *mat.Dense) {
	for i := 0; i < int(math.Exp2(float64(len(hp.GetRationalThinking().GetYuma().GetSubprocesses())))); i++ {
		choices := make([]randutil.Choice, 0)
		historicSum := float64(len(hp.GetRationalThinking().GetYuma().GetSubprocesses()))
		for _, action := range hp.GetRationalThinking().GetYuma().Actions(i) {
			c := randutil.Choice{}
			c.Item = action
			if historicSum == 0.0 || history.At(i, action) == 0.0 {
				c.Weight = 10
			} else {
				if int((100 - (history.At(i, action) / historicSum) * 10.0) / 10.0) <= 0 {
					c. Weight = 1
				} else {
					c.Weight = int((100 - (history.At(i, action) / historicSum) * 10.0) / 10.0)
				}
			}
			choices = append(choices, c)
		}
		hp.suggestions[i] = choices
	}
}