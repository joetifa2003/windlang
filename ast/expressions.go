package ast

import (
	"bytes"
	"strings"
	"wind-vm-go/token"
)

type Identifier struct {
	Expression

	Token token.Token // the token.IDENT token
	Value string
}

func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

type IntegerLiteral struct {
	Expression

	Token token.Token
	Value int64
}

func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.TokenLiteral() }

type FloatLiteral struct {
	Expression

	Token token.Token
	Value float64
}

func (fl *FloatLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FloatLiteral) String() string       { return fl.TokenLiteral() }

type Boolean struct {
	Expression

	Token token.Token
	Value bool
}

func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) String() string       { return b.Token.Literal }

type PrefixExpression struct {
	Expression

	Token    token.Token // The prefix token, e.g. !
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String()
}

type InfixExpression struct {
	Expression

	Token    token.Token // The operator token, e.g. +
	Left     Expression
	Operator string
	Right    Expression
}

func (oe *InfixExpression) TokenLiteral() string { return oe.Token.Literal }
func (oe *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(oe.Left.String())
	out.WriteString(" " + oe.Operator + " ")
	out.WriteString(oe.Right.String())
	out.WriteString(")")
	return out.String()
}

type IfExpression struct {
	Expression

	Token      token.Token // The 'if' token
	Condition  Expression
	ThenBranch Statement
	ElseBranch Statement
}

func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) String() string {
	var out bytes.Buffer
	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.ThenBranch.String())

	if ie.ElseBranch != nil {
		out.WriteString("else ")
		out.WriteString(ie.ElseBranch.String())
	}

	return out.String()
}

type FunctionLiteral struct {
	Expression

	Token      token.Token // The 'fn' token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(fl.Body.String())

	return out.String()
}

type CallExpression struct {
	Expression

	Token     token.Token // The '(' token
	Function  Expression  // Identifier or FunctionLiteral
	Arguments []Expression
}

func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var out bytes.Buffer
	args := []string{}

	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

type StringLiteral struct {
	Expression

	Token token.Token
	Value string
}

func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return sl.Token.Literal }

type PostfixExpression struct {
	Expression

	Token    token.Token
	Left     Expression
	Operator string
}

func (pe *PostfixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PostfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Left.String())
	out.WriteString(" " + pe.Operator + ")")

	return out.String()
}

type AssignExpression struct {
	Expression

	Token token.Token
	Name  Expression
	Value Expression
}

func (ae *AssignExpression) TokenLiteral() string { return ae.Token.Literal }
func (ae *AssignExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ae.Name.String())
	out.WriteString(" = " + ae.Value.String() + ")")

	return out.String()
}

type ArrayLiteral struct {
	Expression

	Token token.Token
	Value []Expression
}

func (ae *ArrayLiteral) TokenLiteral() string { return ae.Token.Literal }
func (a *ArrayLiteral) Inspect() string {
	var out bytes.Buffer

	out.WriteString("[")
	for _, obj := range a.Value {
		out.WriteString(obj.String())
		out.WriteString(",")
	}
	out.WriteString("]")

	return out.String()
}

type IndexExpression struct {
	Expression

	Token token.Token
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IndexExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("])")

	return out.String()
}

type NilLiteral struct{ Expression, Token token.Token }

func (ne *NilLiteral) TokenLiteral() string { return ne.Token.Literal }
func (ne *NilLiteral) String() string       { return "nil" }

type HashLiteral struct {
	Expression

	Token token.Token
	Pairs map[Expression]Expression
}

func (hl *HashLiteral) TokenLiteral() string { return hl.Token.Literal }
func (hl *HashLiteral) String() string       { return "hash" }
