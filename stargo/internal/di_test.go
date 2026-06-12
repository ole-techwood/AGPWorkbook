package internal

import (
	"reflect"
	"strings"
	"testing"
)

type testService struct {
	Name string
}

type testCommand struct {
	Service *testService `di:"inject"`
}

type testCommandWithUntaggedField struct {
	Service *testService
}

type testCommandWithPrivateField struct {
	service *testService `di:"inject"`
}

type testLoggerInterface interface {
	Info(format string, args ...interface{})
}

type testCommandWithInterfaceField struct {
	Logger testLoggerInterface `di:"inject"`
}

type testGreeter interface {
	Greet() string
}

type testGreeterA struct{}

type testGreeterB struct{}

func (g *testGreeterA) Greet() string {
	return "a"
}

func (g *testGreeterB) Greet() string {
	return "b"
}

type testCommandWithAmbiguousInterface struct {
	Greeter testGreeter `di:"inject"`
}

type otherService struct{}

func TestInjectRejectsInvalidTargets(t *testing.T) {
	container := NewContainer()

	t.Run("nil target", func(t *testing.T) {
		err := container.Inject(nil)
		if err == nil {
			t.Fatal("expected error for nil target")
		}
	})

	t.Run("non-pointer target", func(t *testing.T) {
		err := container.Inject(testCommand{})
		if err == nil {
			t.Fatal("expected error for non-pointer target")
		}
	})

	t.Run("nil pointer target", func(t *testing.T) {
		var target *testCommand
		err := container.Inject(target)
		if err == nil {
			t.Fatal("expected error for nil pointer target")
		}
	})

	t.Run("pointer to non-struct", func(t *testing.T) {
		value := 7
		err := container.Inject(&value)
		if err == nil {
			t.Fatal("expected error for pointer-to-non-struct target")
		}
	})
}

func TestInjectSetsTaggedFieldByConcreteType(t *testing.T) {
	container := NewContainer()
	service := &testService{Name: "cargo"}
	container.Register(service)

	cmd := &testCommand{}
	if err := container.Inject(cmd); err != nil {
		t.Fatalf("inject returned error: %v", err)
	}

	if cmd.Service != service {
		t.Fatal("expected tagged field to be injected with registered singleton")
	}
}

func TestInjectLeavesUntaggedFieldUnchanged(t *testing.T) {
	container := NewContainer()
	container.Register(&testService{Name: "cargo"})

	cmd := &testCommandWithUntaggedField{}
	if err := container.Inject(cmd); err != nil {
		t.Fatalf("inject returned error: %v", err)
	}

	if cmd.Service != nil {
		t.Fatal("expected untagged field to remain unchanged")
	}
}

func TestInjectUsesInterfaceResolution(t *testing.T) {
	container := NewContainer()
	logger := NewLogger()
	container.Register(logger)

	cmd := &testCommandWithInterfaceField{}
	if err := container.Inject(cmd); err != nil {
		t.Fatalf("inject returned error: %v", err)
	}

	if cmd.Logger == nil {
		t.Fatal("expected interface field to be injected")
	}
}

func TestInjectFailsOnUnsettableTaggedField(t *testing.T) {
	container := NewContainer()
	container.Register(&testService{Name: "cargo"})

	cmd := &testCommandWithPrivateField{}
	err := container.Inject(cmd)
	if err == nil {
		t.Fatal("expected error for unsettable tagged field")
	}

	if !strings.Contains(err.Error(), "not settable") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestInjectFailsWhenDependencyMissing(t *testing.T) {
	container := NewContainer()

	cmd := &testCommand{}
	err := container.Inject(cmd)
	if err == nil {
		t.Fatal("expected missing dependency error")
	}

	if !strings.Contains(err.Error(), "no dependency registered") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestInjectFailsOnTypeMismatch(t *testing.T) {
	container := NewContainer()
	container.registry[reflect.TypeFor[*testService]()] = reflect.ValueOf(&otherService{})

	cmd := &testCommand{}
	err := container.Inject(cmd)
	if err == nil {
		t.Fatal("expected type mismatch error")
	}

	if !strings.Contains(err.Error(), "not assignable") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestInjectFailsOnAmbiguousInterfaceResolution(t *testing.T) {
	container := NewContainer()
	container.Register(&testGreeterA{})
	container.Register(&testGreeterB{})

	cmd := &testCommandWithAmbiguousInterface{}
	err := container.Inject(cmd)
	if err == nil {
		t.Fatal("expected ambiguous interface resolution error")
	}

	if !strings.Contains(err.Error(), "ambiguous") {
		t.Fatalf("unexpected error: %v", err)
	}
}
