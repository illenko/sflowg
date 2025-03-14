package task

import "github.com/go-resty/resty/v2"

const (
	HeaderPrefix = "header."
	QueryPrefix  = "query."
	BodyPrefix   = "body."
)

// httpRequestConfig is a struct that holds the configuration for an HTTP request
type httpRequestConfig struct {
	uri         string
	method      string
	headers     map[string]string
	queryParams map[string]string
	body        map[string]any
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

	output := make(map[string]any)

	output["http-status"] = resp.Status()
	output["http-status-code"] = resp.StatusCode()
	output["success"] = !resp.IsError()

	if resp.IsError() {
		for k, v := range errorResponse {
			output[k] = v
		}
	} else {
		for k, v := range response {
			output[k] = v
		}
	}
	return output, nil
}
