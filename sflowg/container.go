package sflowg

import (
	"github.com/go-resty/resty/v2"
)

type Container struct {
	HttpClient *resty.Client
}

func NewContainer(httpClient *resty.Client) *Container {
	return &Container{
		HttpClient: httpClient,
	}
}
