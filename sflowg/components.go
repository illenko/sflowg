package sflowg

import (
	"github.com/google/uuid"
)

type Flow struct {
	ID         string         `yaml:"id"`
	Entrypoint Entrypoint     `yaml:"entrypoint"`
	Steps      []Step         `yaml:"steps"`
	Properties map[string]any `yaml:"properties"`
	Return     Return         `yaml:"return"`
}

type Entrypoint struct {
	Type   string `yaml:"type"`
	Path   string `yaml:"path"`
	Result string `yaml:"result"`
}

type Step struct {
	ID     string `yaml:"id"`
	Type   string `yaml:"type"`
	Args   any    `yaml:"args"`
	Next   string `yaml:"next,omitempty"`
	Result string `yaml:"result,omitempty"`
}

type Return struct {
	Type string         `yaml:"type"`
	Args map[string]any `yaml:"args"`
}

type Execution struct {
	ID      string
	FlowId  string
	Context map[string]any
}

func NewExecution(flow Flow) Execution {
	return Execution{
		ID:      uuid.New().String(),
		FlowId:  flow.ID,
		Context: make(map[string]any),
	}
}

func (e *Execution) AddContext(k string, v any) {
	e.Context[FormatKey(k)] = v
}
