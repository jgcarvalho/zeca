package dist

import (
	"github.com/jgcarvalho/zeca/ca"
	"github.com/jgcarvalho/zeca/proteindb"
)

type Config struct {
	Title     string
	Algorithm algoConfig
	EDA       edaConfig
	Rules     ruleConfig
	ProteinDB proteindb.ProtdbConfig
	CA        ca.Config
	Dist      distConfig
}

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
	NatsServerURL   string `toml:"nats-server"`
	TopicFromMaster string `toml:"topic-master"`
	TopicFromSlave  string `toml:"topic-slave"`
}
