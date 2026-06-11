package internal

import (
	"fmt"
	"os"
	"reflect"
)

// Cargo represents a persisted cargo item.
// Extends LoadCommand with an auto-incremented ID.
type Cargo struct {
	ID string

	LoadCommand
}

const (
	cargoDataDir  = "data"
	cargoFilePath = "data/cargo.json"
	weightLimit   = 4_000_000
)

func LoadCargoList(path string) ([]Cargo, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return []Cargo{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("loadCargoList: %w", err)
	}

	val, err := UnmarshalFromJSON(data, reflect.TypeFor[[]Cargo]())
	if err != nil {
		return nil, fmt.Errorf("loadCargoList: %w", err)
	}

	return val.Interface().([]Cargo), nil
}

func SaveCargoList(path string, items []Cargo) error {
	data, err := MarshalToJSON(items)
	if err != nil {
		return fmt.Errorf("saveCargoList: %w", err)
	}

	return os.WriteFile(path, data, 0o644)
}
