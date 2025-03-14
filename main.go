package main

import (
	"fmt"
	"os"
	"sflowg/sflowg"
	"sflowg/sflowg/entrypoint"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
)

func main() {
	yamlFile, err := os.ReadFile("flows/test_flow.yaml")
	if err != nil {
		fmt.Printf("Error reading YAML file: %v\n", err)
		return
	}

	var flow sflowg.Flow
	err = yaml.Unmarshal(yamlFile, &flow)
	if err != nil {
		fmt.Printf("Error unmarshalling YAML: %v\n", err)
		return
	}

	fmt.Printf("Flow: %+v\n", flow)

	g := gin.Default()

	entrypoint.NewHttpHandler(flow, g)

	err = g.Run(":8080")

	if err != nil {
		fmt.Printf("Error running server: %v", err)
	}

	//execution := sflowg.NewExecution(flow)
	//
	//for k, v := range flow.Properties {
	//	execution.AddVal(fmt.Sprintf("properties.%s", k), v)
	//}
	//

}
