package main

import (
	"context"
	"fmt"
	"os"
)

func main() {
	err := Provider().Run(context.Background(), Name, Version)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s", err.Error())
		os.Exit(1)
	}
}
