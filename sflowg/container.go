package sflowg

import (
	"github.com/go-resty/resty/v2"
)

type Container struct {
	HttpClient *resty.Client
	Tasks      map[string]Task
}

func NewContainer(httpClient *resty.Client, tasks map[string]Task) *Container {
	return &Container{
		HttpClient: httpClient,
		Tasks:      tasks,
	}
}

func (c *Container) GetTask(name string) Task {
	task, ok := c.Tasks[name]
	if !ok {
		return nil
	}
	return task
}

func (c *Container) SetTask(name string, task Task) {
	c.Tasks[name] = task
}
