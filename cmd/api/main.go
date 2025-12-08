package main

import (
	"fmt"

	"github.com/kralle333/keyvaluestore/internal/app"
)

func main() {

	config := app.GetTestingConfig()
	app, err := app.NewApp(config)
	if err != nil {
		panic(fmt.Sprintf("Failed to start app: %v", err))
	}
	err = app.Run()
	if err != nil {
		panic(fmt.Sprintf("App stopped with error: %v", err))
	}
}
