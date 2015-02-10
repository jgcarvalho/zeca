package disteda

import (
	"bytes"
	"fmt"
	"math/rand"

	"github.com/jgcarvalho/zeca/ca"
	"github.com/jgcarvalho/zeca/metrics"
	"github.com/jgcarvalho/zeca/rules"
)

type Probabilities struct {
	PID        uint32
	Generation int
	Data       [][][][]float64
}

type Probs struct {
	probs   [][][][]float64
	rulePrm rules.Params
}

type Population struct {
	rule    []*rules.Rule
	fitness []float64
}

type Tournament struct {
	rule    []*rules.Rule
	fitness []float64
}

type Individual struct {
	PID        uint32
	Generation int
	Rule       *rules.Rule
	Fitness    float64
}

func Fitness(c *ca.CellAuto1D) float64 {
	c.Run()
	cm := c.ConfusionMatrix()
	cba := metrics.CBA(cm)
	// fmt.Println("CBA: ", cba)
	return cba
}

func NewProbs(prm rules.Params) *Probs {
	var p Probs
	p.rulePrm = prm
	st := rules.RuleStates(prm)
	p.probs = make([][][][]float64, len(st))
	for ln := range st {
		p.probs[ln] = make([][][]float64, len(st))
		for c := range st {
			p.probs[ln][c] = make([][]float64, len(st))
			for rn := range st {
				//p.probs[ln][c][rn] = prm.transitionStates[rand.Intn(len(prm.transitionStates))]
				p.probs[ln][c][rn] = make([]float64, len(prm.TransitionStates))
				for pv := range p.probs[ln][c][rn] {
					if c != 0 {
						p.probs[ln][c][rn][pv] = 1.0 / float64(len(p.probs[ln][c][rn]))
					} else {
						if pv == len(p.probs[ln][c][rn])-1 {
							p.probs[ln][c][rn][pv] = 1.0
						} else {
							p.probs[ln][c][rn][pv] = 0.0
						}

					}

				}
			}
		}
	}

	return &p
}

func (p *Probs) AdjustProbs(pop Population) {
	// n := len(pop.rule)

	var count [][][][]float64
	st := rules.RuleStates(p.rulePrm)
	count = make([][][][]float64, len(st))
	for ln := range st {
		count[ln] = make([][][]float64, len(st))
		for c := range st {
			count[ln][c] = make([][]float64, len(st))
			for rn := range st {
				//p.probs[ln][c][rn] = prm.transitionStates[rand.Intn(len(prm.transitionStates))]
				count[ln][c][rn] = make([]float64, len(p.rulePrm.TransitionStates))
			}
		}
	}
	for j := 0; j < len(pop.rule); j++ {
		n := len(pop.rule[0].Code)
		for ln := 0; ln < n; ln++ {
			for c := 0; c < n; c++ {
				for rn := 0; rn < n; rn++ {
					index := bytes.IndexByte(p.rulePrm.TransitionStates, pop.rule[j].Code[ln][c][rn])
					count[ln][c][rn][index] += 1
				}
			}
		}
	}

	n := len(pop.rule[0].Code)
	for ln := 0; ln < n; ln++ {
		for c := 0; c < n; c++ {
			for rn := 0; rn < n; rn++ {
				for i := 0; i < len(count[ln][c][rn]); i++ {
					p.probs[ln][c][rn][i] = count[ln][c][rn][i] / float64(len(pop.rule))
				}
				// fmt.Printf("Contagem %f %f %f %f \n", count[ln][c][rn][0], count[ln][c][rn][1], count[ln][c][rn][2], count[ln][c][rn][3])
				// fmt.Printf("len pop rule %d \n", len(pop.rule))
				fmt.Printf("Probabilidade %f %f %f %f \n", p.probs[ln][c][rn][0], p.probs[ln][c][rn][1], p.probs[ln][c][rn][2], p.probs[ln][c][rn][3])
			}
		}
	}

}

func (p *Probs) GenRule() *rules.Rule {
	var rule rules.Rule
	rule.Prm = p.rulePrm
	st := rules.RuleStates(rule.Prm)
	rule.Code = make([][][]byte, len(st))
	for ln := range st {
		rule.Code[ln] = make([][]byte, len(st))
		for c := range st {
			rule.Code[ln][c] = make([]byte, len(st))
			for rn := range st {
				randv := rand.Float64()
				for i, v := range p.probs[ln][c][rn] {
					// fmt.Println("valor", v)
					// fmt.Println("valor", randv)
					randv -= v
					if randv < 0.0 {
						rule.Code[ln][c][rn] = rule.Prm.TransitionStates[i]
						break
					}
				}
			}
		}

	}

	return &rule
}

func (t Tournament) Len() int {
	return len(t.rule)
}

func (t Tournament) Swap(i, j int) {
	t.rule[i], t.rule[j] = t.rule[j], t.rule[i]
	t.fitness[i], t.fitness[j] = t.fitness[j], t.fitness[i]
}

func (t Tournament) Less(i, j int) bool {
	return t.fitness[i] < t.fitness[j]
}
