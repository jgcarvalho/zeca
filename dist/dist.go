package dist

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hash/adler32"
	"math/rand"
	"sort"
	"time"

	"github.com/apcera/nats"
	"github.com/jgcarvalho/zeca/ca"
	"github.com/jgcarvalho/zeca/metrics"
	"github.com/jgcarvalho/zeca/proteindb"
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

func RunMaster(conf Config) {

	conn, err := nats.Connect(conf.Dist.NatsServerURL)
	defer conn.Close()
	if err != nil {
		fmt.Println("Erro na conexao do master ao nats")
		return
	}

	sub, err := conn.SubscribeSync(conf.Dist.TopicFromSlave)
	if err != nil {
		fmt.Println("Erro na tentativa do master assinar o tópico fitness")
		return
	}

	// gerar (ou ler -> TODO) as probabilidades iniciais
	r, _ := rules.Create(conf.CA.InitStates, conf.CA.TransStates, conf.CA.HasJoker, conf.CA.R)
	p := NewProbs(r.Prm)

	var pop Population
	pop.rule = make([]*rules.Rule, conf.EDA.Population/conf.EDA.Tournament)
	pop.fitness = make([]float64, conf.EDA.Population/conf.EDA.Tournament)

	// var ind Individual

	fmt.Println("RUNNING MASTER")
	// para cada geracao
	for g := 0; g < conf.EDA.Generations; g++ {
		fmt.Println("GERACAO", g)
		// atualizar as probabilidades de acordo com a populacao dessa geracao
		if g != 0 {
			// TODO refazer a funcao para trabalhar com toda a populacao
			p.AdjustProbs(pop)
		}

		// Publicar as probabilidades
		tmp, _ := json.Marshal(p.probs)
		pid := adler32.Checksum(tmp)
		prob := &Probabilities{PID: pid, Generation: g, Data: p.probs}
		b, _ := json.Marshal(prob)
		conn.Publish(conf.Dist.TopicFromMaster, b)

		// até a populacao minima ser atingida
		for i := 0; i < len(pop.rule); {
			// fmt.Println("GERACAO", g, "INDIVIDUO", i)
			// recebe regras e fitness dos slaves
			m, err := sub.NextMsg(1 * time.Second)
			if err == nil {
				var ind Individual
				json.Unmarshal(m.Data, &ind)
				if prob.PID == ind.PID {
					//É preciso copiar
					pop.rule[i] = ind.Rule
					pop.fitness[i] = ind.Fitness
					fmt.Printf("Individuo id: %d rid: %d g: %d, score: %f\n", g*len(pop.rule)+i, ind.PID, ind.Generation, ind.Fitness)
					i++
				}
			}
		}

		// imprimir e as estatisticas
	}
	// salvar a melhor regra

}

func RunSlaveAsync(conf Config) {
	rand.Seed(time.Now().UTC().UnixNano())

	fmt.Println("Loading proteins...")
	id, start, end, err := proteindb.GetProteins(conf.ProteinDB)
	if err != nil {
		fmt.Println("Erro no banco de DADOS")
		panic(err)
	}

	conn, err := nats.Connect(conf.Dist.NatsServerURL)
	defer conn.Close()
	if err != nil {
		fmt.Println("Erro na conexao do master ao nats")
		return
	}

	var prob Probabilities
	sub, _ := conn.Subscribe(conf.Dist.TopicFromMaster, func(m *nats.Msg) {
		json.Unmarshal(m.Data, &prob)
	})

	var tourn Tournament
	tourn.rule = make([]*rules.Rule, conf.EDA.Tournament)
	tourn.fitness = make([]float64, conf.EDA.Tournament)

	r, _ := rules.Create(conf.CA.InitStates, conf.CA.TransStates, conf.CA.HasJoker, conf.CA.R)
	p_tmp := NewProbs(r.Prm)
	cellAuto := make([]*ca.CellAuto1D, conf.EDA.Tournament)
	for i := 0; i < conf.EDA.Tournament; i++ {
		tourn.rule[i] = p_tmp.GenRule()
		cellAuto[i], _ = ca.Create1D(id, start, end, tourn.rule[i], conf.CA.Steps, conf.CA.Consensus)
	}

	for sub.IsValid() {

		if prob.PID != 0 {
			for i := 0; i < len(tourn.rule); i++ {
				copy(p_tmp.probs, prob.Data)
				tourn.rule[i] = p_tmp.GenRule()
				cellAuto[i].SetRule(tourn.rule[i])
				tourn.fitness[i] = Fitness(cellAuto[i])
				fmt.Println("Individuo", i, "Fitness", tourn.fitness[i])
			}
			sort.Sort(sort.Reverse(tourn))
			ind := &Individual{PID: prob.PID, Generation: prob.Generation, Rule: tourn.rule[0], Fitness: tourn.fitness[0]}
			b, _ := json.Marshal(ind)
			fmt.Println("Fitness selecionado", tourn.fitness[0])
			conn.Publish(conf.Dist.TopicFromSlave, b)
		}
	}
}

func RunSlave(conf Config) {
	rand.Seed(time.Now().UTC().UnixNano())

	fmt.Println("Loading proteins...")
	id, start, end, err := proteindb.GetProteins(conf.ProteinDB)
	if err != nil {
		fmt.Println("Erro no banco de DADOS")
		panic(err)
	}

	conn, err := nats.Connect(conf.Dist.NatsServerURL)
	defer conn.Close()
	if err != nil {
		fmt.Println("Erro na conexao do master ao nats")
		return
	}

	sub, err := conn.SubscribeSync(conf.Dist.TopicFromMaster)
	if err != nil {
		fmt.Println("Erro na tentativa do slave assinar o tópico probabilities")
		return
	}
	var prob Probabilities

	//
	var tourn Tournament
	tourn.rule = make([]*rules.Rule, conf.EDA.Tournament)
	tourn.fitness = make([]float64, conf.EDA.Tournament)

	r, _ := rules.Create(conf.CA.InitStates, conf.CA.TransStates, conf.CA.HasJoker, conf.CA.R)
	p_tmp := NewProbs(r.Prm)
	cellAuto := make([]*ca.CellAuto1D, conf.EDA.Tournament)
	for i := 0; i < conf.EDA.Tournament; i++ {
		tourn.rule[i] = p_tmp.GenRule()
		cellAuto[i], _ = ca.Create1D(id, start, end, tourn.rule[i], conf.CA.Steps, conf.CA.Consensus)
	}
	// checa e recebe novas probabilidades
	for {
		m, err := sub.NextMsg(5 * time.Millisecond)
		if err == nil {
			// para cada individuo no torneio
			// gera uma regra de acordo com a probabilidade atual
			// roda o automato celular
			// calcula o fitness
			// seleciona o vencedor do torneio
			// retorna sua regra e seu fitness)
			json.Unmarshal(m.Data, &prob)
			fmt.Printf("PID: %d, Geracacao: %d\n", prob.PID, prob.Generation)
			for i := 0; i < len(tourn.rule); i++ {
				p_tmp.probs = prob.Data
				tourn.rule[i] = p_tmp.GenRule()
				// cellAuto[i], _ = ca.Create1D(id, start, end, tourn.rule[i], conf.CA.Steps, conf.CA.Consensus)
				// fmt.Println(cellAuto[i])
				// fmt.Println(prob.Data)
				cellAuto[i].SetRule(tourn.rule[i])
				tourn.fitness[i] = Fitness(cellAuto[i])
				fmt.Println("Individuo", i, "Fitness", tourn.fitness[i])
			}
			sort.Sort(sort.Reverse(tourn))
			ind := &Individual{PID: prob.PID, Generation: prob.Generation, Rule: tourn.rule[0], Fitness: tourn.fitness[0]}
			b, _ := json.Marshal(ind)
			fmt.Println("Fitness selecionado", tourn.fitness[0])
			conn.Publish(conf.Dist.TopicFromSlave, b)
			// i++
		} else if prob.PID != 0 {
			for i := 0; i < len(tourn.rule); i++ {
				p_tmp.probs = prob.Data
				tourn.rule[i] = p_tmp.GenRule()
				cellAuto[i].SetRule(tourn.rule[i])
				tourn.fitness[i] = Fitness(cellAuto[i])
				fmt.Println("Individuo", i, "Fitness", tourn.fitness[i])
			}
			sort.Sort(sort.Reverse(tourn))
			ind := &Individual{PID: prob.PID, Generation: prob.Generation, Rule: tourn.rule[0], Fitness: tourn.fitness[0]}
			b, _ := json.Marshal(ind)
			fmt.Println("Fitness selecionado", tourn.fitness[0])
			conn.Publish(conf.Dist.TopicFromSlave, b)
			// i++
		}
	}
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
				fmt.Printf("Contagem %f %f %f %f \n", count[ln][c][rn][0], count[ln][c][rn][1], count[ln][c][rn][2], count[ln][c][rn][3])
				fmt.Printf("len pop rule %d \n", len(pop.rule))
				fmt.Printf("Probabilidade %f %f %f %f \n", p.probs[ln][c][rn][0], p.probs[ln][c][rn][1], p.probs[ln][c][rn][2], p.probs[ln][c][rn][3])
			}
		}
	}

}

// func (p *Probs) AdjustProbs(pop Population, n_tournament, n_selection int) {
// 	// selection := make([]*ca.CellAuto1D, n_selection)
// 	mean_fitness := 0.0
//
// 	var tournament Tournament
// 	tournament.rule = make([]*rules.Rule, n_tournament)
// 	tournament.fitness = make([]float64, n_tournament)
//
// 	var count [][][][]float64
// 	st := rules.RuleStates(p.rulePrm)
// 	count = make([][][][]float64, len(st))
// 	for ln := range st {
// 		count[ln] = make([][][]float64, len(st))
// 		for c := range st {
// 			count[ln][c] = make([][]float64, len(st))
// 			for rn := range st {
// 				//p.probs[ln][c][rn] = prm.transitionStates[rand.Intn(len(prm.transitionStates))]
// 				count[ln][c][rn] = make([]float64, len(p.rulePrm.TransitionStates))
// 			}
// 		}
// 	}
//
// 	for j := 0; j < n_selection; j++ {
// 		for i := 0; i < n_tournament; i++ {
// 			x := rand.Intn(len(pop.rule))
// 			tournament.rule[i] = pop.rule[x]
// 			tournament.fitness[i] = pop.fitness[x]
// 		}
//
// 		sort.Sort(sort.Reverse(tournament))
//
// 		mean_fitness += tournament.fitness[0]
//
// 		n := len(tournament.rule[0].Code)
// 		for ln := 0; ln < n; ln++ {
// 			for c := 0; c < n; c++ {
// 				for rn := 0; rn < n; rn++ {
// 					index := bytes.IndexByte(p.rulePrm.TransitionStates, tournament.rule[0].Code[ln][c][rn])
// 					count[ln][c][rn][index] += 1
// 				}
// 			}
// 		}
// 	}
//
// 	n := len(tournament.rule[0].Code)
// 	for ln := 0; ln < n; ln++ {
// 		for c := 0; c < n; c++ {
// 			for rn := 0; rn < n; rn++ {
// 				for i := 0; i < len(count[ln][c][rn]); i++ {
// 					p.probs[ln][c][rn][i] = count[ln][c][rn][i] / float64(n_selection)
// 				}
// 			}
// 		}
// 	}
//
// 	fmt.Println("Mean Fitness:", mean_fitness/float64(n_selection))
// }

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
