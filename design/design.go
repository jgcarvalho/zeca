package design

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/jgcarvalho/zeca/ca"
	"github.com/jgcarvalho/zeca/metrics"
	"github.com/jgcarvalho/zeca/proteindb"
	"github.com/jgcarvalho/zeca/rules"
)

type Probs struct {
	probs   [][][][]float64
	rulePrm rules.Params
}

type Population struct {
	rule    []*rules.Rule
	fitness []float64
}

//
// type Selection struct {
// 	rule    []*rules.Rule
// 	fitness []float64
// }

func Run(conf Config) error {
	rand.Seed(time.Now().UTC().UnixNano())

	fmt.Println("Loading proteins...")
	id, start, end, err := proteindb.GetProteins(conf.ProteinDB)
	if err != nil {
		panic(err)
	}

	fmt.Println("Initializing probabilities...")
	r, _ := rules.Create(conf.CA.InitStates, conf.CA.TransStates, conf.CA.HasJoker, conf.CA.R)
	probs := NewProbs(r.Prm)
	// fmt.Println(probs)

	var pop Population

	pop.rule = make([]*rules.Rule, conf.Design.Population)
	pop.fitness = make([]float64, conf.Design.Population)

	pop.rule[0], _ = rules.Create(conf.CA.InitStates, conf.CA.TransStates, conf.CA.HasJoker, conf.CA.R)
	cellauto, _ := ca.Create1D(id, start, end, pop.rule[0], conf.CA.Steps, conf.CA.Consensus)
	pop.fitness[0] = Fitness(cellauto)

	// cellAuto := make([]*ca.CellAuto1D, conf.EDA.Population)

	// var wg1 sync.WaitGroup
	// tmp := 0
	for i := 1; i < conf.Design.Population; i++ {
		// wg1.Add(1)
		// go func(pop *Population, i int) {
		// 	defer wg1.Done()
		pop.rule[i] = probs.GenRule()
		cellauto.SetRule(pop.rule[i])
		pop.fitness[i] = Fitness(cellauto)
		// }(&pop, i) //preciso definir o que por aqui
		// if i%100 == 0 && i > 0 {
		// 	fmt.Println("Waiting", i)
		// 	wg1.Wait()
		// }
	}
	// wg1.Wait()
	pop.save("")

	// var selection Selection
	// selection.rule = make([]*rules.Rule, conf.Design.Selection)
	// selection.fitness = make([]float64, conf.Design.Selection)

	// var tournament Tournament
	// tournament.rule = make([]*rules.Rule, conf.Design.Tournament)
	// tournament.fitness = make([]float64, conf.Design.Tournament)

	// sort.Sort(sort.Reverse(pop))
	// for j := 0; j < conf.Design.Selection; j++ {
	// 	winner := len(pop.rule)
	// 	for k := 0; k < conf.Design.Tournament; k++ {
	// 		x := rand.Intn(len(pop.rule))
	// 		if x < winner {
	// 			winner = x
	// 		}
	// 	}
	//
	// 	//sort.Sort(sort.Reverse(tournament))
	//
	// 	selection.rule[j] = pop.rule[winner]
	// 	selection.fitness[j] = pop.fitness[winner]
	// }
	// fmt.Println("Selection OK")
	// plot.Histogram(pop.fitness, selection.fitness, 0)
	//
	return nil
}

func (p Population) save(fn string) error {
	f, _ := os.Create("./population")
	w := bufio.NewWriter(f)

	codes := append(p.rule[0].Prm.StrStartStates, p.rule[0].Prm.StrTransitionStates...)
	var toprint string
	for i := 0; i < len(p.fitness); i++ {
		for c := 0; c < len(p.rule[i].Code); c++ {
			for ln := 0; ln < len(p.rule[i].Code); ln++ {
				for rn := 0; rn < len(p.rule[i].Code); rn++ {
					toprint += fmt.Sprintf("%s, ", codes[p.rule[i].Code[ln][c][rn]])
				}
			}
		}
		toprint += fmt.Sprintf("%f\n", p.fitness[i])
		w.WriteString(toprint)
		toprint = ""
	}
	return nil
}

func (p Population) Len() int {
	return len(p.rule)
}

func (p Population) Swap(i, j int) {
	p.rule[i], p.rule[j] = p.rule[j], p.rule[i]
	p.fitness[i], p.fitness[j] = p.fitness[j], p.fitness[i]
}

func (p Population) Less(i, j int) bool {
	return p.fitness[i] < p.fitness[j]
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
				if c == 0 {
					rule.Code[ln][c][rn] = 0
				} else {
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

	}

	return &rule
}
