package vm

import "github.com/joetifa2003/windlang/value"

type Stack struct {
	Value []value.Value
	P     int
}

func NewStack() Stack {
	return Stack{
		Value: make([]value.Value, 2048),
	}
}

func (s *Stack) pop() value.Value {
	lastEle := (s.Value)[s.P-1]
	(s.Value)[s.P-1] = nil
	s.P--
	return lastEle
}

func (s *Stack) push(value value.Value) {
	s.Value[s.P] = value
	s.P++
}
