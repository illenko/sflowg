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
