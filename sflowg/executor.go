package sflowg

import (
	"fmt"
	"strconv"
	"time"
)

func ExecuteSteps(e *Execution) error {
	nextStep := ""

	for _, s := range e.Flow.Steps {
		if nextStep != "" {
			if s.ID != nextStep {
				fmt.Printf("Skipping step: %s, searching for %s\n", s.ID, nextStep)
				continue
			}
			nextStep = ""
			fmt.Printf("Found expected step: %s\n", s.ID)
		}

		if err := evaluateCondition(e, s); err != nil {
			fmt.Printf("Skipping step: %s\n", s.ID)
			continue
		}

		if err := executeStepType(e, s, &nextStep); err != nil {
			return err
		}

		if err := handleRetry(e, s); err != nil {
			return err
		}
	}

	return nil
}

func evaluateCondition(e *Execution, s Step) error {
	if s.Condition == "" {
		return nil
	}

	result, err := Eval(s.Condition, e.Values)
	if err != nil {
		fmt.Printf("Error evaluating expression: %v\n", err)
		return err
	}

	if !result.(bool) {
		fmt.Printf("Skipping step: %s\n", s.ID)
		return fmt.Errorf("condition not met") // return an error so that the calling function knows the step was skipped, but continue.
	}
	fmt.Printf("Condition %s is satisfied for step: %s\n", s.Condition, s.ID)
	return nil
}

func executeStepType(e *Execution, s Step, nextStep *string) error {
	switch s.Type {
	case "assign":
		return handleAssign(e, s)
	case "switch":
		return handleSwitch(e, s, nextStep)
	default:
		return handleTask(e, s)
	}
}

func handleAssign(e *Execution, s Step) error {
	for k, v := range s.Args {
		result, err := Eval(v.(string), e.Values)
		if err != nil {
			fmt.Printf("Error evaluating expression: %v\n", err)
			return err
		}
		e.AddVal(k, result)
	}
	return nil
}

func handleSwitch(e *Execution, s Step, nextStep *string) error {
	for n, c := range s.Args {
		condition := c.(string)
		result, err := Eval(condition, e.Values)
		if err != nil {
			fmt.Printf("Error evaluating expression: %v\n", err)
			return err
		}
		resultBool, ok := result.(bool)
		if !ok {
			return fmt.Errorf("condition %s is not a boolean", condition)
		}
		if resultBool {
			fmt.Printf("Resolving switch: %s is true, next step is %s\n", condition, n)
			*nextStep = n
			return nil
		}
		fmt.Printf("Resolving switch: %s is false\n", condition)
	}
	return nil
}

func handleTask(e *Execution, s Step) error {
	task, ok := e.Container.Tasks[s.Type]
	if !ok {
		fmt.Printf("Task type: %s not found\n", s.Type)
		return fmt.Errorf("task type: %s not found", s.Type)
	}
	executeTask(e, task, s)
	fmt.Printf("Step %s executed\n", s.ID)
	return nil
}

func handleRetry(e *Execution, s Step) error {
	if s.Retry == nil {
		return nil
	}

	task, ok := e.Container.Tasks[s.Type]
	if !ok {
		return fmt.Errorf("Task type: %s not found", s.Type)
	}

	for i := 0; i < s.Retry.MaxRetries; i++ {
		condition, err := Eval(s.Retry.Condition, e.Values)
		fmt.Printf("[%s/%s] Retrying step: %s, condition: %v\n", strconv.Itoa(i+1), strconv.Itoa(s.Retry.MaxRetries), s.ID, condition)

		if err != nil {
			fmt.Printf("Error evaluating expression: %v\n", err)
			return err
		}

		if !condition.(bool) {
			break
		}

		delay := time.Duration(s.Retry.Delay) * time.Millisecond
		if s.Retry.Backoff {
			delay = time.Duration(i*s.Retry.Delay) * time.Millisecond
		}
		time.Sleep(delay)
		executeTask(e, task, s)
	}
	return nil
}

func executeTask(e *Execution, task Task, s Step) {
	output, err := task.Execute(e, s.Args)

	if err != nil {
		e.AddVal(fmt.Sprintf("%s.error", s.ID), err.Error())
	}

	for k, v := range output {
		e.AddVal(fmt.Sprintf("%s.result.%s", s.ID, k), v)
	}
}
