package sa

type Config struct {
	Title     string
	Algorithm algoConfig
	SA        saConfig
	Rules     ruleConfig
	ProteinDB protdbConfig
	CA        caConfig
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
