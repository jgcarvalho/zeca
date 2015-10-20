package ga

import (
	"bitbucket.org/jgcarvalho/zeca/ca"
	"bitbucket.org/jgcarvalho/zeca/proteindb"
)

type Config struct {
	Title     string
	Algorithm algoConfig
	GA        gaConfig
	Rules     ruleConfig
	ProteinDB proteindb.ProtdbConfig
	CA        ca.Config
}

type algoConfig struct {
	Method string `toml:"method"`
}

type gaConfig struct {
	Generations int
	Population  int
	Selection   int
	Tournament  int
	CrossOver   int
	Mutation    float64
	OutputProbs string `toml:"output-probabilities"`
	SaveSteps   int    `toml:"save-steps"`
}

type ruleConfig struct {
	Input  string `toml:"input"`
	Output string `toml:"output"`
}
