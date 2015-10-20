package cga

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"sort"
	"sync"
	"time"

	"bitbucket.org/jgcarvalho/zeca/ca"
	"bitbucket.org/jgcarvalho/zeca/db"
	"bitbucket.org/jgcarvalho/zeca/metrics"
	"bitbucket.org/jgcarvalho/zeca/rules"
)

type Probs struct {
	probs   [][][][]float64
	rulePrm rules.Params
}

// func Run(selby string, fnrulein string, fnruleout string, fnprobout string, gen int, pop int, steps int, ca *ca.CellAuto1D, prm rules.Params) error {
func Run(conf Config) error {
	rand.Seed(time.Now().UTC().UnixNano())

	fmt.Println("Loading proteins...")
	id, start, end, err := db.GetProteins(conf.DB)
	if err != nil {
		panic(err)
	}

	fmt.Println("Initializing probabilities...")
	rule, _ := rules.Create(conf.CA.InitStates, conf.CA.TransStates, conf.CA.HasJoker, conf.CA.R)
	probs := NewProbs(rule.Prm)
	fmt.Println(probs)

	var calist []*ca.CellAuto1D

	calist = make([]*ca.CellAuto1D, conf.CGA.Selection)

	for i := 0; i < len(calist); i++ {
		r := probs.GenRule()
		calist[i], _ = ca.Create1D(id, start, end, r, conf.CA.Steps, conf.CA.Consensus)
	}

	var wg1 sync.WaitGroup
	for i := 0; i < conf.CGA.Generations; i++ {
		fmt.Println("Generation", i)
		fmt.Println("Adjusting probabilities...")
		probs.AdjustByRanking(calist, 1.0/float64(conf.CGA.Population))
		fmt.Println("OK")

		if probs.Converged() {
			fmt.Println("Probabilities converged\nDONE")
			break
		}

		for i := 0; i < len(calist); i++ {
			wg1.Add(1)
			go func(i int) {
				defer wg1.Done()
				calist[i].SetRule(probs.GenRule())
			}(i)
			// calist[i], _ = ca.Create1D(id, start, end, r, conf.CA.Steps)
		}
		fmt.Printf("Waiting ")
		wg1.Wait()
		fmt.Println("OK")
	}

	err = ioutil.WriteFile(conf.CGA.OutputProbs, []byte(probs.String()), 0644)
	if err != nil {
		fmt.Println("Erro gravar as probabilidades")
		fmt.Println(probs)
	}

	return nil
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

func Fitness(c *ca.CellAuto1D) float64 {
	c.Run()
	cm := c.ConfusionMatrix()
	cba := metrics.CBA(cm)
	fmt.Println("CBA: ", cba)
	// p, n, unk := 0, 0, 0
	// for i, v := range c.Expected {
	// 	if c.Rule.Prm.Hasjoker == false || (c.Rule.Prm.Hasjoker == true && v != c.Rule.Prm.TransitionStates[len(c.Rule.Prm.TransitionStates)-1]) {
	// 		if c.Expected[i] == c.End[i] {
	// 			p += 1
	// 		} else {
	// 			n += 1
	// 		}
	// 	} else {
	// 		unk += 1
	// 	}
	// }
	// fmt.Println("P", p, "N", n, "Unk", unk, "Total", p+n+unk)
	// return float64(p) / float64(p+n)
	return cba
}

func (p *Probs) String() string {
	codes := append(p.rulePrm.StrStartStates, p.rulePrm.StrTransitionStates...)
	var toprint string
	toprint += fmt.Sprintf("[l][c][r] ->")
	for _, v := range p.rulePrm.StrTransitionStates {
		toprint += fmt.Sprintf(" %s", v)
	}
	toprint += fmt.Sprintln()
	for c := 0; c < len(p.probs); c++ {
		for ln := 0; ln < len(p.probs); ln++ {
			for rn := 0; rn < len(p.probs); rn++ {
				toprint += fmt.Sprintf("[%s][%s][%s] ->", codes[ln], codes[c], codes[rn])
				for _, v := range p.probs[ln][c][rn] {
					toprint += fmt.Sprintf(" %.4f", v)
				}
				toprint += fmt.Sprintln()

			}
		}
	}
	return toprint
}

type Ranking struct {
	c       []*ca.CellAuto1D
	fitness []float64
}

func (r Ranking) Len() int {
	return len(r.c)
}

func (r Ranking) Swap(i, j int) {
	r.c[i], r.c[j] = r.c[j], r.c[i]
	r.fitness[i], r.fitness[j] = r.fitness[j], r.fitness[i]
}

func (r Ranking) Less(i, j int) bool {
	return r.fitness[i] < r.fitness[j]
}

func (p *Probs) AdjustByRanking(c []*ca.CellAuto1D, weight float64) {
	var rank Ranking
	var wg2 sync.WaitGroup
	// var wg3 sync.WaitGroup
	rank.c = make([]*ca.CellAuto1D, len(c))
	copy(rank.c, c)
	rank.fitness = make([]float64, len(c))
	for i := 0; i < len(rank.c); i++ {
		wg2.Add(1)
		go func(r *Ranking, i int) {
			defer wg2.Done()
			r.fitness[i] = Fitness(r.c[i])
		}(&rank, i)
	}
	wg2.Wait()
	sort.Sort(rank)
	fmt.Println("Sorted Ranking", rank.fitness)

	s := (1 + len(rank.c)/2) * (len(rank.c) / 2) / 2
	unit := weight / float64(s)
	w := len(rank.c) / 2

	for i, j := 0, len(rank.c)-1; i < len(rank.c)/2 && j >= len(rank.c)/2; i, j = i+1, j-1 {
		// wg3.Add(1)
		// go func(i, j, w int) {
		// 	defer wg3.Done()
		fmt.Println("elementos", i, j, float64(w)*unit)
		p.duel(rank.c[i], rank.c[j], float64(w)*unit)
		// }(i, j, w)
		w--
	}
	// wg3.Wait()
}

// var mu = &sync.Mutex{}

func (p *Probs) duel(loser *ca.CellAuto1D, winner *ca.CellAuto1D, weight float64) {
	n := len(winner.Rule.Code)
	for ln := 0; ln < n; ln++ {
		for c := 0; c < n; c++ {
			for rn := 0; rn < n; rn++ {
				if winner.Rule.Code[ln][c][rn] != loser.Rule.Code[ln][c][rn] {
					windex := bytes.IndexByte(p.rulePrm.TransitionStates, winner.Rule.Code[ln][c][rn])
					lindex := bytes.IndexByte(p.rulePrm.TransitionStates, loser.Rule.Code[ln][c][rn])
					// mu.Lock()
					if p.probs[ln][c][rn][lindex] < weight {
						p.probs[ln][c][rn][windex] += p.probs[ln][c][rn][lindex]
						p.probs[ln][c][rn][lindex] -= p.probs[ln][c][rn][lindex]
					} else {
						p.probs[ln][c][rn][windex] += weight
						p.probs[ln][c][rn][lindex] -= weight
					}
					// mu.Unlock()
				}
			}
		}
	}
}

func (p *Probs) Converged() bool {
	n := len(p.probs)
	for ln := 0; ln < n; ln++ {
		for c := 0; c < n; c++ {
			for rn := 0; rn < n; rn++ {
				check := false
				for _, v := range p.probs[ln][c][rn] {
					if v > 0.999 {
						check = true
					}
				}
				if check == false {
					fmt.Println("Não convergiu")
					return false
				}
			}
		}
	}
	return true
}
