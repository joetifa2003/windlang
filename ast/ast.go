package ast

type Node interface {
	TokenLiteral() string
}

type Statement interface{ Node }
type Expression interface {
	Node
	String() string
}
