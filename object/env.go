package object

type Environment struct {
	Store    map[string]Object
	Outer    *Environment
	Includes []*Environment
}

func NewEnvironment() *Environment {
	s := make(map[string]Object)

	return &Environment{Store: s}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.Outer = outer

	return env
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.Store[name]
	if !ok && e.Outer != nil {
		obj, ok = e.Outer.Get(name)
	}

	// If not found in the local env it will search in the includes
	if !ok {
		for _, include := range e.Includes {
			if obj, ok := include.Store[name]; ok {
				return obj, ok
			}
		}
	}

	return obj, ok
}

// For assigning
func (e *Environment) Set(name string, val Object) Object {
	_, ok := e.Store[name]
	if !ok && e.Outer != nil {
		return e.Outer.Set(name, val)
	}

	e.Store[name] = val

	return val
}

// For local scope variables
func (e *Environment) Let(name string, val Object) Object {
	e.Store[name] = val

	return val
}
