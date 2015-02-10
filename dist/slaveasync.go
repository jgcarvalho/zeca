package dist

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sort"
	"time"

	"github.com/apcera/nats"
	"github.com/jgcarvalho/zeca/ca"
	"github.com/jgcarvalho/zeca/proteindb"
	"github.com/jgcarvalho/zeca/rules"
)

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
