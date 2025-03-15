package sflowg

import (
	"fmt"
	"os"

	"github.com/go-resty/resty/v2"
	"gopkg.in/yaml.v3"
)

type App struct {
	Container *Container
	Flows     map[string]Flow
}

func NewApp() App {
	// todo: implement reading multiple flows
	yamlFile, err := os.ReadFile("flows/test_flow.yaml")
	if err != nil {
		fmt.Printf("Error reading YAML file: %v\n", err)
	}

	var flow Flow
	err = yaml.Unmarshal(yamlFile, &flow)
	if err != nil {
		fmt.Printf("Error unmarshalling YAML: %v\n", err)
	}

	httpClient := resty.New().SetDebug(true)

	return App{
		Container: NewContainer(httpClient),
		Flows:     map[string]Flow{flow.ID: flow},
	}
}
