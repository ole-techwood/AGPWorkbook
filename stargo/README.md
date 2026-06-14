# Stargo

## Overview

Manage cargo logistics through entire galaxy using Stargo CLI.

## Requirements

### Chapter 1

Let’s start from using reflection to parse the command line and dynamically map data into structures using custom struct tags.

Go’s Flag package does something similar, but for the sake of studying, we’ll implement parsing on our own. Of course, our version is super simple comparing to Flag package

#### CLI Arguments

Each CLI argument must be a struct with custom tags

- **`cli**:"name"`: Establishes the flag name (e.g., `--weight`).
- **`pos**:"position"`: Maps positional text to a struct field based on its index (e.g., the first argument after the command name).
- **`desc**:"description"`: Contains help text for that specific field.

Example:

```go
type LoadCommand struct {
	// Positional argument (the first word after the command name)
	Name string `pos:"1" desc:"Name of the cargo"`

	// Flag --weight (named argument)
	Weight int `cli:"weight" desc:"Weight of the cargo in kg"`

	// Flag --hazardous (if present - true, if not — false)
	Hazardous bool `cli:"hazardous" desc:"Is the cargo dangerous?"`
}
```

#### Arguments parsing

Implement a core parsing function that accepts command type (`command reflect.Type`) and a slice of command arguments (`os.Args[2:]`).

Using the `reflect` package, it must:

1. Dynamically instantiate a brand-new instance of the incoming command:

```go
cmdPtr := reflect.New(cmdType)
cmd := cmdPtr.Elem()
```

1. Inspect the fields and tags of the incoming command
2. Parse the raw command-line strings
3. Bind them into the correct fields (Name, Weight, Hazardous)
4. Return the newly created and populated command.

#### **Dynamic Type Conversion**

The parser must handle conversions from raw strings into at least: `string`, `int`, and `bool` (e.g., treating `-hazardous` or `-hazardous=true` as a boolean flag).

#### **Action Execution**

Each command struct must expose an execution trigger via an interface with method like `Execute()` which is called by the runner after successful binding:

```go
type Command interface {
	Execute() error
}

func (c *LoadCommand) Execute() error {
	fmt.Printf("🚀 Cargo successfully loaded: %s (Weight: %d, Hazardous: %t)\n", c.Name, c.Weight, c.Hazardous)

	return nil
}
```

#### **Auto-Help Generation**

Create a centralized registry of commands. If the user passes `--help` to any command, the engine must use reflection to inspect the target struct's fields and tags to print a beautifully formatted help menu displaying all available flags, arguments, and descriptions.

### Chapter 2

We’ve loaded the cargo and missed it right off the bat…

Let’s make our CLI remember what we’ve loaded into our spaceship.

#### Custom JSON Marshaller

Implement a function `MarshalToJSON(v interface{}) ([]byte, error)` that uses reflection to inspect the value `v`. It must read all public capitalized fields and format them into a valid JSON string (e.g., `{"FieldName": value}`).

Yes, we re-invent the wheel again, but we still do it with good intentions - improve reflection understanding and learn a bit of how JSON marshalling and unmarshalling works under the hood

#### Custom JSON Unmarshaller

Implement a function `UnmarshalFromJSON(data []byte, t reflect.Type) (reflect.Value, error)`. It must accept the `reflect.Type` of the target structure, parse the JSON text, and return a brand-new, dynamically instantiated and populated structure.

#### Cargo State Persistence

- When `load` is executed, the application must check if a `cargo.json` file already exists. If it does, use `UnmarshalFromJSON` to read the state of previously loaded cargo. If it doesn’t, create `cargo.json` and save first loaded cargo.
- The execution logic must validate that the **total weight** of the cargo (existing weight from `cargo.json` + new cargo weight) **does not exceed 4,000,000 kg**. If it exceeds this limit, return a clear validation error and abort the operation.
- If validation passes, use `MarshalToJSON` to save the new state back into `cargo.json`.

### Chapter 3

Our `LoadCommand` is getting heavy, knowing too much about errors printing, which makes it hard to test or swap implementations. Let's inject custom logger into `LoadCommand` using a lightweight Dependency Injection (DI) Container. The goal of this chapter is not to teach all of the DI concepts, its focus on reflection, that’s why we intentionally keep the DI simple.

#### Logger

Define `Logger` infrastructure service to isolate core business logic from outer side-effects like logging into console

#### Runtime DI Container

Implement a `Container` registry struct that maps component types (`reflect.Type`) to their corresponding active concrete singleton instances (`reflect.Value`). It must expose a registration method:

```go
func (c *Container) Register(dependency any)
```

#### Reflection-Driven Injection

Implement an injection engine function `func (c *Container) Inject(target any) error`. Before executing any dynamically instantiated command struct, the DI engine must:

1. Accept a pointer to the newly created command instance.
2. Iterate over its fields using reflection to look up a new custom tag: `di:"inject"`.
3. Query the dependency registry using the field's type.
4. Dynamically assign and populate the matching service singleton to that command field using `reflect.Value.Set()`.

#### Command Decoupling

Update `LoadCommand` so it contains no direct invocation of `fmt` functions. Instead, it must safely request its dependencies using the `di:"inject"` tag and utilize them inside `Execute()`:

```go
type LoadCommand struct {
	Name      string `pos:"1" desc:"Name of the cargo"`
	Weight    int    `cli:"weight" desc:"Weight of the cargo in kg"`
	Hazardous bool   `cli:"hazardous" desc:"Is the cargo dangerous?"`

	// Injected Infrastructure Services
  Logger *Logger `di:"inject"`
}
```

### **Chapter 4**

To explore metaprogramming in Go, we will move away from manually writing repetitive boilerplate for command structs. Instead, you will implement a code generation tool that uses Go templates (`text/template`) to programmatically output functional Go source code based on a JSON declaration file.

#### Command Definition Schema

Define a JSON format that describes a command, its fields, tags, and dependencies. For example, a file named `commands.json` will look like this:

```json
[
  {
    "CommandName": "Cancel",
    "Fields": [
      {
        "Name": "ID",
        "Type": "string",
        "TagKey": "pos",
        "TagVal": "1",
        "Desc": "ID of the order to cancel"
      }
    ],
    "InjectLogger": true
  }
]
```

#### **Template-Based Generator**

Implement a separate standalone generator logic that:

1. Reads the all the commands from `commands.json` definition file.
2. Uses our custom JSON unmarshaller we implemented in Chapter 2 to parse it into an internal definition structure.
3. Feeds this metadata into a defined text template string using Go's standard `text/template` package:

   ```go
   const commandTemplate = `package internal

   type {{.CommandName}}Command struct {
   {{- range .Fields}}
   	{{.Name}} {{.Type}} ` + "`" + `{{.TagKey}}:"{{.TagVal}}" desc:"{{.Desc}}"` + "`" + `
   {{- end}}
   {{- if .InjectLogger}}
   	Logger *Logger ` + "`" + `di:"inject"` + "`" + `
   {{- end}}
   }
   `
   ```

4. The template must programmatically generate a complete, valid Go file (e.g., `generated_commands.go`), automatically appending the correct struct definitions, field types, custom tags (`cli`, `pos`, `desc`), and the `di:"inject"` tag for the logger if `InjectLogger` is set to true.

#### **Automation Integration**

Wire this generator up so it can be seamlessly triggered through native Go automation using the `//go:generate` directive.

## CLI Interface

### Chapter 1

For now we’ll have only one command called `load`:

```bash
# Adding cargo using mixed positional and flagged arguments
stargo load "Dilithium Crystals" --weight 450 --hazardous

# Requesting auto-generated help for a specific command
stargo load --help
```

### Chapter 2

The `load` command now handles persistence:

```bash
# If current total weight + 500,000 <= 4,000,000 -> Saves to cargo.json
stargo load "Plutonium Cells" --weight 500000 --hazardous

# If this pushes total weight over 4,000,000 -> Aborts with error
stargo load "Heavy Dark Matter" --weight 3600000
```

### Chapter 3

No changes to CLI interface in this chapter

### Chapter 4

```bash
# Triggers the template metaprogramming engine to read command.json and compile generated_commands.go
go generate ./...
```

### Success Criteria

- Adding a brand-new CLI command requires _only_ defining a new Go struct with appropriate tags and registering it in a central map. No changes to the parsing loop or manual variable bindings should be needed.
- The application gracefully handles and reports user errors, such as missing required positional arguments or passing invalid types (e.g., passing a string `"heavy"` to an integer `-weight` field) without panicking.
- The auto-generated help menu dynamically changes its output based on the struct definitions and tags alone.
- Running `stargo load` successfully persists data to `cargo.json` in JSON format using your custom reflection-based engine.
- Attempting to load cargo that pushes the total cumulative weight above 4,000,000 kg blocks the save operation and returns a user-friendly error message.
- The metaprogramming pipeline successfully reads an external structural JSON layout and synthesizes a compilable, syntax-valid Go source file containing fully formed struct structures.
- The generated structures contain identical structural properties and tags as manually typed ones, completely integrating with the Chapter 1 Parser, Chapter 2 Serializer, and Chapter 3 DI container.

### Constraints

- **Standard Library Only:** You must rely exclusively on built-in Go packages (`reflect`, `strconv`, `os`, `fmt`).
- Marshalling must be done without using 3rd party libs
- Unmarshal raw JSON into `any` using `encoding/json.Unmarshal` to avoid writing JSON tokenizer from scratch as it goes beyond the scope of this project
- Command structures have zero awareness of how loggers are instantiated; they request components structurally via `di:"inject"` tags.
- The DI container successfully finds registered services by type and maps them securely into command struct fields without type-mismatch panics.
- **No Third-Party DI Tools:** The wiring container logic and dependency assignment must be written manually using the `reflect` package without relying on external DI frameworks.
