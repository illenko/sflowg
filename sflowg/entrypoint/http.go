package entrypoint

import (
	"fmt"
	"io"
	"net/http"
	"sflowg/sflowg"
	"strings"

	"github.com/Jeffail/gabs/v2"
	"github.com/gin-gonic/gin"
)

func NewHttpHandler(f sflowg.Flow, g *gin.Engine) {
	config := f.Entrypoint.Config
	method := strings.ToLower(config["method"].(string))
	path := config["path"].(string)

	fmt.Printf("registering HTTP entrypoint for %s %s \n", method, path)

	switch method {
	case "get":
		g.GET(path, handleRequest(f, false))
	case "post":
		g.POST(path, handleRequest(f, true))
	default:
		fmt.Printf("Method %s is not supported", method)
	}
}

func handleRequest(f sflowg.Flow, body bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		e := sflowg.NewExecution(f)

		handlePathVariables(c, f, &e)
		handleQueryParameters(c, f, &e)
		handleHeaders(c, f, &e)

		if body {
			handleBody(c, f, &e)
		}

		err := sflowg.Execute(f, &e)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Error in task execution: " + err.Error(),
			})
			return
		}

		toResponse(c, &e)
	}
}

func handlePathVariables(c *gin.Context, f sflowg.Flow, e *sflowg.Execution) {
	result := f.Entrypoint.Result
	pathVariables := f.Entrypoint.Config["pathVariables"].([]any)
	for _, p := range pathVariables {
		v := p.(string)
		e.AddVal(fmt.Sprintf("%s.pathVariables.%s", result, v), c.Param(v))
	}
}

func handleQueryParameters(c *gin.Context, f sflowg.Flow, e *sflowg.Execution) {
	result := f.Entrypoint.Result
	queryParameters := f.Entrypoint.Config["queryParameters"].([]any)
	for _, q := range queryParameters {
		v := q.(string)
		e.AddVal(fmt.Sprintf("%s.queryParameters.%s", result, v), c.Query(v))
	}
}

func handleHeaders(c *gin.Context, f sflowg.Flow, e *sflowg.Execution) {
	result := f.Entrypoint.Result
	headers := f.Entrypoint.Config["headers"].([]any)
	for _, h := range headers {
		v := h.(string)
		e.AddVal(fmt.Sprintf("%s.headers.%s", result, v), c.GetHeader(v))
	}
}

func handleBody(c *gin.Context, f sflowg.Flow, e *sflowg.Execution) {
	bodyConfig := f.Entrypoint.Config["body"].(map[string]any)
	bodyType := bodyConfig["type"].(string)

	if bodyType == "json" {
		handleJSONBody(c, f, e)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Body type is not supported"})
	}
}

var wrongBodyFormatRes = gin.H{"message": "Wrong request body format"}

func handleJSONBody(c *gin.Context, f sflowg.Flow, e *sflowg.Execution) {
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

	result := f.Entrypoint.Result
	for k, v := range values {
		e.AddVal(fmt.Sprintf("%s.body.%s", result, k), v)
	}
}

func toResponse(c *gin.Context, e *sflowg.Execution) {
	jsonObj := gabs.New()

	res := make(map[string]any)
	res["test"] = "testVal"

	for k, v := range res {
		_, _ = jsonObj.SetP(v, k)
	}

	c.JSON(http.StatusOK, e)
}
