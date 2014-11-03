package design

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/jgcarvalho/zeca/ca"
	"github.com/jgcarvalho/zeca/metrics"
	"github.com/jgcarvalho/zeca/plot"
	"github.com/jgcarvalho/zeca/proteindb"
	"github.com/jgcarvalho/zeca/rules"
)

type Population struct {
	rule    []*rules.Rule
	fitness []float64
}

type Selection struct {
	rule    []*rules.Rule
	fitness []float64
}

func Run(conf Config) error {
	rand.Seed(time.Now().UTC().UnixNano())

	fmt.Println("Loading proteins...")
	id, start, end, err := proteindb.GetProteins(conf.ProteinDB)
	if err != nil {
		panic(err)
	}

	var pop Population

	pop.rule = make([]*rules.Rule, conf.Design.Population)
	pop.fitness = make([]float64, conf.Design.Population)

	var wg1 sync.WaitGroup
	// tmp := 0
	for i := 0; i < conf.Design.Population; i++ {
		wg1.Add(1)
		go func(pop *Population, i int) {
			defer wg1.Done()
			pop.rule[i], _ = rules.Create(conf.CA.InitStates, conf.CA.TransStates, conf.CA.HasJoker, conf.CA.R)
			ca, _ := ca.Create1D(id, start, end, pop.rule[i], conf.CA.Steps, conf.CA.Consensus)
			pop.fitness[i] = Fitness(ca)
		}(&pop, i) //preciso definir o que por aqui
		if i%100 == 0 && i > 0 {
			fmt.Println("Waiting", i)
			wg1.Wait()
		}
	}
	wg1.Wait()
	pop.save("")

	var selection Selection
	selection.rule = make([]*rules.Rule, conf.Design.Selection)
	selection.fitness = make([]float64, conf.Design.Selection)

	// var tournament Tournament
	// tournament.rule = make([]*rules.Rule, conf.Design.Tournament)
	// tournament.fitness = make([]float64, conf.Design.Tournament)

	sort.Sort(sort.Reverse(pop))
	for j := 0; j < conf.Design.Selection; j++ {
		winner := len(pop.rule)
		for k := 0; k < conf.Design.Tournament; k++ {
			x := rand.Intn(len(pop.rule))
			if x < winner {
				winner = x
			}
		}

		//sort.Sort(sort.Reverse(tournament))

		selection.rule[j] = pop.rule[winner]
		selection.fitness[j] = pop.fitness[winner]
	}
	fmt.Println("Selection OK")
	plot.Histogram(pop.fitness, selection.fitness, 0)

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
