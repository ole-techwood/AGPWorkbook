package internal

import (
	"fmt"
	"reflect"
	"sort"
)

// Container is a runtime DI registry mapping concrete types to singleton instances.
// Not safe for concurrent use.
type Container struct {
	registry map[reflect.Type]reflect.Value
}

// NewContainer creates a new Container with an empty registry.
func NewContainer() *Container {
	return &Container{
		registry: make(map[reflect.Type]reflect.Value),
	}
}

// Register stores dependency in the container keyed by its concrete runtime type.
// Silently ignores untyped nil. Panics on typed nil pointers.
// Re-registering the same concrete type replaces the prior entry.
func (c *Container) Register(dependency any) {
	if dependency == nil {
		return
	}

	v := reflect.ValueOf(dependency)

	// Typed nil pointer: fail fast — invalid singleton value.
	if v.Kind() == reflect.Pointer && v.IsNil() {
		panic("di: typed nil pointer is not a valid dependency")
	}

	c.registry[v.Type()] = v
}

// Resolve returns the registered singleton whose concrete type implements ifaceType.
// ifaceType must be an interface type obtained via reflect.TypeOf((*T)(nil)).Elem().
// Returns an error when no match exists or multiple registered types match.
// Ambiguity errors list candidate type names in stable alphabetical order.
func (c *Container) Resolve(ifaceType reflect.Type) (reflect.Value, error) {
	if ifaceType.Kind() != reflect.Interface {
		return reflect.Value{}, fmt.Errorf("di: resolve expects interface type, got %s", ifaceType)
	}

	var matches []reflect.Type
	for t := range c.registry {
		if t.Implements(ifaceType) {
			matches = append(matches, t)
		}
	}

	if len(matches) == 0 {
		return reflect.Value{}, fmt.Errorf("di: no registered type implements %s", ifaceType)
	}

	if len(matches) > 1 {
		names := make([]string, len(matches))
		for i, t := range matches {
			names[i] = t.String()
		}
		sort.Strings(names)
		return reflect.Value{}, fmt.Errorf("di: ambiguous resolution for %s: candidates %v", ifaceType, names)
	}

	return c.registry[matches[0]], nil
}
