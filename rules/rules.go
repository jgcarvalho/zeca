package rules

import (
	"fmt"
	"math/rand"
)

type Params struct {
	StrStartStates      []string
	StrTransitionStates []string
	StartStates         []byte
	TransitionStates    []byte
	Hasjoker            bool
	R                   uint8
}

type Rule struct {
	Code  [][][]byte
	Fixed [][][]bool
	Prm   Params
}

var PrmDefault Params

func init() {
	PrmDefault.StartStates = []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	PrmDefault.TransitionStates = []byte{21, 22, 23, 24}
	PrmDefault.Hasjoker = true
	PrmDefault.R = 3
}

//	ruleStates recebe os parâmetros para a criação da regra e calcula quantos e quais estados
//	tem que estar presentes na regra. Esses estados são formados da união do conjunto de estados
//	de início (StartStates) com o conjunto de estados de transição (TransitionStates)
//		[ls][s][rs]
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

	st := RuleStates(ru.Prm)
	ru.Code = make([][][]byte, len(st))
	ru.Fixed = make([][][]bool, len(st))
	for c := range st {
		ru.Code[c] = make([][]byte, len(st))
		ru.Fixed[c] = make([][]bool, len(st))
		for ln := range st {
			ru.Code[c][ln] = make([]byte, len(st))
			ru.Fixed[c][ln] = make([]bool, len(st))
			for rn := range st {
				ru.Code[c][ln][rn] = ru.Prm.TransitionStates[rand.Intn(len(ru.Prm.TransitionStates))]
				//
				if c == 0 {
					ru.Fixed[c][ln][rn] = true
				} else {
					ru.Fixed[c][ln][rn] = false
				}
			}
		}

	}
	// fmt.Println("states", st)
	// fmt.Println("rules", rule.Code)
	// fmt.Println(ru)
	return &ru, nil
}

func (r *Rule) String() string {
	codes := append(r.Prm.StrStartStates, r.Prm.StrTransitionStates...)
	var toprint string
	for c := 0; c < len(r.Code); c++ {
		for ln := 0; ln < len(r.Code); ln++ {
			for rn := 0; rn < len(r.Code); rn++ {
				toprint += fmt.Sprintf("[%s][%s][%s] -> [%s]\n", codes[ln], codes[c], codes[rn], codes[r.Code[c][ln][rn]])
			}
		}
	}
	return toprint
}
