package parser

import (
	"fmt"
	"strconv"
	"wind-vm-go/ast"
	"wind-vm-go/lexer"
	"wind-vm-go/token"
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

const (
	_ int = iota
	LOWEST
	ASSIGN      // =
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	POSTFIX     // x++ or x--
	CALL        // myFunction(X)
)

type Parser struct {
	lexer *lexer.Lexer

	Errors []string

	curToken  token.Token
	peekToken token.Token
}

func New(l *lexer.Lexer) *Parser {
	p := Parser{lexer: l}

	p.nextToken()
	p.nextToken()

	return &p
}

func (p *Parser) getPrecedence(tokenType token.TokenType) int {
	switch tokenType {
	case token.EQ, token.NOT_EQ:
		return EQUALS
	case token.LT, token.GT:
		return LESSGREATER
	case token.PLUS, token.MINUS:
		return SUM
	case token.SLASH, token.ASTERISK, token.MODULO:
		return PRODUCT
	case token.LPAREN:
		return CALL
	case token.PLUSPLUS, token.MINUSMINUS:
		return POSTFIX
	}

	return LOWEST
}

func (p *Parser) getPrefixParseFn(tokenType token.TokenType) prefixParseFn {
	switch tokenType {
	case token.IDENT:
		return p.parseIdentifier
	case token.INT:
		return p.parseIntegerLiteral
	case token.BANG, token.MINUS:
		return p.parsePrefixExpression
	case token.TRUE, token.FALSE:
		return p.parseBoolean
	case token.LPAREN:
		return p.parseGroupedExpression
	case token.IF:
		return p.parseIfExpression
	case token.FUNCTION:
		return p.parseFunctionLiteral
	case token.STRING:
		return p.parseStringLiteral
	}

	return nil
}

func (p *Parser) getInfixParseFn(tokenType token.TokenType) infixParseFn {
	switch tokenType {
	case token.PLUS, token.MINUS, token.SLASH, token.ASTERISK, token.EQ, token.NOT_EQ, token.LT, token.GT, token.MODULO:
		return p.parseInfixExpression
	case token.LPAREN:
		return p.parseCallExpression
	case token.PLUSPLUS, token.MINUSMINUS:
		return p.parsePostfixExpression
	case token.ASSIGN:
		return p.parseAssignExpression
	}

	return nil
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := ast.Program{
		Statements: []ast.Statement{},
	}

	for p.curToken.Type != token.EOF {
		statement := p.parseStatement()

		if statement != nil {
			program.Statements = append(program.Statements, statement)
		}
	}

	return &program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	case token.FOR:
		return p.parseForStatement()
	case token.LBRACE:
		return p.parseBlockStatement()
	case token.INCLUDE:
		return p.parseIncludeStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() ast.Statement {
	stmt := ast.LetStatement{Token: p.curToken}

	p.nextToken()

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectCurrent(token.IDENT) {
		return nil
	}

	if !p.expectCurrent(token.ASSIGN) {
		return nil
	}

	stmt.Value = p.parseExpression(LOWEST)

	if !p.expectCurrent(token.SEMICOLON) {
		return nil
	}

	return &stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)

	if !p.expectCurrent(token.SEMICOLON) {
		return nil
	}

	return &stmt
}

func (p *Parser) parseForStatement() ast.Statement {
	stmt := ast.ForStatement{Token: p.curToken}

	p.nextToken()

	if !p.expectCurrent(token.LPAREN) {
		return nil
	}

	stmt.Initializer = p.parseStatement()

	stmt.Condition = p.parseExpression(LOWEST)

	if !p.expectCurrent(token.SEMICOLON) {
		return nil
	}

	stmt.Increment = p.parseExpression(LOWEST)

	if !p.expectCurrent(token.RPAREN) {
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	return &stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.currentTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return &stmt
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	if !p.expectCurrent(token.LBRACE) {
		return nil
	}

	for !p.currentTokenIs(token.RBRACE) && !p.currentTokenIs(token.EOF) {
		stmt := p.parseStatement()

		block.Statements = append(block.Statements, stmt)
	}

	if !p.expectCurrent(token.RBRACE) {
		return nil
	}

	return &block
}

func (p *Parser) parseIncludeStatement() *ast.IncludeStatement {
	stmt := ast.IncludeStatement{Token: p.curToken}

	p.nextToken()

	stmt.Path = p.curToken.Literal

	p.nextToken()

	if !p.expectCurrent(token.SEMICOLON) {
		return nil
	}

	return &stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.getPrefixParseFn(p.curToken.Type)
	if prefix == nil {
		msg := fmt.Sprintf("cannot parse %s as an expression", p.curToken.Literal)
		p.Errors = append(p.Errors, msg)
		return nil
	}

	leftExp := prefix()

	for !p.currentTokenIs(token.SEMICOLON) && p.curPrecedence() >= precedence {
		infix := p.getInfixParseFn(p.curToken.Type)
		if infix == nil {
			return leftExp
		}

		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()

	p.nextToken()

	expression.Right = p.parseExpression(precedence)

	return &expression
}

func (p *Parser) parseIdentifier() ast.Expression {
	ident := ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	p.nextToken()

	return &ident
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	integer := ast.IntegerLiteral{Token: p.curToken}
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)

	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.Errors = append(p.Errors, msg)
		return nil
	}

	integer.Value = value

	p.nextToken()

	return &integer
}

func (p *Parser) parseStringLiteral() ast.Expression {
	str := ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}

	p.nextToken()

	return &str
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()

	expression.Right = p.parseExpression(PREFIX)

	return &expression
}

func (p *Parser) parseBoolean() ast.Expression {
	expr := ast.Boolean{Token: p.curToken, Value: p.currentTokenIs(token.TRUE)}

	p.nextToken()

	return &expr
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}

func (p *Parser) parseIfExpression() ast.Expression {
	expression := ast.IfExpression{Token: p.curToken}

	p.nextToken()

	if !p.expectCurrent(token.LPAREN) {
		return nil
	}

	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectCurrent(token.RPAREN) {
		return nil
	}

	thenStatement := p.parseStatement()
	expression.ThenBranch = thenStatement

	if p.currentTokenIs(token.ELSE) {
		p.nextToken()

		elseStatement := p.parseStatement()
		expression.ElseBranch = elseStatement
	}

	return &expression

}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := ast.FunctionLiteral{Token: p.curToken}

	p.nextToken()

	if !p.expectCurrent(token.LPAREN) {
		return nil
	}

	lit.Parameters = p.parseFunctionParameters()

	lit.Body = p.parseBlockStatement()

	return &lit
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	if p.currentTokenIs(token.RPAREN) {
		p.nextToken()
		return identifiers
	}

	for p.peekTokenIs(token.COMMA) {
		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, ident)

		p.nextToken() // consume IDENT
		p.nextToken() // consume COMMA
	}

	// last IDENT
	ident := ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	identifiers = append(identifiers, &ident)
	p.nextToken()

	if !p.expectCurrent(token.RPAREN) {
		return nil
	}

	return identifiers
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := ast.CallExpression{Token: p.curToken, Function: function}

	p.nextToken()

	exp.Arguments = p.parseCallArguments()

	return &exp
}

func (p *Parser) parseCallArguments() []ast.Expression {
	args := []ast.Expression{}

	if p.currentTokenIs(token.RPAREN) {
		p.nextToken()
		return args
	}

	args = append(args, p.parseExpression(LOWEST))

	for p.currentTokenIs(token.COMMA) {
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}

	if !p.expectCurrent(token.RPAREN) {
		return nil
	}

	return args
}

func (p *Parser) parsePostfixExpression(left ast.Expression) ast.Expression {
	expression := ast.PostfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	p.nextToken()

	return &expression
}

func (p *Parser) parseAssignExpression(left ast.Expression) ast.Expression {
	ident, ok := left.(*ast.Identifier)
	if !ok {
		msg := fmt.Sprintf("left-hand side of assignment must be an identifier, got %T", left)
		p.Errors = append(p.Errors, msg)
		return nil
	}

	expression := ast.AssignExpression{
		Token: p.curToken,
		Name:  ident,
	}

	precedence := p.curPrecedence()

	p.nextToken() // consume ASSIGN

	expression.Value = p.parseExpression(precedence)

	return &expression
}

func (p *Parser) curPrecedence() int {
	return p.getPrecedence(p.curToken.Type)
}

func (p *Parser) expectPeek(tokenType token.TokenType) bool {
	if p.peekTokenIs(tokenType) {
		p.nextToken()
		return true
	}

	p.peekError(tokenType)
	return false
}

func (p *Parser) expectCurrent(tokenType token.TokenType) bool {
	if p.currentTokenIs(tokenType) {
		p.nextToken()
		return true
	}

	p.currentError(tokenType)
	return false
}

func (p *Parser) peekTokenIs(tokenType token.TokenType) bool {
	return p.peekToken.Type == tokenType
}

func (p *Parser) currentTokenIs(tokenType token.TokenType) bool {
	return p.curToken.Type == tokenType
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t.String(), p.peekToken.Type.String())

	p.Errors = append(p.Errors, msg)
}

func (p *Parser) currentError(t token.TokenType) {
	msg := fmt.Sprintf("expected token to be %s, got %s instead",
		t.String(), p.curToken.Type.String())

	p.Errors = append(p.Errors, msg)
}
