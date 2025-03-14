package task

import (
	"sflowg/sflowg"
)

type Task interface {
	Execute(*sflowg.Container, *sflowg.Execution, map[string]any) (map[string]any, error)
}
