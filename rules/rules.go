package rules

import (
	"fmt"
	"math/rand"
)

// Rule represents a CA rule
type Rule struct {
	Code  [][][]byte
	Fixed [][][]bool
	Prm   Params
}

// Params represents the parameters of the CA rule, with information about start
// and transition rules, if there is a "wild card" in the transition states, and
// the neighborhood size R
type Params struct {
	StrStartStates      []string
	StrTransitionStates []string
	StartStates         []byte
	TransitionStates    []byte
	Hasjoker            bool
	R                   uint8
}

//PrmDefault was used just for test. I'm not sure! :|
var PrmDefault Params

func init() {
	PrmDefault.StartStates = []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	PrmDefault.TransitionStates = []byte{21, 22, 23, 24}
	PrmDefault.Hasjoker = true
	PrmDefault.R = 3
}

// RuleStates recebe os parâmetros para a criação da regra e calcula quantos e quais estados
// tem que estar presentes na regra. Esses estados são formados da união do conjunto de estados
// de início (StartStates) com o conjunto de estados de transição (TransitionStates)
//      [ls][s][rs]
//			 ↓
//		 	[t]
func RuleStates(prm Params) []byte {
	st := make([]byte, len(prm.StartStates), len(prm.StartStates)+len(prm.TransitionStates))
	copy(st, prm.StartStates)
	for i, vi := range prm.TransitionStates {
		if prm.Hasjoker && i == (len(prm.TransitionStates)-1) {
			break
		}
		isIn := false
		for _, vj := range prm.StartStates {
			if vi == vj {
				isIn = true
			}
		}
		if !isIn {
			st = append(st, vi)
		}
	}
	return st
}

// Create a new rule given the start states, transition states, if there is a joker
// (which must be the last element in the transition states) and neighborhood r.
func Create(sStates []string, tStates []string, hasjoker bool, r int) (*Rule, error) {
	var ru Rule

	ru.Prm.StrStartStates = sStates
	ru.Prm.StrTransitionStates = tStates
	ru.Prm.StartStates = make([]byte, len(sStates))
	for i := 0; i < len(sStates); i++ {
		ru.Prm.StartStates[i] = byte(i)
	}
	ru.Prm.TransitionStates = make([]byte, len(tStates))
	for i := 0; i < len(tStates); i++ {
		ru.Prm.TransitionStates[i] = byte(len(sStates) + i)
	}
	ru.Prm.Hasjoker = hasjoker
	ru.Prm.R = uint8(r)

	//all possible states
	st := RuleStates(ru.Prm)
	ru.Code = make([][][]byte, len(st))
	ru.Fixed = make([][][]bool, len(st))
	for ln := range st {
		ru.Code[ln] = make([][]byte, len(st))
		ru.Fixed[ln] = make([][]bool, len(st))
		for c := range st {
			ru.Code[ln][c] = make([]byte, len(st))
			ru.Fixed[ln][c] = make([]bool, len(st))
			for rn := range st {
				ru.Code[ln][c][rn] = ru.Prm.TransitionStates[rand.Intn(len(ru.Prm.TransitionStates))]
				// c == 0 usually means # that represents the n and c terminal. it's not a residue
				// but it's essencial to CA representation and must be fixed (never change)
				if c == 0 {
					ru.Code[ln][c][rn] = 0
					ru.Fixed[ln][c][rn] = true
				} else {
					ru.Fixed[ln][c][rn] = false
				}
			}
		}

	}
	return &ru, nil
}

func (r *Rule) String() string {
	codes := append(r.Prm.StrStartStates, r.Prm.StrTransitionStates...)
	var toprint string
	for c := 0; c < len(r.Code); c++ {
		for ln := 0; ln < len(r.Code); ln++ {
			for rn := 0; rn < len(r.Code); rn++ {
				toprint += fmt.Sprintf("[%s][%s][%s] -> [%s]\n", codes[ln], codes[c], codes[rn], codes[r.Code[ln][c][rn]])
			}
		}
	}
	return toprint
}
