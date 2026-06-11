package internal

import (
	"fmt"
	"os"
	"reflect"
	"sort"
	"strconv"
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
	existing, err := LoadCargoList(cargoFilePath)
	if err != nil {
		return fmt.Errorf("load: read cargo state: %w", err)
	}

	totalWeight := c.Weight
	for _, item := range existing {
		totalWeight += item.Weight
	}

	if totalWeight > weightLimit {
		return fmt.Errorf("load: total weight would be %d kg, exceeds %d kg limit", totalWeight, weightLimit)
	}

	newCargo := Cargo{
		ID:          strconv.Itoa(len(existing)),
		LoadCommand: *c,
	}
	existing = append(existing, newCargo)

	if err := os.MkdirAll(cargoDataDir, 0o755); err != nil {
		return fmt.Errorf("load: create data dir: %w", err)
	}

	if err := SaveCargoList(cargoFilePath, existing); err != nil {
		return fmt.Errorf("load: save cargo state: %w", err)
	}

	fmt.Printf("Cargo loaded: %s (ID: %s, total weight: %d kg)\n", c.Name, newCargo.ID, totalWeight)
	return nil
}

func (c *CancelCommand) Execute() error {
	fmt.Printf("❌ Order successfully canceled: %s\n", c.ID)
	return nil
}
