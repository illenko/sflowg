package sflowg

import (
	"context"
	"time"

	"github.com/google/uuid"
)

var _ context.Context = &Execution{}

type Execution struct {
	ID        string
	Values    map[string]any
	Flow      *Flow
	Container *Container
}

func (e *Execution) Deadline() (deadline time.Time, ok bool) {
	return time.Time{}, false
}

func (e *Execution) Done() <-chan struct{} {
	return nil
}

func (e *Execution) Err() error {
	if e.Container == nil {
		return nil
	}

	return nil
}

func (e *Execution) AddValue(k string, v any) {
	e.Values[FormatKey(k)] = v
}

func (e *Execution) Value(key any) any {
	k, ok := key.(string)

	if !ok {
		return nil
	}

	return e.Values[FormatKey(k)]
}

func NewExecution(flow *Flow, container *Container) Execution {
	id := uuid.New().String()
	return Execution{
		ID:        id,
		Values:    make(map[string]any),
		Flow:      flow,
		Container: container,
	}
}
