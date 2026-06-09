package internal

import (
	"fmt"
	"reflect"
)

func NormalizeCommandType(command reflect.Type) (reflect.Type, error) {
	// If consumer of the function passed a pointer to a struct (*LoadCommand),
	// we get the struct's base type (LoadCommand) to read its fields.
	if command.Kind() == reflect.Pointer {
		command = command.Elem()
	}

	// The command must be ONLY a struct
	// Return an error if the consumer of the function passed a non-struct type (e.g., int, string, etc.)
	if command.Kind() != reflect.Struct {
		return nil, fmt.Errorf("command type must be a struct")
	}

	return command, nil
}
