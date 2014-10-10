package main

import (
	// "github.com/jgcarvalho/zeca/ca"
	"github.com/jgcarvalho/zeca/cga"
	"github.com/jgcarvalho/zeca/eda"
	"github.com/jgcarvalho/zeca/sa"
	// "github.com/jgcarvalho/zeca/rules"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/BurntSushi/toml"
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
			RunCGA(os.Getenv("GOPATH") + "/src/github.com/jgcarvalho/zeca/cgaconfig.toml")
		} else {
			RunCGA(*fnconfig)
		}
	case 2:
		if *fnconfig == "default" {
			RunEDA(os.Getenv("GOPATH") + "/src/github.com/jgcarvalho/zeca/edaconfig.toml")
		} else {
			RunEDA(*fnconfig)
		}
	case 3:
		if *fnconfig == "default" {
			runSA(os.Getenv("GOPATH") + "/src/github.com/jgcarvalho/zeca/saconfig.toml")
		} else {
			runSA(*fnconfig)
		}
	default:
		flag.PrintDefaults()
	}

}
