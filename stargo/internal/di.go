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

// Inject populates struct fields tagged with `di:"inject"` using registered
// singleton dependencies from the container.
func (c *Container) Inject(target any) error {
	if c == nil {
		return fmt.Errorf("di: container is nil")
	}

	if target == nil {
		return fmt.Errorf("di: inject target must be non-nil pointer to struct")
	}

	targetValue := reflect.ValueOf(target)
	if targetValue.Kind() != reflect.Pointer {
		return fmt.Errorf("di: inject target must be pointer to struct, got %s", targetValue.Kind())
	}

	if targetValue.IsNil() {
		return fmt.Errorf("di: inject target pointer must be non-nil")
	}

	structValue := targetValue.Elem()
	if structValue.Kind() != reflect.Struct {
		return fmt.Errorf("di: inject target must point to struct, got %s", structValue.Kind())
	}

	structType := structValue.Type()

	for i := range structType.NumField() {
		fieldType := structType.Field(i)
		if fieldType.Tag.Get("di") != "inject" {
			continue
		}

		fieldValue := structValue.Field(i)
		if !fieldValue.CanSet() {
			return fmt.Errorf("di: inject field %s on %s is not settable", fieldType.Name, structType)
		}

		resolvedValue, err := c.resolveByFieldType(fieldType.Type)
		if err != nil {
			return fmt.Errorf("di: inject field %s (%s): %w", fieldType.Name, fieldType.Type, err)
		}

		if !resolvedValue.Type().AssignableTo(fieldValue.Type()) {
			return fmt.Errorf("di: inject field %s (%s): resolved %s is not assignable", fieldType.Name, fieldType.Type, resolvedValue.Type())
		}

		fieldValue.Set(resolvedValue)
	}

	return nil
}

func (c *Container) resolveByFieldType(fieldType reflect.Type) (reflect.Value, error) {
	if value, ok := c.registry[fieldType]; ok {
		return value, nil
	}

	if fieldType.Kind() == reflect.Interface {
		return c.Resolve(fieldType)
	}

	return reflect.Value{}, fmt.Errorf("di: no dependency registered for %s", fieldType)
}
