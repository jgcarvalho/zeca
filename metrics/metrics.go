package metrics

import (
	"math"
)

func MatthewsCorr(cm [][]int) float64 {
	return -1.0
}

// Class Balance Accuracy
func CBA(cm [][]int) float64 {
	nr := len(cm)
	np := len(cm[0])

	cba := make([]float64, nr)

	for i := 0; i < nr; i++ {
		ci_, c_i := 0.0, 0.0
		for j := 0; j < nr; j++ {
			ci_ += float64(cm[i][j])
			c_i += float64(cm[j][i])
		}
		if nr != np {
			ci_ += float64(cm[i][nr])
		}
		cba[i] = float64(cm[i][i]) / math.Max(ci_, c_i)
	}

	total := 0.0
	for t := 0; t < nr; t++ {
		total += cba[t]
	}
	return total / float64(nr)
}
