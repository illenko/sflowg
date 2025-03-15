package sflowg

type Flow struct {
	ID         string         `yaml:"id"`
	Entrypoint Entrypoint     `yaml:"entrypoint"`
	Steps      []Step         `yaml:"steps"`
	Properties map[string]any `yaml:"properties"`
	Return     Return         `yaml:"return"`
}

type Entrypoint struct {
	Type   string         `yaml:"type"`
	Config map[string]any `yaml:"config"`
}

type Step struct {
	ID        string         `yaml:"id"`
	Type      string         `yaml:"type"`
	Condition string         `yaml:"condition,omitempty"`
	Args      map[string]any `yaml:"args"`
	Next      string         `yaml:"next,omitempty"`
	Retry     *RetryConfig   `yaml:"retry,omitempty"`
}

type Return struct {
	Type string         `yaml:"type"`
	Args map[string]any `yaml:"args"`
}

type RetryConfig struct {
	MaxRetries int    `yaml:"maxRetries"`
	Delay      int    `yaml:"delay"`
	Backoff    bool   `yaml:"backoff"`
	Condition  string `yaml:"condition"`
}
