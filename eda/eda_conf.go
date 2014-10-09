package eda

type Config struct {
	Title     string
	Algorithm algoConfig
	EDA       edaConfig
	Rules     ruleConfig
	ProteinDB protdbConfig
	CA        caConfig
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

type protdbConfig struct {
	Ip         string `toml:"db-ip"`
	Name       string `toml:"db-name"`
	Collection string `toml:"collection-name"`
	Target     string `toml:"target"`
}

type caConfig struct {
	InitStates  []string `toml:"initial-states"`
	TransStates []string `toml:"transition-states"`
	HasJoker    bool     `toml:"has-joker"`
	R           int      `toml:"r"`
	Steps       int      `toml:"steps"`
	Consensus   int      `toml:"consensus"`
}
