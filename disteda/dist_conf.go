package disteda

import (
	"github.com/jgcarvalho/zeca/ca"
	"github.com/jgcarvalho/zeca/db"
)

type Config struct {
	Title     string
	Algorithm algoConfig
	EDA       edaConfig
	Rules     ruleConfig
	DB        db.Config
	CA        ca.Config
	Dist      distConfig
}

// type dbConfig struct {
//
// }

type algoConfig struct {
	Method string `toml:"method"`
}

type edaConfig struct {
	Generations int
	Population  int
	Selection   int
	Tournament  int
	OutputProbs string `toml:"output-probabilities"`
	SaveSteps   int    `toml:"save-steps"`
}

type ruleConfig struct {
	Input  string `toml:"input"`
	Output string `toml:"output"`
}

type distConfig struct {
	MasterURL string `toml:"master-url"`
	PortA     string `toml:"port-a"`
	PortB     string `toml:"port-b"`
}
