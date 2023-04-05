package yuma

import (
	"gonum.org/v1/gonum/mat"
)

type RationalThinking interface {
	GetYuma() *Yuma
	Learn(int, *mat.Dense, *mat.Dense, *mat.Dense, *mat.Dense) error
	Solve(int, *mat.Dense, *mat.Dense, *mat.Dense, *mat.Dense) []string
	ArgMaxAction(*mat.Dense, int, []int) int
	log(string, string) error
}