// Package main provides the mbaigo CLI tool for managing Arrowhead systems
package main

import (
	"fmt"
	"os"

	"github.com/eislab-cps/mbaigo/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
