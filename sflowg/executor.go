package sflowg

import (
	"fmt"
)

func Execute(f Flow, e *Execution) error {

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
			args := s.Args.(map[string]any)

			for k, v := range args {
				result, err := Eval(v.(string), e.Values)
				if err != nil {
					fmt.Printf("Error evaluating expression: %v\n", err)
					return err
				}
				e.AddVal(k, result)
			}
		} else if s.Type == "switch" {
			args := s.Args.([]any)

			for _, c := range args {
				cond := c.(map[string]any)

				condition := cond["condition"].(string)
				next := cond["next"].(string)

				result, err := Eval(condition, e.Values)
				if err != nil {
					fmt.Printf("Error evaluating expression: %v\n", err)
					return err
				}

				if result.(bool) {
					fmt.Printf("Resolving switch: %s is true, next step is %s\n", condition, next)
					nextStep = next
					break
				} else {
					fmt.Printf("Resolving switch: %s is false\n", condition)
				}
			}
		} else {
			fmt.Printf("Task type: %s not supported yet\n", s.Type)
		}
	}

	return nil
}
