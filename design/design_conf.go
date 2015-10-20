package design

import (
	"bitbucket.org/jgcarvalho/zeca/ca"
	"bitbucket.org/jgcarvalho/zeca/db"
)

type Config struct {
	Title     string
	Algorithm algoConfig
	Design    designConfig
	Rules     ruleConfig
	DB        db.Config
	CA        ca.Config
}

type algoConfig struct {
	Method string `toml:"method"`
}

type designConfig struct {
	//Generations int
	Population int
	Selection  int
	//Tournament  int
	//CrossOver   int
	//Mutation    float64
	OutputProbs string `toml:"output-probabilities"`
	//SaveSteps   int    `toml:"save-steps"`
}

type ruleConfig struct {
	Input  string `toml:"input"`
	Output string `toml:"output"`
}
