package sflowg

import (
	"context"

	"github.com/google/uuid"
)

type Execution struct {
	ID        string
	Context   context.Context
	Values    map[string]any
	Flow      *Flow
	Container *Container
}

func NewExecution(flow *Flow, container *Container) Execution {
	id := uuid.New().String()
	ctx := context.Background()
	ctx = context.WithValue(ctx, "flowId", flow.ID)
	ctx = context.WithValue(ctx, "executionId", id)
	return Execution{
		ID:        id,
		Context:   ctx,
		Values:    make(map[string]any),
		Flow:      flow,
		Container: container,
	}
}

func (e *Execution) AddVal(k string, v any) {
	e.Values[FormatKey(k)] = v
}
