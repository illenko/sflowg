package sflowg

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Jeffail/gabs/v2"
	"github.com/gin-gonic/gin"
)

func NewHttpHandler(flow *Flow, container *Container, g *gin.Engine) {
	config := flow.Entrypoint.Config
	method := strings.ToLower(config["method"].(string))
	path := config["path"].(string)

	fmt.Printf("registering HTTP entrypoint for %s %s \n", method, path)

	switch method {
	case "get":
		g.GET(path, handleRequest(flow, container, false))
	case "post":
		g.POST(path, handleRequest(flow, container, true))
	default:
		fmt.Printf("Method %s is not supported", method)
	}
}

func handleRequest(flow *Flow, container *Container, withBody bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		e := NewExecution(flow, container)

		extractRequestData(c, flow, &e, withBody)

		err := ExecuteSteps(&e)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Error in task execution: " + err.Error(),
			})
			return
		}

		toResponse(c, &e)
	}
}

const (
	PathVariablesKey   = "pathVariables"
	QueryParametersKey = "queryParameters"
	HeadersKey         = "headers"

	PathVariablesPrefix   = "request.pathVariables"
	QueryParametersPrefix = "request.queryParameters"
	HeadersPrefix         = "request.headers"
	RequestBodyPrefix     = "request.body"
)

func extractRequestData(c *gin.Context, f *Flow, e *Execution, withBody bool) {
	if pathVariables, ok := f.Entrypoint.Config[PathVariablesKey].([]any); ok {
		extractValues(e, pathVariables, PathVariablesPrefix, c.Param)
	}

	if queryParameters, ok := f.Entrypoint.Config[QueryParametersKey].([]any); ok {
		extractValues(e, queryParameters, QueryParametersPrefix, c.Query)
	}

	if headers, ok := f.Entrypoint.Config[HeadersKey].([]any); ok {
		extractValues(e, headers, HeadersPrefix, c.GetHeader)
	}

	if withBody {
		extractBody(c, f, e)
	}
}

func extractValues(e *Execution, keys []any, prefix string, getValue func(string) string) {
	for _, key := range keys {
		if v, ok := key.(string); ok {
			e.AddVal(fmt.Sprintf("%s.%s", prefix, v), getValue(v))
		}
	}
}

func extractBody(c *gin.Context, f *Flow, e *Execution) {
	bodyConfig := f.Entrypoint.Config["body"].(map[string]any)
	bodyType := bodyConfig["type"].(string)

	if bodyType == "json" {
		extractJsonBody(c, e)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Body type is not supported"})
	}
}

var wrongBodyFormatRes = gin.H{"message": "Wrong request body format"}

func extractJsonBody(c *gin.Context, e *Execution) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, wrongBodyFormatRes)
		return
	}

	bodyParsed, err := gabs.ParseJSON(body)
	if err != nil {
		c.JSON(http.StatusBadRequest, wrongBodyFormatRes)
		return
	}

	values, err := bodyParsed.Flatten()

	if err != nil {
		c.JSON(http.StatusBadRequest, wrongBodyFormatRes)
		return
	}

	for k, v := range values {
		e.AddVal(fmt.Sprintf("%s.%s", RequestBodyPrefix, k), v)
	}
}

func toResponse(c *gin.Context, e *Execution) {
	jsonObj := gabs.New()

	res := make(map[string]any)
	res["test"] = "testVal"

	for k, v := range res {
		_, _ = jsonObj.SetP(v, k)
	}

	c.JSON(http.StatusOK, e)
}
