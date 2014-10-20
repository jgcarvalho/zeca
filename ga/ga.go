package ga

import (
	"fmt"
	"math/rand"
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
	proteins := proteindb.LoadProteinsFromMongo(conf.ProteinDB.Ip, conf.ProteinDB.Name, conf.ProteinDB.Collection)
	prot_id := "20primeiras"
	prot_seq := "#"
	prot_ss := "#"
	for i := 0; i < len(proteins); i++ {
		prot_seq += proteins[i].Chains[0].Seq_pdb + "#"
		prot_ss += proteins[i].Chains[0].Ss3_cons_all + "#"
	}

	var pop Population

	pop.rule = make([]*rules.Rule, conf.GA.Population)
	pop.fitness = make([]float64, conf.GA.Population)

	var wg1 sync.WaitGroup
	// tmp := 0
	for i := 0; i < conf.GA.Population; i++ {
		wg1.Add(1)
		go func(pop *Population, i int) {
			defer wg1.Done()
			pop.rule[i], _ = rules.Create(conf.CA.InitStates, conf.CA.TransStates, conf.CA.HasJoker, conf.CA.R)
			ca, _ := ca.Create1D(prot_id, prot_seq, prot_ss, pop.rule[i], conf.CA.Steps, conf.CA.Consensus)
			pop.fitness[i] = Fitness(ca)
		}(&pop, i) //preciso definir o que por aqui
		if i%100 == 0 && i > 0 {
			fmt.Println("Waiting", i)
			wg1.Wait()
		}
	}
	wg1.Wait()

	var selection Selection
	selection.rule = make([]*rules.Rule, conf.GA.Selection)
	selection.fitness = make([]float64, conf.GA.Selection)

	// var tournament Tournament
	// tournament.rule = make([]*rules.Rule, conf.GA.Tournament)
	// tournament.fitness = make([]float64, conf.GA.Tournament)
	var wg2 sync.WaitGroup
	for i := 0; i < conf.GA.Generations; i++ {
		fmt.Println("Gen", i)
		sort.Sort(sort.Reverse(pop))
		for j := 0; j < conf.GA.Selection; j++ {
			winner := len(pop.rule)
			for k := 0; k < conf.GA.Tournament; k++ {
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
		plot.Histogram(pop.fitness, selection.fitness, i)
		fmt.Println("Plot", i, "ok")
		for p := 0; p < len(pop.rule); p++ {
			wg2.Add(1)
			go func(pop *Population, p int) {
				defer wg2.Done()
				s := rand.Intn(len(selection.rule))
				Mutate(selection.rule[s], conf.GA.Mutation)
				pop.rule[p] = selection.rule[s]
				ca, _ := ca.Create1D(prot_id, prot_seq, prot_ss, pop.rule[p], conf.CA.Steps, conf.CA.Consensus)
				pop.fitness[p] = Fitness(ca)
			}(&pop, p)
			if p%10 == 0 && p > 0 {
				fmt.Println("Waiting", p)
				wg2.Wait()
			}
		}
		fmt.Println("New pop OK")
	}

	return nil
}

func Mutate(rule *rules.Rule, mutRate float64) {
	n := len(rule.Code)
	for c := 0; c < n; c++ {
		for ln := 0; ln < n; ln++ {
			for rn := 0; rn < n; rn++ {
				if rule.Fixed[c][ln][rn] == false && rand.Float64() < mutRate {
					//pode não mudar
					rule.Code[c][ln][rn] = rule.Prm.TransitionStates[rand.Intn(len(rule.Prm.TransitionStates))]
				}
			}
		}
	}
}

func CrossOver(rule1, rule2 *rules.Rule) {
	//n := len(rule1.Code)

}

// type Tournament struct {
// 	rule    []*rules.Rule
// 	fitness []float64
// }

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
