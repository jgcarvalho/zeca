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

	// sub, err := conn.SubscribeSync(conf.Dist.TopicFromMaster)
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
		// m, err := sub.NextMsg(100 * time.Millisecond)
		m, err := conn.Request(conf.Dist.TopicFromMaster, []byte("get"), 2*time.Second)
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
				// p_tmp.probs = prob.Data
				copy(p_tmp.probs, prob.Data)
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
			fmt.Println("Waiting new probs")
			// err := conn.Flush()
			// if err == nil {
			// 	fmt.Println("Flush OK")
			// 	// Everything has been processed by the server for nc *Conn.
			// } else {
			// 	fmt.Println("Not Flush")
			// }
			fmt.Println(prob.PID)
			// for i := 0; i < len(tourn.rule); i++ {
			// 	// p_tmp.probs = prob.Data
			// 	copy(p_tmp.probs, prob.Data)
			// 	tourn.rule[i] = p_tmp.GenRule()
			// 	cellAuto[i].SetRule(tourn.rule[i])
			// 	tourn.fitness[i] = Fitness(cellAuto[i])
			// 	fmt.Println("Individuo", i, "Fitness", tourn.fitness[i])
			// }
			// sort.Sort(sort.Reverse(tourn))
			// ind := &Individual{PID: prob.PID, Generation: prob.Generation, Rule: tourn.rule[0], Fitness: tourn.fitness[0]}
			// b, _ := json.Marshal(ind)
			// fmt.Println("Fitness selecionado", tourn.fitness[0])
			// conn.Publish(conf.Dist.TopicFromSlave, b)
			// i++
		}
	}
}
