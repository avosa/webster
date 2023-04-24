package evaluator

import (
	"fmt"
	"strings"
)

// Environment is a map of string keys to Objects.
type Environment struct {
	store map[string]Object
	outer *Environment
}

// NewEnvironment creates a new Environment.
func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s, outer: nil}
}

// NewEnclosedEnvironment creates a new Environment with an outer Environment.
func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

// Get retrieves the value of a key from the Environment.
func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

// Set sets the value of a key in the Environment.
func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}

// String returns a string representation of the Environment.
func (e *Environment) String() string {
	var out []string
	for name := range e.store {
		out = append(out, name)
	}
	if e.outer != nil {
		out = append(out, e.outer.String())
	}
	return fmt.Sprintf("{%s}", strings.Join(out, ", "))
}
