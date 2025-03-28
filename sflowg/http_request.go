package sflowg

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

type HttpRequestTask struct{}

func (t *HttpRequestTask) Execute(e *Execution, args map[string]any) (map[string]any, error) {
	requestConfig, err := parseArgs(e, args)

	if err != nil {
		return nil, err
	}

	client := e.Container.HttpClient

	response, err := executeRequest(client, requestConfig)

	if err != nil {
		return nil, err
	}

	return response, nil
}

// httpRequestConfig is a struct that holds the configuration for an HTTP request
type httpRequestConfig struct {
	uri         string
	method      string
	headers     map[string]string
	queryParams map[string]string
	body        map[string]any
}

func parseArgs(e *Execution, args map[string]any) (httpRequestConfig, error) {
	uri, ok := args["url"].(string)
	if !ok {
		return httpRequestConfig{}, fmt.Errorf("uri not found or not a string")
	}

	method, ok := args["method"].(string)
	if !ok {
		return httpRequestConfig{}, fmt.Errorf("method not found or not a string")
	}

	headers := make(map[string]any)

	for key, value := range args["headers"].(map[string]any) {
		headerValue, err := Eval(value.(string), e.Values)

		if err != nil {
			return httpRequestConfig{}, err
		}

		headers[key] = headerValue
	}

	queryParameters := make(map[string]any)

	for key, value := range args["queryParameters"].(map[string]any) {
		queryValue, err := Eval(value.(string), e.Values)

		if err != nil {
			return httpRequestConfig{}, err
		}

		queryParameters[key] = queryValue
	}

	body := make(map[string]any)

	for key, value := range args["body"].(map[string]any) {
		bodyValue, err := Eval(value.(string), e.Values)

		if err != nil {
			return httpRequestConfig{}, err
		}

		body[key] = bodyValue
	}

	return httpRequestConfig{
		uri:         uri,
		method:      method,
		headers:     ToStringValueMap(headers),
		queryParams: ToStringValueMap(queryParameters),
		body:        body,
	}, nil
}

// executeRequest executes the HTTP request
func executeRequest(client *resty.Client, config httpRequestConfig) (map[string]any, error) {
	response := map[string]any{}
	errorResponse := map[string]any{}

	resp, err := client.R().
		SetHeaders(config.headers).
		SetQueryParams(config.queryParams).
		SetBody(config.body).
		SetResult(&response).
		SetError(&errorResponse).
		Execute(config.method, config.uri)

	if err != nil {
		return nil, err
	}

	result := make(map[string]any)

	result["status"] = resp.Status()
	result["statusCode"] = resp.StatusCode()
	result["isError"] = resp.IsError()

	if resp.IsError() {
		for k, v := range errorResponse {
			result[fmt.Sprintf("body.%s", k)] = v
		}
	} else {
		for k, v := range response {
			result[fmt.Sprintf("body.%s", k)] = v
		}
	}
	return result, nil
}
