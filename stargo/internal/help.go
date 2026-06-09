package internal

import (
	"fmt"
	"reflect"
	"slices"
	"sort"
	"strconv"
	"strings"
)

type argument struct {
	position    int
	placeholder string
	description string
}

type flag struct {
	name        string
	valueHint   string
	description string
}

func ContainsHelpFlag(args []string) bool {
	return slices.Contains(args, "--help")
}

func BuildHelpCommand(commandName string, commandType reflect.Type) (string, error) {
	normalizedType, err := NormalizeCommandType(commandType)
	if err != nil {
		return "", err
	}

	arguments, flags, err := extractHelpEntries(normalizedType)
	if err != nil {
		return "", err
	}

	builder := strings.Builder{}

	fmt.Fprintf(&builder, "Command: %s\n", commandName)
	fmt.Fprintf(&builder, "Usage:   %s\n", buildUsage(commandName, arguments, flags))

	writeArgumentsSection(&builder, arguments)
	writeFlagsSection(&builder, flags)

	return builder.String(), nil
}

func extractHelpEntries(commandType reflect.Type) ([]argument, []flag, error) {
	arguments := make([]argument, 0)
	flags := make([]flag, 0)

	for field := range commandType.Fields() {
		argument, hasArgument, err := argumentFromField(field)
		if err != nil {
			return nil, nil, err
		}
		if hasArgument {
			arguments = append(arguments, argument)
		}

		flag, hasFlag := flagFromField(field)
		if hasFlag {
			flags = append(flags, flag)
		}
	}

	sort.Slice(arguments, func(i, j int) bool {
		return arguments[i].position < arguments[j].position
	})

	return arguments, flags, nil
}

func argumentFromField(field reflect.StructField) (argument, bool, error) {
	posTag := field.Tag.Get("pos")
	if posTag == "" {
		return argument{}, false, nil
	}

	position, err := strconv.Atoi(posTag)
	if err != nil {
		return argument{}, false, fmt.Errorf("invalid pos tag for %s: %w", field.Name, err)
	}

	return argument{
		position:    position,
		placeholder: fmt.Sprintf("<%s>", strings.ToLower(field.Name)),
		description: getDescription(field.Tag.Get("desc")),
	}, true, nil
}

func flagFromField(field reflect.StructField) (flag, bool) {
	cliTag := field.Tag.Get("cli")
	if cliTag == "" {
		return flag{}, false
	}

	return flag{
		name:        fmt.Sprintf("--%s", cliTag),
		valueHint:   getFlagValueHint(field.Type.Kind()),
		description: getDescription(field.Tag.Get("desc")),
	}, true
}

func buildUsage(commandName string, arguments []argument, flags []flag) string {
	usageParts := make([]string, 0, len(arguments)+len(flags)+1)
	usageParts = append(usageParts, fmt.Sprintf("go run cmd/main.go %s", commandName))

	for _, argument := range arguments {
		usageParts = append(usageParts, argument.placeholder)
	}

	for _, flag := range flags {
		usageParts = append(usageParts, fmt.Sprintf("[%s%s]", flag.name, flag.valueHint))
	}

	return strings.Join(usageParts, " ")
}

func writeArgumentsSection(builder *strings.Builder, arguments []argument) {
	if len(arguments) == 0 {
		return
	}

	builder.WriteString("\nArguments:\n")
	maxWidth := 0
	for _, argument := range arguments {
		if len(argument.placeholder) > maxWidth {
			maxWidth = len(argument.placeholder)
		}
	}

	for _, argument := range arguments {
		fmt.Fprintf(builder, "  %-*s  %s\n", maxWidth, argument.placeholder, argument.description)
	}
}

func writeFlagsSection(builder *strings.Builder, flags []flag) {
	if len(flags) == 0 {
		return
	}

	builder.WriteString("\nFlags:\n")
	maxWidth := 0
	flagLabels := make([]string, 0, len(flags))

	for _, flag := range flags {
		label := fmt.Sprintf("%s%s", flag.name, flag.valueHint)
		flagLabels = append(flagLabels, label)
		if len(label) > maxWidth {
			maxWidth = len(label)
		}
	}

	for i, flag := range flags {
		fmt.Fprintf(builder, "  %-*s  %s\n", maxWidth, flagLabels[i], flag.description)
	}
}

func getFlagValueHint(kind reflect.Kind) string {
	if kind == reflect.Bool {
		return ""
	}

	return fmt.Sprintf(" <%s>", kind)
}

func getDescription(description string) string {
	if description == "" {
		return "No description provided."
	}

	return description
}
