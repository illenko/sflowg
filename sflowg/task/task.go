package task

import (
	"sflowg/sflowg"
	"sflowg/sflowg/container"
)

type Task interface {
	Execute(*container.Container, *sflowg.Execution, map[string]any) (map[string]any, error)
}
