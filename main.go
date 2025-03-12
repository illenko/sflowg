package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Flow struct {
	Entrypoint Entrypoint     `yaml:"entrypoint"`
	Steps      []Step         `yaml:"steps"`
	Properties map[string]any `yaml:"properties"`
	Return     Return         `yaml:"return"`
}

type Entrypoint struct {
	Type   string `yaml:"type"`
	Path   string `yaml:"path"`
	Result string `yaml:"result"`
}

type Step struct {
	ID     string `yaml:"id"`
	Type   string `yaml:"type"`
	Args   any    `yaml:"args"`
	Next   string `yaml:"next,omitempty"`
	Result string `yaml:"result,omitempty"`
}

type Return struct {
	Type string                 `yaml:"type"`
	Args map[string]interface{} `yaml:"args"`
}

type Execution struct {
	Flow    Flow
	Context map[string]any
}

func main() {
	yamlFile, err := os.ReadFile("flows/test_flow.yaml")
	if err != nil {
		fmt.Printf("Error reading YAML file: %v\n", err)
		return
	}

	var flow Flow
	err = yaml.Unmarshal(yamlFile, &flow)
	if err != nil {
		fmt.Printf("Error unmarshalling YAML: %v\n", err)
		return
	}

	fmt.Printf("Flow: %+v\n", flow)
}

func findNextStepIndex(steps []Step, currentStepID string) int {
	for i, step := range steps {
		if step.ID == currentStepID && i+1 < len(steps) {
			return i + 1
		}
	}
	return -1
}

func findStepByID(steps []Step, id string) (Step, bool) {
	for _, step := range steps {
		if step.ID == id {
			return step, true
		}
	}
	return Step{}, false
}
