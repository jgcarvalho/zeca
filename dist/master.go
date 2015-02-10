package dist

import (
	"encoding/json"
	"fmt"
	"hash/adler32"
	"time"

	"github.com/apcera/nats"
	"github.com/jgcarvalho/zeca/rules"
)

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
		// conn.Publish(conf.Dist.TopicFromMaster, b)
		conn.Subscribe(conf.Dist.TopicFromMaster, func(m *nats.Msg) {
			conn.Publish(m.Reply, b)
		})

		// até a populacao minima ser atingida
		for i := 0; i < len(pop.rule); {
			// fmt.Println("Waiting")
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
				} else {
					// fmt.Println("Update")
					// err := conn.Flush()
					// if err == nil {
					// 	fmt.Println("Flush OK")
					// 	// Everything has been processed by the server for nc *Conn.
					// } else {
					// 	fmt.Println("Not Flush")
					// }
					fmt.Println(prob.PID, ind.PID)
					// conn.Publish(conf.Dist.TopicFromMaster, b)

				}
				if (i+1)%10 == 0 {
					fmt.Println("Remember")
					// conn.Publish(conf.Dist.TopicFromMaster, b)
				}
			} else {
				fmt.Println("Publishing")
				// conn.Publish(conf.Dist.TopicFromMaster, b)
			}

		}

		// imprimir e as estatisticas
	}
	// salvar a melhor regra

}
