package disteda

import (
	"encoding/json"
	"fmt"
	"hash/adler32"

	zmq "github.com/pebbe/zmq4"

	"github.com/jgcarvalho/zeca/rules"
)

func RunMaster(conf Config) {

	sender, _ := zmq.NewSocket(zmq.PUSH)
	defer sender.Close()
	sender.Bind("tcp://*:" + conf.Dist.PortA)

	receiver, _ := zmq.NewSocket(zmq.PULL)
	defer receiver.Close()
	receiver.Bind("tcp://*:" + conf.Dist.PortB)

	// gerar (ou ler -> TODO) as probabilidades iniciais
	r, _ := rules.Create(conf.CA.InitStates, conf.CA.TransStates, conf.CA.HasJoker, conf.CA.R)
	p := NewProbs(r.Prm)

	var pop []Individual
	pop = make([]Individual, conf.EDA.Population/conf.EDA.Tournament)

	fmt.Print("Press Enter when the workers are ready: ")
	var line string
	fmt.Scanln(&line)
	fmt.Println("Sending tasks to workers...")

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

		go func(b *[]byte) {
			for i := 0; i < len(pop); i++ {
				sender.Send(string(*b), 0)
			}
		}(&b)

		for i := 0; i < len(pop); {
			m, err := receiver.Recv(0)
			if err == nil {
				json.Unmarshal([]byte(m), &pop[i])
				if prob.PID == pop[i].PID {
					fmt.Printf("Individuo id: %d rid: %d g: %d, score: %f\n", g*len(pop)+i, pop[i].PID, pop[i].Generation, pop[i].Fitness)
					i++
				} else {
					fmt.Println(prob.PID, pop[i].PID)
				}

			} else {
				fmt.Println(err)
			}
		}

		// imprimir e as estatisticas
	}
	// salvar a melhor regra

}
