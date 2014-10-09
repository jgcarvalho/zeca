package rules

// import (
// 	"fmt"
// 	"os"
// 	"log"
// 	"bufio"
// )

// func (r *Rules) PrintRules() {
// 	for e := 0; e < 24; e++{
// 		for c := 0; c < 24; c++ {
// 			for d := 0; d < 24; d++ {
// 				fmt.Printf("%d, %d, %d -> %d\n", e, c, d, r.Code[e][c][d])
// 			}

// 		}
// 	}
// }

// func (r *Rules) SaveToFile(fn string) {
// 	f, err := os.OpenFile(fn, os.O_RDWR|os.O_CREATE, 0666)
//     if err != nil {
//         log.Fatal(err)
//     }
//     defer f.Close()

//     b := bufio.NewWriter(f)
//     defer func() {
//         if err = b.Flush(); err != nil {
//             log.Fatal(err)
//         }
//     }()

// 	for e := 0; e < 24; e++{
// 		for c := 0; c < 24; c++ {
// 			for d := 0; d < 24; d++ {
// 				fmt.Fprintf(b, "%d, %d, %d -> %d\n", e, c, d, r.Code[e][c][d])
// 			}

// 		}
// 	}
// }
