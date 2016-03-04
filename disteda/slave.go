package disteda

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sort"
	"time"

	"github.com/jgcarvalho/zeca/db"

	"github.com/jgcarvalho/zeca/ca"
	"github.com/jgcarvalho/zeca/rules"
	zmq "github.com/pebbe/zmq4"
)

func RunSlave(conf Config) {

	// Cria o receptor que recebe a probabilidade emitida pelo master na porta A
	receiver, _ := zmq.NewSocket(zmq.PULL)
	defer receiver.Close()
	receiver.Connect("tcp://" + conf.Dist.MasterURL + ":" + conf.Dist.PortA)

	// Cria o emissor que envia o individuo vencedor do torneio na rede pela
	// porta B
	sender, _ := zmq.NewSocket(zmq.PUSH)
	defer sender.Close()
	sender.Connect("tcp://" + conf.Dist.MasterURL + ":" + conf.Dist.PortB)

	// semente randomica
	rand.Seed(time.Now().UTC().UnixNano())

	// Le os dados das proteinas no DB
	fmt.Println("Loading proteins...")
	id, start, end, err := db.GetProteins(conf.DB)
	if err != nil {
		fmt.Println("Erro no banco de DADOS")
		panic(err)
	}
	fmt.Println("Done")
	// ? Ha vantagem em enviar um sinal de Ok (proteinas lidas) para o master?

	var prob Probabilities

	var tourn Tournament
	tourn = make([]Individual, conf.EDA.Tournament)
	// tourn.rule = make([]*rules.Rule, conf.EDA.Tournament)
	// tourn.fitness = make([]float64, conf.EDA.Tournament)

	r, _ := rules.Create(conf.CA.InitStates, conf.CA.TransStates, conf.CA.HasJoker, conf.CA.R)
	// probabilidade temporaria para ser substituida pelas recebidas
	p_tmp := NewProbs(r.Prm)
	cellAuto := make([]*ca.CellAuto1D, conf.EDA.Tournament)
	for i := 0; i < conf.EDA.Tournament; i++ {
		// tourn.rule[i] = p_tmp.GenRule()
		tourn[i].Rule = p_tmp.GenRule()

		// cellAuto[i], _ = ca.Create1D(id, start, end, tourn.rule[i], conf.CA.Steps, conf.CA.Consensus)
		cellAuto[i], _ = ca.Create1D(id, start, end, tourn[i].Rule, conf.CA.Steps, conf.CA.Consensus)
	}

	// Individuo vencedor do torneio
	var (
		ind    Individual
		b      []byte
		m      string
		conerr error
	)

	for {
		// m é a mensagem com as probabilidades
		m, conerr = receiver.Recv(0)
		// m, err := conn.Request(conf.Dist.TopicFromMaster, []byte("get"), 2*time.Second)
		if conerr == nil {
			// para cada individuo no torneio
			// gera uma regra de acordo com a probabilidade atual
			// roda o automato celular
			// calcula o fitness
			// seleciona o vencedor do torneio
			// retorna sua regra e seu fitness)

			// converte a probabilidade recebida em JSON para uma estrutura
			json.Unmarshal([]byte(m), &prob)
			fmt.Printf("PID: %d, Geracacao: %d\n", prob.PID, prob.Generation)
			// for i := 0; i < len(tourn.rule); i++ {
			for i := 0; i < len(tourn); i++ {

				// copia a probabilidade recebida para a probabilidade dos individuos
				copy(p_tmp.probs, prob.Data)
				// gera a regra e atribui ao membro do torneio
				tourn[i].Rule = p_tmp.GenRule()
				// define a regra do automato como sendo a nova regra
				cellAuto[i].SetRule(tourn[i].Rule)
				// retorna o fitness e outras medidas de desempenho do autômato
				tourn[i].Fitness, tourn[i].Q3 = FitnessAndQ3(cellAuto[i])

				// fmt.Println("Individuo", i, "Fitness", tourn.fitness[i])
				fmt.Println("Individuo", i, "Fitness", tourn[i].Fitness)

			}

			// Ordena os individuos do torneio de acordo com o fitness (maior primeiro)
			sort.Sort(sort.Reverse(tourn))
			// ind.PID, ind.Generation, ind.Rule, ind.Fitness = prob.PID, prob.Generation, tourn.rule[0], tourn.fitness[0]
			ind.PID, ind.Generation, ind.Rule, ind.Fitness, ind.Q3 = prob.PID, prob.Generation, tourn[0].Rule, tourn[0].Fitness, tourn[0].Q3

			// Codifica o individuo vencedor em JSON e envia para o master
			b, _ = json.Marshal(ind)
			fmt.Println("Fitness selecionado", tourn[0].Fitness)
			sender.Send(string(b), 0)

		} else {
			// Erro na conexão
			fmt.Println(err)
		}

	}
}
