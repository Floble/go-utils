package yuma

import (
	"gonum.org/v1/gonum/mat"
)

type RationalThinking interface {
	GetYuma() *Yuma
	Learn(int) error
	Solve(int) []string
	ArgMaxAction(*mat.Dense, int) int
	log(string, string) error
}