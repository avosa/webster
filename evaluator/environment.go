package evaluator

import (
	"fmt"

	"github.com/avosa/webster/object"
)

type Environment struct {
	store map[string]object.Object
	outer *Environment
}

func NewEnvironment() *Environment {
	s := make(map[string]object.Object)
	return &Environment{store: s, outer: nil}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func (e *Environment) Get(name string) (object.Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

func (e *Environment) Set(name string, val object.Object) object.Object {
	e.store[name] = val
	return val
}

func (e *Environment) String() string {
	var out string
	for k, v := range e.store {
		out += fmt.Sprintf("%s: %s\n", k, v.Inspect())
	}
	if e.outer != nil {
		out += e.outer.String()
	}
	return out
}
