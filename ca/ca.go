package ca

import (
	// "github.com/jgcarvalho/zeca/proteindb"
	"bytes"
	"fmt"

	"github.com/jgcarvalho/zeca/rules"
	// "io/ioutil"
	// "strings"
)

// CellAuto é uma estrutura para armazenar um autômato
type CellAuto1D struct {
	id           string
	Begin        []byte
	Expected     []byte
	End          [][]byte
	EndConsensus []byte
	Rule         *rules.Rule
	steps        int
	consensus    int
}

func Create1D(id string, Begin string, Expected string, r *rules.Rule, step int, consensus int) (*CellAuto1D, error) {
	if len(Begin) != len(Expected) {
		return nil, fmt.Errorf("Estado de entrada e a saida esperada tem comprimentos diferentes")
	}
	b := encode(Begin, &r.Prm)
	exp := encode(Expected, &r.Prm)
	EndConsensus := make([]byte, len(b))
	End := make([][]byte, consensus)
	for i := range End {
		End[i] = make([]byte, len(b))
	}
	return &CellAuto1D{id, b, exp, End, EndConsensus, r, step, consensus}, nil
}

// alterar para retornar tp, tn, fp, fn para facilitar o calculo da correlacao
func (ca *CellAuto1D) Run() {
	currentState := make([]byte, len(ca.Begin))
	nextState := make([]byte, len(ca.Begin))
	copy(currentState, ca.Begin)
	copy(nextState, ca.Begin)

	// fmt.Println(currentState)
	for i := 0; i < (ca.steps - ca.consensus); i++ {
		if i%2 == 0 {
			oneStep(ca.Begin, currentState, nextState, ca.Rule)
			// fmt.Println(i, nextState)
		} else {
			oneStep(ca.Begin, nextState, currentState, ca.Rule)
			// fmt.Println(i, currentState)
		}
	}
	// modificar para utilizar um comite dos "11" ultimos estados
	copy(ca.End[0], currentState)
	for i := 0; i < ca.consensus-1; i++ {
		oneStep(ca.Begin, ca.End[i], ca.End[i+1], ca.Rule)
	}
	if ca.consensus == 1 {
		copy(ca.EndConsensus, ca.End[0])
	} else {
		for i := 0; i < len(ca.End[0]); i++ {
			col := make([]byte, len(ca.End))
			for j := 0; j < len(ca.End); j++ {
				col[j] = ca.End[j][i]
			}
			count := 0
			for _, v := range col {
				// maior caso haja em empate, o mais próximo ao fim terá preferencia
				if bytes.Count(col, []byte{v}) >= count {
					count = bytes.Count(col, []byte{v})
					ca.EndConsensus[i] = v
				}
			}
			// fmt.Println("Elements: ", col, " -> ", ca.EndConsensus[i])
		}
	}
}

func oneStep(seq []byte, currentState []byte, nextState []byte, rule *rules.Rule) {
	var state byte
	for c := 1; c < len(currentState)-1; c++ {
		state = rule.Code[currentState[c-1]][currentState[c]][currentState[c+1]]
		if rule.Prm.Hasjoker && state == rule.Prm.TransitionStates[len(rule.Prm.TransitionStates)-1] {
			nextState[c] = seq[c]
		} else {
			nextState[c] = state
		}
	}
}

func (ca *CellAuto1D) SetRule(r *rules.Rule) {
	ca.Rule = r
}

// ConfusionMatrix é um método que retorna a matrix de confusão. A matrix tem formato NxN, onde N é o número
// de estados de transição do autômato. Quando as regras de transição incluem um "coringa" o formato da matrix
// será (N-1)x(N)
func (ca *CellAuto1D) ConfusionMatrix() [][]int {
	n := len(ca.Rule.Prm.TransitionStates)
	np, nr := n, n
	if ca.Rule.Prm.Hasjoker {
		nr -= 1
	}
	cm := make([][]int, nr)
	for i := 0; i < nr; i++ {
		cm[i] = make([]int, np)
	}
	expected, predicted := 0, 0

	for i := range ca.Expected {
		expected = bytes.IndexByte(ca.Rule.Prm.TransitionStates, ca.Expected[i])
		predicted = bytes.IndexByte(ca.Rule.Prm.TransitionStates, ca.EndConsensus[i])

		if (ca.Rule.Prm.Hasjoker && expected == len(ca.Rule.Prm.TransitionStates)-1) || expected == -1 {
			continue
		}

		if predicted == -1 {
			predicted = len(ca.Rule.Prm.TransitionStates) - 1
		}
		cm[expected][predicted] += 1
	}

	// fmt.Println("Confusion Matrix:")
	// fmt.Println(cm)
	return cm
}

// func Run(c CellAuto, rule *rules.Rules) ([]byte, float64, float64, float64, float64, float64, float64, float64, float64) {
// 	currentState := make([]byte, len(c.cell))
// 	nextState := make([]byte, len(c.cell))
// 	copy(currentState, c.cell)
// 	copy(nextState, c.cell)

// 	//fmt.Println(currentState)
// 	for i := 0; i < len(currentState) * 2; i++ {
// 		if i%2 == 0 {
// 			oneStep(c.seq, currentState, nextState, rule)
// 			//fmt.Println(i, nextState)
// 		} else {
// 			oneStep(c.seq, nextState, currentState, rule)
// 			//fmt.Println(i, currentState)
// 		}

// 		if bytes.Equal(currentState,nextState) {
// 			//fmt.Println("Estabilizou no passo", i)
// 			break
// 		}
// 	}
// func encode(start string, End string, rule *rules.Rule) ([]byte, []byte, error) {
// 	Begin := make([]byte, len(start))
// 	Expected := make([]byte, len(End))
// }

// func CreateOne(fn string) *CellAuto {
// 	/* Funcao que cria um automato celular de acordo com o arquivo passado
// 	INPUT:
// 	Nome do arquivo
// 	OUTPUT:
// 	Automato celular com id (nome=pdb id), sequencia(celulas) e estrutura real */

// 	c := new(CellAuto)
// 	c.id = fn
// 	c.seq, c.cell, c.trueSS = loadFile(fn)

// 	/*DEBUG
// 	fmt.Println(c.cell)
// 	fmt.Println(c.trueSS)
// 	println(lines[0][0:5]) */
// 	return c
// }

// func CreateN(fns []string) []CellAuto {
// 	/* Funcao que cria N automatos celulares de acordo com um vetor contendo o nome dos arquivos
// 	INPUT:
// 	Slice com nome dos arquivos das proteinas
// 	OUTPUT:
// 	Slice de automatos celulares com id (nome=pdb id), sequencia(celulas) e estrutura real */

// 	//cria uma slice de automatos celulares com dimensao igual ao numero de arquivos de proteinas
// 	cas := make([]CellAuto, len(fns))

// 	//inicializa os automatos
// 	for i := 0; i < len(fns); i++ {
// 		cas[i].id = fns[i]
// 		cas[i].seq, cas[i].cell, cas[i].trueSS = loadFile(fns[i])
// 	}
// 	return cas
// }

// func CreateFromProteins(p []proteindb.Protein) []CellAuto {
// 	cas := make([]CellAuto, len(p))
// 	for i := 0; i < len(p); i++ {
// 		cas[i].id = p[i].Pdb_id
// 		cas[i].seq = encode(p[i].Chains[0].Seq_pdb)
// 		cas[i].cell = encode(p[i].Chains[0].Seq_pdb)
// 		cas[i].trueSS = encode(p[i].Chains[0].Ss3_cons_all)
// 	}
// 	return cas
// }

// func loadFile(fn string) (seq []byte, cell []byte, trueSS []byte) {
// 	content, err := ioutil.ReadFile("/home/jgcarvalho/sscago/data/" + fn)
// 	if err != nil {
// 		println("erro na leitura do arquivo", fn, err)
// 	}

// 	//considerando haver duas linhas, a primeira a seq e a segunda a ss
// 	lines := strings.Split(string(content), "\n")
// 	seq = encode(lines[0])
// 	cell = encode(lines[0])
// 	trueSS = encode(lines[1])
// 	return
// }

func encode(seq string, prm *rules.Params) []byte {
	out := make([]byte, len(seq))
	codes := append(prm.StrStartStates, prm.StrTransitionStates...)
	for i, s := range seq {
		for j, c := range codes {
			if string(s) == c {
				out[i] = byte(j)
			}
		}
	}
	return out
}

func decode(cell []byte) string {
	s := make([]byte, len(cell))
	for i, v := range cell {
		switch v {
		case 0:
			s[i] = '#'
		case 1:
			s[i] = '-'
		case 2:
			s[i] = '*'
		case 3:
			s[i] = '|'
		case 4:
			s[i] = 'A'
		case 5:
			s[i] = 'C'
		case 6:
			s[i] = 'D'
		case 7:
			s[i] = 'E'
		case 8:
			s[i] = 'F'
		case 9:
			s[i] = 'G'
		case 10:
			s[i] = 'H'
		case 11:
			s[i] = 'I'
		case 12:
			s[i] = 'K'
		case 13:
			s[i] = 'L'
		case 14:
			s[i] = 'M'
		case 15:
			s[i] = 'N'
		case 16:
			s[i] = 'P'
		case 17:
			s[i] = 'Q'
		case 18:
			s[i] = 'R'
		case 19:
			s[i] = 'S'
		case 20:
			s[i] = 'T'
		case 21:
			s[i] = 'V'
		case 22:
			s[i] = 'W'
		case 23:
			s[i] = 'Y'
		case 24:
			s[i] = '?'
		}
	}
	return string(s)
}
