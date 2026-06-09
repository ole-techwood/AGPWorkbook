package internal

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

type Command interface {
	Execute() error
}

var commandRegistry = map[string]reflect.Type{
	"load":   reflect.TypeFor[*LoadCommand](),
	"cancel": reflect.TypeFor[*CancelCommand](),
}

func GetCommandType(name string) (reflect.Type, bool) {
	commandType, ok := commandRegistry[name]
	return commandType, ok
}

func CommandsUsage() string {
	commandNames := make([]string, 0, len(commandRegistry))
	for commandName := range commandRegistry {
		commandNames = append(commandNames, commandName)
	}

	sort.Strings(commandNames)

	return strings.Join(commandNames, "|")
}

type LoadCommand struct {
	// Positional argument (the first word after the command name)
	Name string `pos:"1" desc:"Name of the cargo"`

	// Flag --weight (named argument)
	Weight int `cli:"weight" desc:"Weight of the cargo in kg"`

	// Flag --hazardous (if present - true, if not — false)
	Hazardous bool `cli:"hazardous" desc:"Is the cargo dangerous?"`
}

type CancelCommand struct {
	// Positional argument (the first word after the command name)
	ID string `pos:"1" desc:"ID of the order to cancel"`
}

func (c *LoadCommand) Execute() error {
	fmt.Printf("🚀 Cargo successfully loaded: %s (Weight: %d, Hazardous: %t)\n", c.Name, c.Weight, c.Hazardous)

	return nil
}

func (c *CancelCommand) Execute() error {
	fmt.Printf("❌ Order successfully canceled: %s\n", c.ID)
	return nil
}
