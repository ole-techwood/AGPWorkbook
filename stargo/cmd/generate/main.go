package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ole-techwood/AGPWorkbook/stargo/internal"
)

func main() {
	input := flag.String("input", "../commands.json", "path to commands.json")
	output := flag.String("output", "generated_commands.go", "output file path")
	flag.Parse()

	if err := internal.GenerateCommands(*input, *output); err != nil {
		fmt.Fprintf(os.Stderr, "generate: %v\n", err)
		os.Exit(1)
	}
}
