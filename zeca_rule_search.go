package main

// import (
// 	"bitbucket.org/zgcarvalho/zeca/ca"
// 	"bitbucket.org/zgcarvalho/zeca/proteindb"
// 	"bitbucket.org/zgcarvalho/zeca/rules"
// 	"bitbucket.org/zgcarvalho/zeca/sa"
// 	"fmt"
// 	"io/ioutil"
// 	"runtime"
// 	"strings"
// 	// _ "net/http/pprof"
// 	// "net/http"
// 	// "log"
// )

// func loadPDBs(fn string) (pdbs []string) {
// 	/* Funcao que le o nome dos arquivos (PDB ID) contendo a sequencia e a estrutura secundaria*/
// 	content, err := ioutil.ReadFile(fn)
// 	if err != nil {
// 		println("Erro na leitura do arquivo que contem nomes de arquivos com a sequencia e a estrutura secundaria", fn, err)
// 	}

// 	//retorna uma lista com os codigos de PDB (nome dos arquivos) #TODO: resolver como inserir o diretorio
// 	pdbs = strings.Split(string(content), "\n")
// 	//Remove o ultimo elemento que e uma string vazia #TODO melhorar essa gambiarra
// 	pdbs = pdbs[:len(pdbs)-1]
// 	return
// }

// func main() {
// 	// go func() {
// 	// 	log.Println(http.ListenAndServe("localhost:6060", nil))
// 	// }()

// 	runtime.GOMAXPROCS(runtime.NumCPU())
// 	//Gera uma regra aleatoria (inicial) para o automato celular
// 	rule := rules.GenRules()

// 	//Le o nome dos arquivos de entrada
// 	//pdbs := loadPDBs("/home/jgcarvalho/sscago/data/pdb_list")

// 	//DEBUG: Imprime o nome dos arquivos de entrada e o numero de arquivos
// 	//fmt.Println(pdbs)
// 	//fmt.Println(len(pdbs))

// 	//Cria N automatos celulares de acordo com os arquivos passados como entrada
// 	//cas := ca.CreateN(pdbs)
// 	proteins := proteindb.LoadProteinsFromMongo("proteindb_dev", "protein")
// 	cas := ca.CreateFromProteins(proteins)
// 	fmt.Println("Estão sendo usadas", len(proteins), "estruturas proteicas.")

// 	//c := ca.CreateOne("./4hti")

// 	//Roda N simulated anealing para otimizar a regra (automatos, regra, temperatura inicial, temperatura final, passos)
// 	ruleBest := sa.RunSANp(cas, rule, 1.0, 0.0000000001, 1000)
// 	//rule.PrintRules()

// 	for i := 0; i < len(cas); i++ {
// 		ca.RunToView(cas[i], ruleBest)
// 	}
// }
