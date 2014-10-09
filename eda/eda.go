package eda

import (
	"bitbucket.org/zgcarvalho/zeca/ca"
	"bitbucket.org/zgcarvalho/zeca/metrics"
	"bitbucket.org/zgcarvalho/zeca/proteindb"
	"bitbucket.org/zgcarvalho/zeca/rules"
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"sort"
	"sync"
	"time"
)

type Probs struct {
	probs   [][][][]float64
	rulePrm rules.Params
}

type Population struct {
	rule    []*rules.Rule
	fitness []float64
}

func Run(conf Config) error {
	rand.Seed(time.Now().UTC().UnixNano())

	fmt.Println("Loading proteins...")
	proteins := proteindb.LoadProteinsFromMongo(conf.ProteinDB.Ip, conf.ProteinDB.Name, conf.ProteinDB.Collection)
	prot_id := "20primeiras"
	prot_seq := "#"
	prot_ss := "#"
	for i := 0; i < len(proteins); i++ {
		prot_seq += proteins[i].Chains[0].Seq_pdb + "#"
		prot_ss += proteins[i].Chains[0].Ss3_cons_all + "#"
	}

	fmt.Println("Initializing probabilities...")
	rule, _ := rules.Create(conf.CA.InitStates, conf.CA.TransStates, conf.CA.HasJoker, conf.CA.R)
	probs := NewProbs(rule.Prm)
	fmt.Println(probs)

	var pop Population

	pop.rule = make([]*rules.Rule, conf.EDA.Population)
	pop.fitness = make([]float64, conf.EDA.Population)

	var wg1 sync.WaitGroup
	// tmp := 0
	for i := 0; i < conf.EDA.Population; i++ {
		wg1.Add(1)
		go func(pop *Population, i int) {
			defer wg1.Done()
			pop.rule[i] = probs.GenRule()
			ca, _ := ca.Create1D(prot_id, prot_seq, prot_ss, pop.rule[i], conf.CA.Steps, conf.CA.Consensus)
			pop.fitness[i] = Fitness(ca)
		}(&pop, i) //preciso definir o que por aqui
		if i % 100 == 0 && i > 0{
			fmt.Println("Waiting", i)
			wg1.Wait()
		}
	}
	wg1.Wait()

	fmt.Println("População inicial = ", conf.EDA.Population, "OK")

	var wg2 sync.WaitGroup
	for i := 0; i < conf.EDA.Generations; i++ {
		fmt.Println("Generation", i+1)
		fmt.Println("Adjusting probabilities...")
		probs.AdjustProbs(pop, conf.EDA.Selection, conf.EDA.Tournament)
		fmt.Println("OK")

		if probs.Converged() {
			fmt.Println("Probabilities converged\nDONE")
			break
		}

		if (i+1)%conf.EDA.SaveSteps == 0 {
			ioutil.WriteFile(fmt.Sprintf("%s_%d", conf.EDA.OutputProbs, i+1), []byte(probs.String()), 0644)
		}


		for j := 0; j < len(pop.rule); j++ {
			wg2.Add(1)
			go func(pop *Population, j int) {
				defer wg2.Done()
				pop.rule[j] = probs.GenRule()
				ca, _ := ca.Create1D(prot_id, prot_seq, prot_ss, pop.rule[j], conf.CA.Steps, conf.CA.Consensus)
				pop.fitness[j] = Fitness(ca)

			}(&pop, j)
			if j % 10 == 0 && j > 0{
				fmt.Println("Waiting", j)
				wg2.Wait()
			}
		}
		fmt.Printf("Wait - Setting new rule")
		wg2.Wait()
		fmt.Println("OK")
	}

	err := ioutil.WriteFile(conf.EDA.OutputProbs, []byte(probs.String()), 0644)
	if err != nil {
		fmt.Println("Erro gravar as probabilidades")
		fmt.Println(probs)
	}

	return nil
}

type Tournament struct {
	rule    []*rules.Rule
	fitness []float64
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

func (p *Probs) AdjustProbs(pop Population, n_selection, n_tournament int) {
	// selection := make([]*ca.CellAuto1D, n_selection)
	mean_fitness := 0.0

	var tournament Tournament
	tournament.rule = make([]*rules.Rule, n_tournament)
	tournament.fitness = make([]float64, n_tournament)

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
	// var wg2 sync.WaitGroup
	// var wg3 sync.WaitGroup

	for j := 0; j < n_selection; j++ {
		for i := 0; i < n_tournament; i++ {
			x := rand.Intn(len(pop.rule))
			tournament.rule[i] = pop.rule[x]
			tournament.fitness[i] = pop.fitness[x]
		}

		sort.Sort(sort.Reverse(tournament))

		mean_fitness += tournament.fitness[0]

		n := len(tournament.rule[0].Code)
		for ln := 0; ln < n; ln++ {
			for c := 0; c < n; c++ {
				for rn := 0; rn < n; rn++ {
					index := bytes.IndexByte(p.rulePrm.TransitionStates, tournament.rule[0].Code[ln][c][rn])
					count[ln][c][rn][index] += 1
				}
			}
		}
	}

	n := len(tournament.rule[0].Code)
	for ln := 0; ln < n; ln++ {
		for c := 0; c < n; c++ {
			for rn := 0; rn < n; rn++ {
				for i := 0; i < len(count[ln][c][rn]); i++ {
					p.probs[ln][c][rn][i] = count[ln][c][rn][i] / float64(n_selection)
				}
			}
		}
	}

	fmt.Println("Mean Fitness:", mean_fitness/float64(n_selection))
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
