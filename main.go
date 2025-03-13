package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"sflowg/sflowg"
)

func main() {
	yamlFile, err := os.ReadFile("flows/test_flow.yaml")
	if err != nil {
		fmt.Printf("Error reading YAML file: %v\n", err)
		return
	}

	var flow sflowg.Flow
	err = yaml.Unmarshal(yamlFile, &flow)
	if err != nil {
		fmt.Printf("Error unmarshalling YAML: %v\n", err)
		return
	}

	fmt.Printf("Flow: %+v\n", flow)

	execution := sflowg.NewExecution(flow)

	for k, v := range flow.Properties {
		execution.AddContext(fmt.Sprintf("properties.%s", k), v)
	}

	for _, step := range flow.Steps {
		if step.Type == "assign" {
			args := step.Args.(map[string]any)

			for k, v := range args {
				result, err := sflowg.Eval(v.(string), execution.Context)
				if err != nil {
					fmt.Printf("Error evaluating expression: %v\n", err)
					return
				}
				execution.AddContext(k, result)
			}
		}
	}
}
