package main

import (
	// "github.com/jgcarvalho/zeca/ca"
	"github.com/jgcarvalho/zeca/cga"
	"github.com/jgcarvalho/zeca/design"
	"github.com/jgcarvalho/zeca/eda"
	"github.com/jgcarvalho/zeca/ga"
	"github.com/jgcarvalho/zeca/sa"
	// "github.com/jgcarvalho/zeca/rules"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/BurntSushi/toml"
)

func runCGA(fnconfig string) {
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

func runGA(fnconfig string) {
	var conf ga.Config
	md, err := toml.DecodeFile(fnconfig, &conf)
	if err != nil {
		log.Fatal(err)
	}
	if len(md.Undecoded()) > 0 {
		fmt.Printf("Chaves desconhecidas no arquivo de configuração: %q\n", md.Undecoded())
		return
	}
	fmt.Println("Configuration:", conf)
	ga.Run(conf)
}

func runEDA(fnconfig string) {
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

func runDesign(fnconfig string) {
	var conf design.Config
	md, err := toml.DecodeFile(fnconfig, &conf)
	if err != nil {
		log.Fatal(err)
	}
	if len(md.Undecoded()) > 0 {
		fmt.Printf("Chaves desconhecidas no arquivo de configuração: %q\n", md.Undecoded())
		return
	}
	fmt.Println("Configuration:", conf)
	design.Run(conf)
}

func printUsage() {
	fmt.Println("Manual")
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	method := flag.Int("method", 0, "Algorithm to be used during cellular automata rule search. Options: "+
		"(1) compact genetic algorithm; (2) EDA; (3) simulated annealing")
	fnconfig := flag.String("config", "default", "Configuration file")
	cpuprofile := flag.String("cpuprofile", "", "write cpu profile to file")
	memprofile := flag.String("memprofile", "", "write memory profile to this file")
	flag.Parse()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.WriteHeapProfile(f)
		defer f.Close()
	}

	switch *method {
	case 1:
		if *fnconfig == "default" {
			runCGA(os.Getenv("GOPATH") + "/src/github.com/jgcarvalho/zeca/cgaconfig.toml")
		} else {
			runCGA(*fnconfig)
		}
	case 2:
		if *fnconfig == "default" {
			runEDA(os.Getenv("GOPATH") + "/src/github.com/jgcarvalho/zeca/edaconfig.toml")
		} else {
			runEDA(*fnconfig)
		}
	case 3:
		if *fnconfig == "default" {
			runSA(os.Getenv("GOPATH") + "/src/github.com/jgcarvalho/zeca/saconfig.toml")
		} else {
			runSA(*fnconfig)
		}
	case 4:
		if *fnconfig == "default" {
			runGA(os.Getenv("GOPATH") + "/src/github.com/jgcarvalho/zeca/gaconfig.toml")
		} else {
			runGA(*fnconfig)
		}
	case 9:
		if *fnconfig == "default" {
			runDesign(os.Getenv("GOPATH") + "/src/github.com/jgcarvalho/zeca/designconfig.toml")
		} else {
			runDesign(*fnconfig)
		}
	default:
		flag.PrintDefaults()
	}

}
