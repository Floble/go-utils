package yuma

import (
	"github.com/jmcvetta/randutil"
	"gonum.org/v1/gonum/mat"
)

type Policy interface {
	GetSuggestions() map[int][]randutil.Choice
	GetWeight(int, int) int
	DerivePolicy(*mat.Dense)
}