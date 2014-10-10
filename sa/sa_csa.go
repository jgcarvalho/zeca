package sa

// import (
// 	"github.com/jgcarvalho/zeca/ca"
// 	"github.com/jgcarvalho/zeca/rules"
// 	"fmt"
// 	"math"
// 	"math/rand"
// 	"time"
// )

/*
func RunSA(c ca.CellAuto, rule *rules.Rules, t0 float64, n int, alpha float64) {
	_, c3, _, _, _,_,_,_,_ := ca.Run(c,rule)
	t := t0
	nSucesso := 0

	ruleNew := new(rules.Rules)
	c3New := 0.0
	deltac3 := 0.0

	rand.Seed(time.Now().UnixNano())

	for i := 0; i < n; i++ {
		ruleNew = rules.PertRules(rule, t)
		_, c3New, _, _, _,_,_,_,_ = ca.Run(c,ruleNew)
		deltac3 = c3New - c3

		if deltac3 > 0.0 {
			rule.Code = ruleNew.Code
			c3 = c3New
			nSucesso += 1
		} else if math.Exp(deltac3/t) > rand.Float64() {
			rule.Code = ruleNew.Code
			c3 = c3New
			nSucesso += 1
		}

		if deltac3 == 1.0 {
			break
		}
		if i%1000 == 0 {
			fmt.Println(i, nSucesso, t, c3, deltac3)
		}

		t = t*alpha
	}

	fmt.Println(ca.Run(c,rule))
}

func RunSAN(cas []ca.CellAuto, rule *rules.Rules, t0 float64, n int, alpha float64) {
	t := t0
	nSucesso := 0

	ruleNew := new(rules.Rules)
	sumc3 := 0.0
	sumc3New := 0.0
	deltaSumc3 := 0.0

	for i := 0; i < len(cas); i++ {
		_, c3, _, _, _,_,_,_,_ := ca.Run(cas[i],rule)
		sumc3 += c3
	}

	rand.Seed(time.Now().UnixNano())

	for i := 0; i < n; i++ {
		ruleNew = rules.PertRules(rule, t)

		sumc3New = 0.0
		for i := 0; i < len(cas); i++ {
			_, c3New, _, _, _ ,_,_,_,_ := ca.Run(cas[i],ruleNew)
			sumc3New += c3New
		}

		deltaSumc3 = sumc3New/float64(len(cas)) - sumc3/float64(len(cas))

		if deltaSumc3 > 0.0 {
			rule.Code = ruleNew.Code
			sumc3 = sumc3New
			nSucesso += 1
		} else if math.Exp(deltaSumc3/t) > rand.Float64() {
			rule.Code = ruleNew.Code
			sumc3 = sumc3New
			nSucesso += 1
		}

		if deltaSumc3 == 1.0 {
			break
		}

		if i%1000 == 0 {
			fmt.Println(i, nSucesso, t, sumc3/float64(len(cas)), deltaSumc3)
		}

		t = t*alpha
	}

	for i := 0; i < len(cas); i++ {
		fmt.Println(ca.Run(cas[i],rule))
	}
}
*/

// Comentarios abaixo Apenas para testar

// func RunSANp(cas []ca.CellAuto, rule *rules.Rules, t0 float64, tn float64, n int) *rules.Rules{
// 	n += 1
// 	t := t0
// 	alpha := math.Pow((tn/t0), 1.0/float64(n))
// 	nSucesso := 0

// 	ruleBest := new(rules.Rules)
// 	avgc3Best := -1.0
// 	avgq3Best := 0.0

// 	ruleNew := new(rules.Rules)
// 	sumc3 := 0.0
// 	sumc3New := 0.0
// 	avgc3 := 0.0
// 	avgc3New := 0.0
// 	deltaAvgc3 := 0.0

// 	sumq3 := 0.0
// 	avgq3 := 0.0

// 	chc3 := make(chan float64, len(cas))
// 	chq3 := make(chan float64, len(cas))

// 	for i := 0; i < len(cas); i++ {
// 		_, c3, _, _, _,_,_,_,_ := ca.Run(cas[i],rule)
// 		sumc3 += c3
// 	}
// 	avgc3 = sumc3/float64(len(cas))

// 	rand.Seed(time.Now().UnixNano())

// 	for i := 0; i < n; i++ {
// 		ruleNew = rules.PertRules(rule, t)

// 		sumc3New = 0.0
// 		sumq3 = 0.0
// 		for i := 0; i < len(cas); i++ {
// 			go ca.RunCh(cas[i],ruleNew, chc3, chq3)
// 		}

// 		for i := 0; i < len(cas); i++ {
// 			sumc3New += <- chc3
// 			sumq3 += <- chq3
// 		}

// 		avgc3New = sumc3New/float64(len(cas))
// 		//fmt.Println("sumq3",sumq3)
// 		//fmt.Println("len(cas)",float64(len(cas)))
// 		avgq3 = sumq3/float64(len(cas))

// 		deltaAvgc3 = avgc3New  - avgc3

// 		if deltaAvgc3 > 0.0 || math.Exp(deltaAvgc3/t) > rand.Float64(){
// 			rule.Code = ruleNew.Code
// 			avgc3 = avgc3New
// 			nSucesso += 1
// 			if avgc3Best < avgc3 {
// 				ruleBest.Code = rule.Code
// 				avgc3Best = avgc3
// 				avgq3Best = avgq3
// 			}
// 		}

// 		if avgc3 == 1.0 {
// 			break
// 		}
// 		//fmt.Println(i, nSucesso, t, avgc3, deltaAvgc3, avgq3)
// 		if i%1000 == 0 {
// 			fmt.Println(i, nSucesso, t, avgc3, deltaAvgc3, avgq3)
// 		}

// 		t = t*alpha
// 	}

// 	fmt.Printf("Best Rule Found: C3 = %v, Q3 = %v\n",avgc3Best, avgq3Best)
// 	ruleBest.SaveToFile("./data/rule")

// 	// for i := 0; i < len(cas); i++ {
// 	// 	fmt.Println(ca.Run(cas[i],ruleBest))
// 	// }

// 	return ruleBest

// }
