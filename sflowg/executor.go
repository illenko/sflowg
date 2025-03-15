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
			} else {
				nextStep = ""
				fmt.Printf("Found expected step: %s\n", s.ID)
			}
		}

		if s.Condition != "" {
			result, err := Eval(s.Condition, e.Values)
			if err != nil {
				fmt.Printf("Error evaluating expression: %v\n", err)
				return err
			}

			if !result.(bool) {
				fmt.Printf("Skipping step: %s\n", s.ID)
				continue
			} else {
				fmt.Printf("Condition %s is satisfied for step: %s\n", s.Condition, s.ID)
			}

		}

		if s.Type == "assign" {
			for k, v := range s.Args {
				result, err := Eval(v.(string), e.Values)
				if err != nil {
					fmt.Printf("Error evaluating expression: %v\n", err)
					return err
				}
				e.AddVal(k, result)
			}
		} else if s.Type == "switch" {
			for n, c := range s.Args {
				condition := c.(string)

				result, err := Eval(condition, e.Values)
				if err != nil {
					fmt.Printf("Error evaluating expression: %v\n", err)
					return err
				}

				resultBool, ok := result.(bool)

				if !ok {
					fmt.Printf("Error evaluating expression: %v\n", err)
					return fmt.Errorf("condition %s is not a boolean", condition)
				}

				if resultBool {
					fmt.Printf("Resolving switch: %s is true, next step is %s\n", condition, n)
					nextStep = n
					break
				} else {
					fmt.Printf("Resolving switch: %s is false\n", condition)
				}
			}
		} else {
			task, ok := e.Container.Tasks[s.Type]

			if !ok {
				fmt.Printf("Task type: %s not found\n", s.Type)
				return fmt.Errorf("Task type: %s not found", s.Type)
			}

			executeTask(e, task, s)

			fmt.Printf("Step %s executed\n", s.ID)

			if s.Retry != nil {
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

					if s.Retry.Backoff {
						time.Sleep(time.Duration(i*s.Retry.Delay) * time.Millisecond)
					} else {
						time.Sleep(time.Duration(s.Retry.Delay) * time.Millisecond)
					}

					executeTask(e, task, s)

				}
			}
		}
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
