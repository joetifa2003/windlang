package evaluator

type Environment struct {
	Store           map[string]Object
	ConstantStore   map[string]Object
	Outer           *Environment
	Includes        []*Environment
	IncludesAliased map[string]*IncludeObject
}

func NewEnvironment() *Environment {
	return &Environment{
		Store:           nil,
		ConstantStore:   nil,
		Outer:           nil,
		Includes:        nil,
		IncludesAliased: nil,
	}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.Outer = outer

	return env
}

func (e *Environment) Get(name string) (Object, bool) {
	if e.Store == nil && e.ConstantStore == nil {
		if e.Outer != nil {
			return e.Outer.Get(name)
		} else {
			return e.getIncludes(name)
		}
	}

	obj, ok := e.ConstantStore[name]
	if ok {
		return obj, ok
	}

	obj, ok = e.Store[name]
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
	aliasedInclude, ok := e.IncludesAliased[name]
	if ok {
		return aliasedInclude, true
	}

	for _, include := range e.Includes {
		if obj, ok := include.Store[name]; ok {
			return obj, ok
		}
	}

	return nil, false
}

// Set For assigning
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

// Let For local scope variables
func (e *Environment) Let(name string, val Object) Object {
	if e.Store == nil {
		e.Store = make(map[string]Object)
	}

	e.Store[name] = val

	return val
}

func (e *Environment) LetConstant(name string, val Object) Object {
	if e.ConstantStore == nil {
		e.ConstantStore = make(map[string]Object)
	}

	e.ConstantStore[name] = val

	return val
}

func (e *Environment) IsConstant(name string) bool {
	if e.ConstantStore == nil {
		if e.Outer != nil {
			return e.Outer.IsConstant(name)
		} else {
			return false
		}
	}

	_, ok := e.ConstantStore[name]
	if ok {
		return true
	} else {
		if e.Outer != nil {
			return e.Outer.IsConstant(name)
		} else {
			return false
		}
	}
}

func (e *Environment) ClearStore() {
	for k := range e.Store {
		delete(e.Store, k)
	}
}

func (e *Environment) AddAlias(name string, include *IncludeObject) {
	if e.IncludesAliased == nil {
		e.IncludesAliased = make(map[string]*IncludeObject)
	}

	e.IncludesAliased[name] = include
}
