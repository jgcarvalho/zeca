package rules

// import (
// 	"math/rand"
// 	"time"
// )

// // type Rules struct {
// // 	Code [24][24][24]byte
// // }

// func GenRules() *Rules {
// 	rule := new(Rules)

// 	rand.Seed(time.Now().UnixNano())

// 	for e := 0; e < 24; e++ {
// 		for c := 0; c < 24; c++ {
// 			for d := 0; d < 24; d++ {
// 				if c == 0 {
// 					rule.Code[e][c][d] = 0
// 				} else if (e == 0 || e == 1 || e == 2 || e == 3) && (c == 1 || c == 2 || c == 3) && (d == 0 || d == 1 || d == 2 || d == 3) {
// 					rule.Code[e][c][d] = byte(c)
// 				} else {
// 					if rand.Float64() < 0.3 {
// 						randChoice := rand.Intn(3)
// 						rule.Code[e][c][d] = byte(randChoice + 1)
// 					} else {
// 						rule.Code[e][c][d] = 24
// 					}

// 				}
// 			}
// 		}
// 	}
// 	return rule
// }

// func PertRules(ruleOld *Rules, t float64) *Rules {
// 	rule := new(Rules)
// 	rule.Code = ruleOld.Code

// 	rand.Seed(time.Now().UnixNano())

// 	changes := 0

// 	for e := 0; e < 24; e++ {
// 		for c := 0; c < 24; c++ {
// 			for d := 0; d < 24; d++ {
// 				if c == 0 {
// 					continue
// 				} else if (e == 0 || e == 1 || e == 2 || e == 3) && (c == 1 || c == 2 || c == 3) && (d == 0 || d == 1 || d == 2 || d == 3) {
// 					continue
// 				} else {
// 					if rand.Float64() < 0.0005 {
// 						changes += 1
// 						randChoice := rand.Intn(4)
// 						if randChoice == 0 {
// 							rule.Code[e][c][d] = 24
// 						} else {
// 							rule.Code[e][c][d] = byte(randChoice)
// 						}
// 					}
// 				}
// 			}
// 		}
// 	}
// 	//println("Changes = ", changes)
// 	return rule
// }
