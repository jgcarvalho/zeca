package cga

import (
	"github.com/jgcarvalho/zeca/ca"
	"github.com/jgcarvalho/zeca/db"
)

type Config struct {
	Title     string
	Algorithm algoConfig
	CGA       cgaConfig
	Rules     ruleConfig
	DB        db.Config
	CA        ca.Config
}

type algoConfig struct {
	Method string `toml:"method"`
}

type cgaConfig struct {
	Generations int
	Population  int
	Selection   int
	OutputProbs string `toml:"output-probabilities"`
	SaveSteps   int    `toml:"save-steps"`
}

type ruleConfig struct {
	Input  string `toml:"input"`
	Output string `toml:"output"`
}
