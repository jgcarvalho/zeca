package ca

import (
	"bytes"
	"fmt"

	"github.com/jgcarvalho/zeca/rules"
)

// CellAuto1D é uma estrutura para armazenar um autômato
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

type CaConfig struct {
	InitStates  []string `toml:"initial-states"`
	TransStates []string `toml:"transition-states"`
	HasJoker    bool     `toml:"has-joker"`
	R           int      `toml:"r"`
	Steps       int      `toml:"steps"`
	Consensus   int      `toml:"consensus"`
}

// Create1D creates a 1D cellular automata with id, initial state, expected final
// state, transition rules, number of steps to evolve, and number of rows (last
// rows) to create a consensus that will be compared with the expected final state
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

// Run the 1D cellular automata to determine the final state ("end and endconsensus")
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
				// maior ou igual pq caso haja em empate, o mais próximo ao fim terá preferencia
				if bytes.Count(col, []byte{v}) >= count {
					count = bytes.Count(col, []byte{v})
					ca.EndConsensus[i] = v
					//se a contagem for maior que a metade, então ela é o consenso e podemos
					//continuar
					if count > ca.consensus/2 {
						break
					}
				}
			}
			// fmt.Println("Elements: ", col, " -> ", ca.EndConsensus[i])
		}
	}
}

// oneStep only evolves the CA for one step
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

// SetRule changes the CA rule to a new one
func (ca *CellAuto1D) SetRule(newRule *rules.Rule) {
	ca.Rule = newRule
}

// ConfusionMatrix returns a confusion matrix. The matrix has a dimension NxN
// where N is the number of transitions states. When the transition rules have a
// "wild card" the dimension of the matrix will be (N-1)x(N)
func (ca *CellAuto1D) ConfusionMatrix() [][]int {
	n := len(ca.Rule.Prm.TransitionStates)
	//np -> number of predicted; nr -> number of real
	np, nr := n, n
	if ca.Rule.Prm.Hasjoker {
		nr--
	}
	cm := make([][]int, nr)
	for i := 0; i < nr; i++ {
		cm[i] = make([]int, np)
	}
	expected, predicted := 0, 0

	// Isto pode ser otimizado para evitar a chamada dessa função "INDEXBYTE"
	for i := range ca.Expected {
		expected = bytes.IndexByte(ca.Rule.Prm.TransitionStates, ca.Expected[i])
		predicted = bytes.IndexByte(ca.Rule.Prm.TransitionStates, ca.EndConsensus[i])

		if (ca.Rule.Prm.Hasjoker && expected == len(ca.Rule.Prm.TransitionStates)-1) || expected == -1 {
			continue
		}

		if predicted == -1 {
			predicted = len(ca.Rule.Prm.TransitionStates) - 1
		}
		cm[expected][predicted]++
	}

	// fmt.Println("Confusion Matrix:")
	// fmt.Println(cm)
	return cm
}

// encode changes the sequence that represents the initial state to a slice of
// bytes
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

// decode changes a CA state (set of cells) from byte code to string code
// with the character that was set in the rule
func decode(cell []byte, prm *rules.Params) string {
	var seq string
	codes := append(prm.StrStartStates, prm.StrTransitionStates...)
	for _, c := range cell {
		for j, code := range codes {
			if int(c) == j {
				seq += code
			}
		}
	}
	return seq
}

/* old decode function kept just to precaution
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
*/
