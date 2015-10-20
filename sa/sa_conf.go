package sa

import (
	"bitbucket.org/jgcarvalho/zeca/ca"
	"bitbucket.org/jgcarvalho/zeca/proteindb"
)

type Config struct {
	Title     string
	Algorithm algoConfig
	SA        saConfig
	Rules     ruleConfig
	ProteinDB proteindb.ProtdbConfig
	CA        ca.Config
}

type algoConfig struct {
	Method string `toml:"method"`
}

type saConfig struct {
	OuterLoop int     `toml:"outer-loop"`
	InnerLoop int     `toml:"inner-loop"`
	Tini      float64 `toml:"temp-start"`
	Tfinal    float64 `toml:"temp-final"`
	SaveSteps int     `toml:"save-steps"`
}

type ruleConfig struct {
	Input  string `toml:"input"`
	Output string `toml:"output"`
}
