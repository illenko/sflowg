package main

import (
	"fmt"
	"log/slog"
	"os"
	"sflowg/sflowg"

	"github.com/gin-gonic/gin"
)

func main() {
	app, err := sflowg.NewApp("flows")

	if err != nil {
		fmt.Printf("Error initializing app: %v", err)
		return
	}

	g := gin.Default()

	flow := app.Flows["test_flow"]

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	executor := sflowg.NewExecutor(logger)

	sflowg.NewHttpHandler(&flow, app.Container, executor, g)

	err = g.Run(":8080")

	if err != nil {
		fmt.Printf("Error running server: %v", err)
	}
}
