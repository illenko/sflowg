package sflowg

import (
	"fmt"
	"log/slog"
	"strconv"
	"time"
)

type Executor struct {
	l *slog.Logger
}

func NewExecutor(l *slog.Logger) *Executor {
	return &Executor{l: l}
}

func (e *Executor) ExecuteSteps(execution *Execution) error {
	nextStep := ""

	for _, s := range execution.Flow.Steps {
		if nextStep != "" {
			if s.ID != nextStep {
				e.l.InfoContext(execution, fmt.Sprintf("Skipping step: %s", s.ID))
				continue
			}
			nextStep = ""
			e.l.InfoContext(execution, fmt.Sprintf("Resuming flow at step: %s", s.ID))
		}

		if err := e.evaluateCondition(execution, s); err != nil {
			e.l.InfoContext(execution, fmt.Sprintf("Skipping step: %s", s.ID))
			continue
		}

		if err := e.executeStepType(execution, s, &nextStep); err != nil {
			return fmt.Errorf("error executing step: %s", s.ID)
		}

		if err := e.handleRetry(execution, s); err != nil {
			return fmt.Errorf("error retrying step: %s", s.ID)
		}
	}

	return nil
}

func (e *Executor) evaluateCondition(execution *Execution, step Step) error {
	if step.Condition == "" {
		return nil
	}

	result, err := Eval(step.Condition, execution.Values)
	if err != nil {
		return fmt.Errorf("error evaluating condition: %v", err)
	}

	if !result.(bool) {
		return fmt.Errorf("condition not met: %s", step.Condition)
	}
	e.l.InfoContext(execution, fmt.Sprintf("Condition met: %s", step.Condition))
	return nil
}

func (e *Executor) executeStepType(execution *Execution, step Step, nextStep *string) error {
	switch step.Type {
	case "assign":
		return e.handleAssign(execution, step)
	case "switch":
		return e.handleSwitch(execution, step, nextStep)
	default:
		return e.handleTask(execution, step)
	}
}

func (e *Executor) handleAssign(execution *Execution, step Step) error {
	for k, v := range step.Args {
		result, err := Eval(v.(string), execution.Values)
		if err != nil {
			return fmt.Errorf("error evaluating expression: %v", err)
		}
		execution.AddValue(k, result)
	}
	return nil
}

func (e *Executor) handleSwitch(execution *Execution, step Step, nextStep *string) error {
	for n, c := range step.Args {
		condition := c.(string)
		result, err := Eval(condition, execution.Values)
		if err != nil {
			return fmt.Errorf("error evaluating expression: %v", err)
		}
		resultBool, ok := result.(bool)
		if !ok {
			return fmt.Errorf("condition %s is not a boolean", condition)
		}
		if resultBool {
			e.l.InfoContext(execution, fmt.Sprintf("Resolving switch: %s is true", condition))
			*nextStep = n
			return nil
		}
		e.l.InfoContext(execution, fmt.Sprintf("Resolving switch: %s is false", condition))
	}
	return nil
}

func (e *Executor) handleTask(execution *Execution, step Step) error {
	task, ok := execution.Container.Tasks[step.Type]
	if !ok {
		return fmt.Errorf("task type: %s not found", step.Type)
	}
	e.executeTask(execution, task, step)
	e.l.InfoContext(execution, fmt.Sprintf("Executed task: %s", step.Type))
	return nil
}

func (e *Executor) handleRetry(execution *Execution, step Step) error {
	if step.Retry == nil {
		return nil
	}

	task, ok := execution.Container.Tasks[step.Type]
	if !ok {
		return fmt.Errorf("Task type: %step not found", step.Type)
	}

	for i := 0; i < step.Retry.MaxRetries; i++ {
		condition, err := Eval(step.Retry.Condition, execution.Values)
		e.l.InfoContext(execution, fmt.Sprintf("[%s/%s] Retrying step: %step, condition: %v\n", strconv.Itoa(i+1), strconv.Itoa(step.Retry.MaxRetries), step.ID, condition))

		if err != nil {
			return fmt.Errorf("error evaluating retry condition: %v", err)
		}

		if !condition.(bool) {
			break
		}

		delay := time.Duration(step.Retry.Delay) * time.Millisecond
		if step.Retry.Backoff {
			delay = time.Duration(i*step.Retry.Delay) * time.Millisecond
		}
		time.Sleep(delay)
		e.executeTask(execution, task, step)
	}
	return nil
}

func (e *Executor) executeTask(execution *Execution, task Task, s Step) {
	output, err := task.Execute(execution, s.Args)

	if err != nil {
		execution.AddValue(fmt.Sprintf("%s.error", s.ID), err.Error())
	}

	for k, v := range output {
		execution.AddValue(fmt.Sprintf("%s.result.%s", s.ID, k), v)
	}
}
