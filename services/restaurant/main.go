package main

import (
	"fmt"
	"os"

	"github.com/katatrina/food-delivery/services/restaurant/cmd/app"
)

func main() {
	if err := app.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
