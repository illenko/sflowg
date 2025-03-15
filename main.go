package main

import (
	"fmt"
	"sflowg/sflowg"

	"github.com/gin-gonic/gin"
)

func main() {
	app := sflowg.NewApp()

	g := gin.Default()

	flow := app.Flows["test_flow"]

	sflowg.NewHttpHandler(&flow, app.Container, g)

	err := g.Run(":8080")

	if err != nil {
		fmt.Printf("Error running server: %v", err)
	}
}
