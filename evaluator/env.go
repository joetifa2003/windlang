package evaluator

import (
	"sync"
)

var envPool = sync.Pool{
	New: func() interface{} {
		return &Environment{}
	},
}

type Environment struct {
	Store    map[string]Object
	Outer    *Environment
	Includes []*Environment
}

func NewEnvironment() *Environment {
	return &Environment{}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := envPool.Get().(*Environment)
	env.Outer = outer

	return env
}

func (e *Environment) Get(name string) (Object, bool) {
	if e.Store == nil {
		if e.Outer != nil {
			return e.Outer.Get(name)
		} else {
			return e.getIncludes(name)
		}
	}

	obj, ok := e.Store[name]
	if !ok {
		if e.Outer != nil {
			obj, ok = e.Outer.Get(name)
		} else {
			return e.getIncludes(name)
		}
	}

	return obj, ok
}

func (e *Environment) getIncludes(name string) (Object, bool) {
	for _, include := range e.Includes {
		if obj, ok := include.Store[name]; ok {
			return obj, ok
		}
	}

	return nil, false
}

// For assigning
func (e *Environment) Set(name string, val Object) (Object, bool) {
	if e.Store == nil {
		if e.Outer != nil {
			return e.Outer.Set(name, val)
		} else {
			return nil, false
		}
	}

	_, ok := e.Store[name]
	if !ok {
		if e.Outer != nil {
			return e.Outer.Set(name, val)
		} else {
			return nil, false
		}
	}

	e.Store[name] = val

	return val, true
}

// For local scope variables
func (e *Environment) Let(name string, val Object) Object {
	if e.Store == nil {
		e.Store = make(map[string]Object)
	}

	e.Store[name] = val

	return val
}

func (e *Environment) ClearStore() {
	for k := range e.Store {
		delete(e.Store, k)
	}
}

func (e *Environment) Dispose() {
	envPool.Put(e)
}
