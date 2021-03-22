package machinelearning

import (
	helper "go-utils/helper"
	"math"
	"math/rand"
	"time"
	"gonum.org/v1/gonum/mat"
)

type KMeans struct {
	clusters, maxSteps int
	centroids []*mat.Dense
	assignments map[*mat.Dense][]*mat.Dense
}

func NewKMeans(clusters, maxSteps int) *KMeans {
	km := new(KMeans)
	km.clusters = clusters
	km.maxSteps = maxSteps

	return km
}

func (km *KMeans) init(d int) {
	randSource := rand.NewSource(time.Now().UnixNano())
	randGen := rand.New(randSource)

	for i := 0; i < km.clusters; i++ {
		km.centroids = append(km.centroids, mat.NewDense(d, 1, nil))

		for j := 0; j < d; j++ {
			km.centroids[i].Set(j, 0, randGen.Float64())
		}
	}
}

func (km *KMeans) setAssignments(data []*mat.Dense) {
	km.assignments = make(map[*mat.Dense][]*mat.Dense, 0)
	applySquared := func(_, _ int, n float64) float64 { return n * n }

	for _, centroid := range km.centroids {
		km.assignments[centroid] = make([]*mat.Dense, 0)
	}

	for _, date := range data {
		min := math.MaxFloat64
		var assignment *mat.Dense

		for _, centroid := range km.centroids {
			tmp := new(mat.Dense)
			tmp.Sub(date, centroid)
			tmp.Apply(applySquared, tmp)
			tmp = helper.SumAlongColumn(tmp)
			distance := tmp.At(0, 0)
			distance = math.Sqrt(distance)
			
			if distance < min {
				min = distance
				assignment = centroid
			}
		}

		km.assignments[assignment] = append(km.assignments[assignment], date)
	}
}

func (km *KMeans) setCentroids() {
	km.centroids = make([]*mat.Dense, 0)

	for centroid, assignment := range km.assignments {
		if len(assignment) == 0 {
			km.centroids = append(km.centroids, centroid)
			continue
		}

		r, c := assignment[0].Dims()
		tmp := mat.NewDense(r, c, nil)
		tmp.Zero()

		for _, date := range assignment {
			tmp.Add(tmp, date)
		}

		tmp.Scale(1.0 / float64(len(assignment)), tmp)
		km.centroids = append(km.centroids, tmp)
	}
}

func (km *KMeans) Run(x []*mat.Dense) (map[*mat.Dense][]*mat.Dense, float64) {
	d, _ := x[0].Dims()
	km.init(d)

	for i := 0; i < km.maxSteps; i++ {
		km.setAssignments(x)
		km.setCentroids()
	}

	return km.assignments, km.loss(x)
}

func (km *KMeans) loss(x []*mat.Dense) float64 {
	applySquared := func(_, _ int, n float64) float64 { return n * n }
	loss := 0.0

	for centroid, assignment := range km.assignments {
		for _, date := range assignment {
			tmp := new(mat.Dense)
			tmp.Sub(date, centroid)
			tmp.Apply(applySquared, tmp)
			tmp = helper.SumAlongColumn(tmp)
			distance := tmp.At(0, 0)
			loss += distance
		}
	}

	return math.Sqrt(loss)
}