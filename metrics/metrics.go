package metrics

import (
	"math"
)

func MCC(cm [][]int) float64 {
	// mcc := 0.0

	return -1.0
}

func Q3(cm [][]int) float64 {
	n_trues := 0
	n_falses := 0
	for i := range cm {
		for j := range cm[i] {
			if i == j {
				n_trues += cm[i][j]
			} else {
				n_falses += cm[i][j]
			}
		}
	}
	return float64(n_trues) / float64(n_trues+n_falses)
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
