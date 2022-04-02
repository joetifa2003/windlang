package ast

import (
	"bytes"
	"wind-vm-go/token"
)

type Program struct {
	Statements []Statement
}

func (p Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (p Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
		out.WriteString("\n")
	}

	return out.String()
}

type LetStatement struct {
	Statement

	Token token.Token // the token.LET token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	out.WriteString(";")
	return out.String()
}

type ReturnStatement struct {
	Statement

	Token       token.Token // the 'return' token
	ReturnValue Expression
}

func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}

	out.WriteString(";")
	return out.String()
}

type ExpressionStatement struct {
	Statement

	Token      token.Token // the first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}

	return ""
}

type BlockStatement struct {
	Statement

	Token      token.Token // the { token
	Statements []Statement
}

func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

type ForStatement struct {
	Statement

	Token       token.Token // the 'for' token
	Initializer Statement
	Condition   Expression
	Increment   Expression
	Body        Statement
}

func (fs *ForStatement) TokenLiteral() string { return fs.Token.Literal }
func (fs *ForStatement) String() string {
	return ""
}

type IncludeStatement struct {
	Statement

	Token token.Token // the 'include' token
	Path  string
	Alias *Identifier
}

func (is *IncludeStatement) TokenLiteral() string { return is.Token.Literal }
func (is *IncludeStatement) String() string {
	var out bytes.Buffer

	out.WriteString(is.TokenLiteral() + " ")
	out.WriteString(is.Path)
	out.WriteString(";")

	return out.String()
}

type WhileStatement struct {
	Statement

	Token     token.Token
	Condition Expression
	Body      Statement
}

func (ws *WhileStatement) TokenLiteral() string { return ws.Token.Literal }
func (ws *WhileStatement) String() string {
	return ""
}
