package yuma

import (
	"github.com/jmcvetta/randutil"
	"gonum.org/v1/gonum/mat"
)

type Policy interface {
	GetSuggestions(int) []randutil.Choice
	GetWeight(int, int) int
	SetWeight(int, int, int)
	DerivePolicy(*mat.Dense, *mat.Dense)
}