package dist

import (
	"encoding/json"
	"fmt"
	"hash/adler32"
	"math/rand"
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
	Dat        Probs
}

type Probs struct {
	probs   [][][][]float64
	rulePrm rules.Params
}

type Population struct {
  rule    []*rules.Rule
  fitness []float64
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

  pop := make([]Population,conf.EDA.Population)

	// para cada geracao
	for g := 0; g < conf.EDA.Generations; g++ {

		// atualizar as probabilidades de acordo com a populacao dessa geracao
		if g != 0 {
      // TODO refazer a funcao para trabalhar com toda a populacao
      p.AdjustProbs(pop Population, n_selection, n_tournament int)
		}

		// Publicar as probabilidades
		tmp, _ := json.Marshal(p)
		pid := adler32.Checksum(tmp)
		prob := &Probabilities{PID: pid, Generation: g}
		b, _ := json.Marshal(prob)
		conn.Publish(conf.Dist.TopicFromMaster, b)

    // até a populacao minima ser atingida
    for i := 0; i < len(pop); i++{
      // recebe regras e fitness dos slaves
    }



		// imprimir e as estatisticas
	}
	// salvar a melhor regra

}

func RunSlave(conf Config) {
	rand.Seed(time.Now().UTC().UnixNano())

	fmt.Println("Loading proteins...")
	id, start, end, err := proteindb.GetProteins(conf.ProteinDB)
	if err != nil {
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

	// checa e recebe novas probabilidades
	// para cada individuo no torneio
	// gera uma regra de acordo com a probabilidade atual
	// roda o automato celular
	// calcula o fitness
	// seleciona o vencedor do torneio
	// retorna sua regra e seu fitness

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
