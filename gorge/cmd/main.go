package main

import (
	"fmt"
	"io"
	"os"

	"github.com/ole-techwood/AGPWorkbook/gorge/internal/cli"
)

// commandFunc is contract every command must satisfy to be dispatched.
type commandFunc func(args []string, stdout, stderr io.Writer) int

// commands registers available commands. Adding a new command means adding
// an entry here, not modifying the dispatch logic in run().
var commands = map[string]commandFunc{
	"run": cli.Run,
}

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		printUsage(stderr)
		return 1
	}

	command, ok := commands[args[0]]
	if !ok {
		fmt.Fprintf(stderr, "unknown command: %s\n", args[0])
		printUsage(stderr)
		return 1
	}

	return command(args[1:], stdout, stderr)
}

func printUsage(out io.Writer) {
	fmt.Fprintln(out, "Usage:")
	fmt.Fprintln(out, "  gorge run --stream=\"<comma_separated_raw_events>\"")
}
