package vm

import (
	"github.com/joetifa2003/windlang/value"
)

type Stack struct {
	Value []value.Value
}

func NewStack() Stack {
	return Stack{
		Value: make([]value.Value, 0, 2048),
	}
}

func (s *Stack) pop() value.Value {
	lastEle := s.Value[len(s.Value)-1]
	s.Value = s.Value[:len(s.Value)-1]
	return lastEle
}

func (s *Stack) push(value value.Value) {
	s.Value = append(s.Value, value)
}
