package vm

import "github.com/joetifa2003/windlang/value"

type Environment struct {
	Store []value.Value
}

func NewEnvironment(varCount int) Environment {
	store := make([]value.Value, varCount)
	return Environment{
		Store: store,
	}
}

type EnvironmentStack struct {
	Value []Environment
	p     int
}

func NewEnvironmentStack() EnvironmentStack {
	return EnvironmentStack{
		Value: make([]Environment, 2048),
	}
}

func (s *EnvironmentStack) pop() Environment {
	lastEle := (s.Value)[s.p-1]
	(s.Value)[s.p-1] = Environment{Store: nil}
	s.p--
	return lastEle
}

func (s *EnvironmentStack) push(env Environment) {
	s.Value[s.p] = env
	s.p++
}

func (s *EnvironmentStack) let(index int, val value.Value) {
	env := &s.Value[s.p-1]

	env.Store[index] = val
}

func (s *EnvironmentStack) get(scopeIndex, index int) value.Value {
	env := &s.Value[scopeIndex]

	return env.Store[index]
}

func (s *EnvironmentStack) set(scopeIndex, index int, value value.Value) value.Value {
	env := &s.Value[scopeIndex]
	env.Store[index] = value

	return env.Store[index]
}
