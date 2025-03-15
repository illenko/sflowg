package sflowg

type Task interface {
	Execute(*Execution, map[string]any) (map[string]any, error)
}
