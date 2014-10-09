package ca

// import (
// 	"bitbucket.org/zgcarvalho/zeca/rules"
// 	"bytes"
// 	"github.com/ajstarks/svgo"
//     "os"
//     "log"
// 	"bufio"
// 	"fmt"
// )

// func oneStep(seq []byte, currentState []byte, nextState []byte, rule *rules.Rules){
// 	var state byte
// 	for c := 1; c < len(currentState) - 1 ; c++ {
// 		state = rule.Code[currentState[c-1]][currentState[c]][currentState[c+1]]
// 		if state == 24 {
// 			nextState[c] = seq[c]
// 		} else {
// 			nextState[c] = state
// 		}
// 	}
// }

// /* Talvez utilizar um worker para paralelizar
// func workerOneStep(seq []byte, currentState []byte, nextState []byte, rule *rules.Rules){
// 	var state byte
// 	for c := 1; c < len(currentState) - 1 ; c++ {
// 		state = rule.Code[currentState[c-1]][currentState[c]][currentState[c+1]]
// 		if state == 24 {
// 			nextState[c] = seq[c]
// 		} else {
// 			nextState[c] = state
// 		}
// 	}
// }
// */

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
// 	//fmt.Println(currentState)
// 	c3, cL, cH, cS, q3, qL, qH, qS := fitness(currentState, c.trueSS)
// 	return currentState, c3, cL, cH, cS, q3, qL, qH, qS
// }

// func RunCh(c CellAuto, rule *rules.Rules, ch1 chan float64, ch2 chan float64) {
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
// 	//fmt.Println(currentState)
// 	c3,q3 := fitnessSimple(currentState, c.trueSS)
// 	//fmt.Println("pdbid", c.id)
// 	//fmt.Println("q3", q3)
// 	ch1 <- c3
// 	ch2 <- q3
// }

// func RunToView(c CellAuto, rule *rules.Rules) ([]byte, float64, float64, float64, float64, float64, float64, float64, float64) {
// 	currentState := make([]byte, len(c.cell))
// 	nextState := make([]byte, len(c.cell))
// 	copy(currentState, c.cell)
// 	copy(nextState, c.cell)

// 	f, err := os.OpenFile("./data/ssca_"+c.id+".svg", os.O_RDWR|os.O_CREATE, 0666)
// 	if err != nil {
// 	    log.Fatal(err)
// 	}
// 	defer f.Close()

// 	b := bufio.NewWriter(f)
// 	defer func() {
// 	    if err = b.Flush(); err != nil {
// 	        log.Fatal(err)
// 	    }
// 	}()

// 	width := 30 * len(c.cell)
// 	height := len(c.cell)
// 	canvas := svg.New(b)
// 	canvas.Start(width, height)
// 	x := 30
// 	y := 60
// 	canvas.Text(x + 30/2, y + 25, "PDB id: "+c.id, "font-size:40px;fill:black")
// 	x = 0
// 	y += 50

// 	state := decode(currentState)
// 	for _, v := range state {
// 		switch string(v) {
// 		case "#":
// 			canvas.Square(x, y, 30, "fill:black")
// 		case "-":
// 			canvas.Square(x, y, 30, "fill:yellow")
// 		case "*":
// 			canvas.Square(x, y, 30, "fill:red")
// 		case "|":
// 			canvas.Square(x, y, 30, "fill:blue")
// 		default:
// 			canvas.Square(x, y, 30, "fill:white")
// 			canvas.Text(x + 30/2, y + 25, string(v), "text-anchor:middle;font-size:20px;fill:black")
// 		}

// 		x += 30
// 	}

// 	//fmt.Println(currentState)
// 	for i := 0; i < len(currentState) * 2; i++ {
// 		if i%2 == 0 {
// 			y += 30
// 			x = 0
// 			oneStep(c.seq, currentState, nextState, rule)

// 			state := decode(nextState)
// 			for _, v := range state {
// 				switch string(v) {
// 				case "#":
// 					canvas.Square(x, y, 30, "fill:black")
// 				case "-":
// 					canvas.Square(x, y, 30, "fill:yellow")
// 				case "*":
// 					canvas.Square(x, y, 30, "fill:red")
// 				case "|":
// 					canvas.Square(x, y, 30, "fill:blue")
// 				default:
// 					canvas.Square(x, y, 30, "fill:white")
// 					canvas.Text(x + 30/2, y + 25, string(v), "text-anchor:middle;font-size:20px;fill:black")
// 				}
// 				x += 30
// 			}

// 			//fmt.Println(i, nextState)
// 		} else {
// 			y += 30
// 			x = 0
// 			oneStep(c.seq, nextState, currentState, rule)

// 			state := decode(currentState)
// 			for _, v := range state {
// 				switch string(v) {
// 				case "#":
// 					canvas.Square(x, y, 30, "fill:black")
// 				case "-":
// 					canvas.Square(x, y, 30, "fill:yellow")
// 				case "*":
// 					canvas.Square(x, y, 30, "fill:red")
// 				case "|":
// 					canvas.Square(x, y, 30, "fill:blue")
// 				default:
// 					canvas.Square(x, y, 30, "fill:white")
// 					canvas.Text(x + 30/2, y + 25, string(v), "text-anchor:middle;font-size:20px;fill:black")
// 				}
// 				x += 30
// 			//fmt.Println(i, currentState)
// 			}
// 		}

// 		if bytes.Equal(currentState,nextState) {
// 			//fmt.Println("Estabilizou no passo", i)
// 			break
// 		}
// 	}

// 	y += 70
// 	x = 30
// 	canvas.Text(x + 30/2, y + 25, "Secondary Structure Assigned", "font-size:40px;fill:black")
// 	y += 50
// 	x = 0

// 	trueSS := decode(c.trueSS)
// 	for _, v := range trueSS {
// 		switch string(v) {
// 		case "#":
// 			canvas.Square(x, y, 30, "fill:black")
// 		case "-":
// 			canvas.Square(x, y, 30, "fill:yellow")
// 		case "*":
// 			canvas.Square(x, y, 30, "fill:red")
// 		case "|":
// 			canvas.Square(x, y, 30, "fill:blue")
// 		default:
// 			canvas.Square(x, y, 30, "fill:white")
// 			canvas.Text(x + 30/2, y + 25, string(v), "text-anchor:middle;font-size:20px;fill:black")
// 		}

// 		x += 30
// 	}

// 	y += 70
// 	x = 30

// 	//fmt.Println(currentState)
// 	c3, cL, cH, cS, q3, qL, qH, qS := fitness(currentState, c.trueSS)
// 	//canvas.Text(x + 30/2, y + 25, fmt.Sprintf("C3= %3.f, Cl= %3.f, Ch= %3.f, Cs= %3.f, Q3= %3.f, Ql= %3.f, Qh= %3.f, Qs= %3.f" , "font-size:30px;fill:black")
// 	canvas.Text(x + 30/2, y + 25, fmt.Sprintf("C3= %.3f, Cl= %.3f, Ch= %.3f, Cs= %.3f, Q3= %.3f, Ql= %.3f, Qh= %.3f, Qs= %.3f", c3, cL, cH, cS, q3, qL, qH, qS) , "font-size:30px;fill:black")
// 	canvas.End()
// 	return currentState, c3, cL, cH, cS, q3, qL, qH, qS
// }
