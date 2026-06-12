package main

import (
	"fmt"
	"os"

	"github.com/ole-techwood/AGPWorkbook/stargo/internal"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("usage: go run cmd/main.go <%s> ...\n", internal.CommandsUsage())
		os.Exit(1)
	}

	commandName := os.Args[1]
	args := os.Args[2:]

	commandType, ok := internal.GetCommandType(commandName)
	if !ok {
		fmt.Printf("unknown command: %s (supported: %s)\n", commandName, internal.CommandsUsage())
		os.Exit(1)
	}

	if internal.ContainsHelpFlag(args) {
		helpText, err := internal.BuildHelpCommand(commandName, commandType)
		if err != nil {
			fmt.Printf("failed to build help for %s command: %v\n", commandName, err)
			os.Exit(1)
		}

		fmt.Print(helpText)
		return
	}

	command, err := internal.ParseArgs(commandType, args)
	if err != nil {
		fmt.Printf("failed to parse %s command: %v\n", commandName, err)
		os.Exit(1)
	}

	container := internal.NewContainer()
	container.Register(internal.NewLogger())

	if err := container.Inject(command); err != nil {
		fmt.Printf("failed to inject dependencies for %s command: %v\n", commandName, err)
		os.Exit(1)
	}

	if err := command.Execute(); err != nil {
		fmt.Printf("failed to execute %s command: %v\n", commandName, err)
		os.Exit(1)
	}
}
