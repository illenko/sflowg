package sflowg

import (
	"fmt"
	"sflowg/sflowg/container"
	"sflowg/sflowg/task"
)

func Execute(c container.Container, f Flow, e *Execution) error {

	nextStep := ""

	for _, s := range f.Steps {
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
		} else if s.Type == "http" {
			httpTask := task.HttpRequest{}

			output, err := httpTask.Execute(&c, e, s.Args)

			if err != nil {
				return err
			}

			for k, v := range output {
				e.AddVal(k, v)
			}
		} else {
			fmt.Printf("Task type: %s not supported yet\n", s.Type)
		}
	}

	return nil
}
