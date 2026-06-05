package main

import (
	"fmt"
	"reflect"
)

type Command interface {
	Execute() error
}

type LoadCommand struct {
	// Positional argument (the first word after the command name)
	Name string `pos:"1" desc:"Name of the cargo"`

	// Flag --weight (named argument)
	Weight int `cli:"weight" desc:"Weight of the cargo in kg"`

	// Flag --hazardous (if present - true, if not — false)
	Hazardous bool `cli:"hazardous" desc:"Is the cargo dangerous?"`
}

func (c *LoadCommand) Execute() error {
	fmt.Printf("🚀 Cargo successfully loaded: %s (Weight: %d, Hazardous: %t)\n", c.Name, c.Weight, c.Hazardous)

	return nil
}

func main() {
	// TODO: implement the main logic
}

func parseArgs(command reflect.Type, args []string) (Command, error) {
	// TODO: implement argument parsing based on struct tags

	return nil, nil
}
