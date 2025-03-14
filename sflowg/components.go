package sflowg

import (
	"context"

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
	Type   string         `yaml:"type"`
	Config map[string]any `yaml:"config"`
	Result string         `yaml:"result"`
}

type Step struct {
	ID        string `yaml:"id"`
	Type      string `yaml:"type"`
	Condition string `yaml:"condition,omitempty"`
	Args      any    `yaml:"args"`
	Next      string `yaml:"next,omitempty"`
	Result    string `yaml:"result,omitempty"`
}

type Return struct {
	Type string         `yaml:"type"`
	Args map[string]any `yaml:"args"`
}

type Execution struct {
	ID      string
	FlowId  string
	Context context.Context
	Values  map[string]any
}

func NewExecution(flow Flow) Execution {
	id := uuid.New().String()
	ctx := context.Background()
	ctx = context.WithValue(ctx, "flowId", flow.ID)
	ctx = context.WithValue(ctx, "executionId", id)
	return Execution{
		ID:      id,
		FlowId:  flow.ID,
		Context: ctx,
		Values:  make(map[string]any),
	}
}

func (e *Execution) AddVal(k string, v any) {
	e.Values[FormatKey(k)] = v
}
