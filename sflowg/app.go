package sflowg

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-resty/resty/v2"
	"gopkg.in/yaml.v3"
)

type App struct {
	Container *Container
	Flows     map[string]Flow
}

func NewApp(flowsDir string) (*App, error) {
	files, err := filepath.Glob(filepath.Join(flowsDir, "*.yaml"))
	if err != nil {
		return nil, fmt.Errorf("error reading directory: %w", err)
	}

	httpClient := resty.New().SetDebug(true)

	app := App{
		Container: NewContainer(httpClient, builtInTasks),
		Flows:     make(map[string]Flow),
	}

	for _, file := range files {
		flow, err := readFlow(file)
		if err != nil {
			return nil, err
		}
		app.RegisterFlow(flow)
	}

	return &app, nil
}

var builtInTasks = map[string]Task{
	"http": &HttpRequestTask{},
}

func (a *App) RegisterTask(name string, task Task) {
	a.Container.SetTask(name, task)
}

func (a *App) RegisterFlow(flow Flow) {
	a.Flows[flow.ID] = flow
}

func readFlow(file string) (Flow, error) {
	yamlFile, err := os.ReadFile(file)
	if err != nil {
		return Flow{}, fmt.Errorf("error reading YAML file: %w", err)
	}

	var flow Flow
	err = yaml.Unmarshal(yamlFile, &flow)
	if err != nil {
		return Flow{}, fmt.Errorf("error unmarshalling YAML: %w", err)
	}

	return flow, nil
}
