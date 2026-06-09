package internal

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func indexFieldsByTags(command reflect.Type) (map[int]int, map[string]int, error) {
	// Maps to track which struct fields correspond to which positional arguments and flags
	// positionalFields: position (1st, 2nd, etc.) -> struct field index
	// flagFields: flag name (e.g., "--weight") -> struct field index
	positionalFields := map[int]int{}
	flagFields := map[string]int{}

	for i := range command.NumField() {
		fieldType := command.Field(i)

		// If the field has a `pos` tag (e.g. `pos:"1"`), convert it to int and
		// store which struct field index corresponds to which position
		if posTag := fieldType.Tag.Get("pos"); posTag != "" {
			position, err := strconv.Atoi(posTag)
			if err != nil {
				return nil, nil, fmt.Errorf("invalid pos tag for %s: %w", fieldType.Name, err)
			}

			positionalFields[position] = i
		}

		// if the field has a `cli` tag (e.g., `cli:"weight"`),
		// store which struct field corresponds to this flag name
		if cliTag := fieldType.Tag.Get("cli"); cliTag != "" {
			flagFields[cliTag] = i
		}
	}

	return positionalFields, flagFields, nil
}

func checkMissedFlagValue(index int, args []string, flagName string) error {
	if index+1 >= len(args) {
		return fmt.Errorf("missing value for flag: %s", flagName)
	}

	return nil
}

type flagFieldParser func(index int, args []string, arg string, field reflect.Value) (int, error)

var flagParsers = map[reflect.Kind]flagFieldParser{
	reflect.Bool:   parseBoolFlag,
	reflect.Int:    parseIntFlag,
	reflect.String: parseStringFlag,
}

func parseBoolFlag(index int, _ []string, _ string, field reflect.Value) (int, error) {
	field.SetBool(true)
	return index, nil
}

func parseIntFlag(index int, args []string, arg string, field reflect.Value) (int, error) {
	if err := checkMissedFlagValue(index, args, arg); err != nil {
		return index, err
	}

	value, err := strconv.Atoi(args[index+1])
	if err != nil {
		return index, fmt.Errorf("invalid value for flag %s: %w", arg, err)
	}

	field.SetInt(int64(value))

	return index + 1, nil
}

func parseStringFlag(index int, args []string, arg string, field reflect.Value) (int, error) {
	if err := checkMissedFlagValue(index, args, arg); err != nil {
		return index, err
	}

	field.SetString(args[index+1])

	return index + 1, nil
}

func handleFlagArgument(index int, arg string, args []string, cmd reflect.Value, flagFields map[string]int) (int, error) {
	flagName := strings.TrimPrefix(arg, "--")

	fieldIndex, ok := flagFields[flagName]
	if !ok {
		return index, fmt.Errorf("unknown flag: %s", arg)
	}

	// Get the struct field corresponding to this flag to determine its type and set its value later
	field := cmd.Field(fieldIndex)

	parser, ok := flagParsers[field.Kind()]
	if !ok {
		return index, fmt.Errorf("unsupported flag field type for %s", flagName)
	}

	nextIndex, err := parser(index, args, arg, field)
	if err != nil {
		return index, err
	}

	return nextIndex, nil
}

func handleCLICommands(args []string, cmd reflect.Value, positionalFields map[int]int, flagFields map[string]int) error {
	// Track the current position of the positional argument we're trying to parse
	position := 1

	for index := 0; index < len(args); index++ {
		arg := args[index]

		// If the argument starts with "--", it's a flag.
		if strings.HasPrefix(arg, "--") {
			nextIndex, err := handleFlagArgument(index, arg, args, cmd, flagFields)
			if err != nil {
				return err
			}

			index = nextIndex

			continue
		}

		fieldIndex, ok := positionalFields[position]
		if !ok {
			return fmt.Errorf("unexpected positional argument: %s", arg)
		}

		field := cmd.Field(fieldIndex)
		switch field.Kind() {
		case reflect.String:
			field.SetString(arg)
		default:
			return fmt.Errorf("only strings allowed for positional arguments %d", position)
		}

		position++
	}

	return nil
}

func ParseArgs(command reflect.Type, args []string) (Command, error) {
	command, err := NormalizeCommandType(command)
	if err != nil {
		return nil, err
	}

	// Instantiate a brand-new instance of the incoming command to avoid mutating the original struct
	cmdPtr := reflect.New(command)
	cmd := cmdPtr.Elem()

	positionalFields, flagFields, err := indexFieldsByTags(command)
	if err != nil {
		return nil, err
	}

	if err := handleCLICommands(args, cmd, positionalFields, flagFields); err != nil {
		return nil, err
	}

	parsedCommand, ok := cmdPtr.Interface().(Command)
	if !ok {
		return nil, fmt.Errorf("parsed command does not implement Command")
	}

	return parsedCommand, nil
}
