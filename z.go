package main

import (
	// "bitbucket.org/zgcarvalho/zeca/ca"
	"bitbucket.org/zgcarvalho/zeca/cga"
	"bitbucket.org/zgcarvalho/zeca/eda"
	"bitbucket.org/zgcarvalho/zeca/sa"
	// "bitbucket.org/zgcarvalho/zeca/rules"
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"log"
	"os"
	"runtime"
)

func RunCGA(fnconfig string) {
	var conf cga.Config
	md, err := toml.DecodeFile(fnconfig, &conf)
	if err != nil {
		log.Fatal(err)
	}
	if len(md.Undecoded()) > 0 {
		fmt.Printf("Chaves desconhecidas no arquivo de configuração: %q\n", md.Undecoded())
		return
	}
	fmt.Println("Configuration:", conf)
	cga.Run(conf)
}

func RunEDA(fnconfig string) {
	var conf eda.Config
	md, err := toml.DecodeFile(fnconfig, &conf)
	if err != nil {
		log.Fatal(err)
	}
	if len(md.Undecoded()) > 0 {
		fmt.Printf("Chaves desconhecidas no arquivo de configuração: %q\n", md.Undecoded())
		return
	}
	fmt.Println("Configuration:", conf)
	eda.Run(conf)
}

func runSA(fnconfig string) {
	var conf sa.Config
	md, err := toml.DecodeFile(fnconfig, &conf)
	if err != nil {
		log.Fatal(err)
	}
	if len(md.Undecoded()) > 0 {
		fmt.Printf("Chaves desconhecidas no arquivo de configuração: %q\n", md.Undecoded())
		return
	}
	fmt.Println("Configuration:", conf)
	sa.Run(conf)
}

func printUsage() {
	fmt.Println("Manual")
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	method := flag.Int("method", 0, "Algorithm to be used during cellular automata rule search. Options: "+
		"(1) compact genetic algorithm; (2) EDA; (3) simulated annealing")
	fnconfig := flag.String("config", "default", "Configuration file")
	flag.Parse()

	switch *method {
	case 1:
		if *fnconfig == "default" {
			RunCGA(os.Getenv("GOPATH") + "/src/bitbucket.org/zgcarvalho/zeca/cgaconfig.toml")
		} else {
			RunCGA(*fnconfig)
		}
	case 2:
		if *fnconfig == "default" {
			RunEDA(os.Getenv("GOPATH") + "/src/bitbucket.org/zgcarvalho/zeca/edaconfig.toml")
		} else {
			RunEDA(*fnconfig)
		}
	case 3:
		if *fnconfig == "default" {
			runSA(os.Getenv("GOPATH") + "/src/bitbucket.org/zgcarvalho/zeca/saconfig.toml")
		} else {
			runSA(*fnconfig)
		}
	default:
		flag.PrintDefaults()
	}

}
