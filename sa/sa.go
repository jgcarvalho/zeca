package sa

import (
	"bitbucket.org/zgcarvalho/zeca/ca"
	"bitbucket.org/zgcarvalho/zeca/metrics"
	"bitbucket.org/zgcarvalho/zeca/proteindb"
	"bitbucket.org/zgcarvalho/zeca/rules"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"time"
)

type Solution struct {
	rule    *rules.Rule
	fitness float64
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

	fmt.Println("Initializing simulated annealing")
	rule, _ := rules.Create(conf.CA.InitStates, conf.CA.TransStates, conf.CA.HasJoker, conf.CA.R)
	cellauto, _ := ca.Create1D(prot_id, prot_seq, prot_ss, rule, conf.CA.Steps, conf.CA.Consensus)

	solution := &Solution{rule, Fitness(cellauto)}
	solution_new := &Solution{rule, Fitness(cellauto)}
	solution_best := &Solution{rule, Fitness(cellauto)}

	alpha := math.Pow((conf.SA.Tfinal / conf.SA.Tini), 1.0/float64(conf.SA.OuterLoop))

	temp := conf.SA.Tini

	for outer := 0; outer < conf.SA.OuterLoop; outer++ {
		fmt.Println("@ Temperature", temp)
		for inner := 0; inner < conf.SA.InnerLoop; inner++ {
			solution_new.Neighbor(solution)
			cellauto.SetRule(solution_new.rule)
			solution_new.fitness = Fitness(cellauto)
			if solution_new.fitness >= solution.fitness {
				fmt.Println("Update", solution.fitness, "->", solution_new.fitness)
				solution.rule = solution_new.rule
				solution.fitness = solution_new.fitness
				//fmt.Println("Check update", solution.fitness, "=", solution_new.fitness)
				if solution.fitness > solution_best.fitness {
					fmt.Println("Update BEST", solution_best.fitness, "->", solution.fitness)
					solution_best.rule = solution.rule
					solution_best.fitness = solution.fitness
				}
			} else if math.Exp(solution_new.fitness-solution.fitness/temp) > rand.Float64() {

				fmt.Println("*Update", solution.fitness, "->", solution_new.fitness)
				solution.rule = solution_new.rule
				solution.fitness = solution_new.fitness
				//fmt.Println("*Check update", solution.fitness, "=", solution_new.fitness)
			}
		}
		temp = alpha * temp
	}

	err := ioutil.WriteFile(conf.Rules.Output, []byte(solution_best.rule.String()), 0644)
	if err != nil {
		fmt.Println("Erro gravar a melhor regra")
		fmt.Println(solution_best.rule)
	}
	fmt.Println("Melhor regra", solution_best.fitness)
	return nil
}

// IMPLEMENTAR A FUNCAO DE GERAR O VIZINHO
func (s_new *Solution) Neighbor(s *Solution) {
	*s_new.rule = *s.rule
	update := false
	n := len(s_new.rule.Code)
	nts := len(s_new.rule.Prm.TransitionStates)
	var c, ln, rn int
	var new_code byte
	for update != true {
		c = rand.Intn(n)
		ln = rand.Intn(n)
		rn = rand.Intn(n)
		if s_new.rule.Fixed[ln][c][rn] == true {
			continue
		} else {
			for update != true {
				new_code = s_new.rule.Prm.TransitionStates[0] + byte(rand.Intn(nts))
				if s_new.rule.Code[ln][c][rn] != new_code {
					s_new.rule.Code[ln][c][rn] = new_code
					update = true
				}
			}
		}
	}
}

func Fitness(c *ca.CellAuto1D) float64 {
	c.Run()
	cm := c.ConfusionMatrix()
	cba := metrics.CBA(cm)
	//fmt.Println("CBA: ", cba)
	return cba
}
