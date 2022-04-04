package lexer

import (
	"testing"

	"github.com/joetifa2003/windlang/token"

	"github.com/stretchr/testify/assert"
)

var nextTokenTestCases = []struct {
	input          string
	expectedTokens []token.Token
}{
	{
		input: "(){},;",
		expectedTokens: []token.Token{
			{Type: token.LPAREN, Literal: "(", Line: 1},
			{Type: token.RPAREN, Literal: ")", Line: 1},
			{Type: token.LBRACE, Literal: "{", Line: 1},
			{Type: token.RBRACE, Literal: "}", Line: 1},
			{Type: token.COMMA, Literal: ",", Line: 1},
			{Type: token.SEMICOLON, Literal: ";", Line: 1},
			{Type: token.EOF, Literal: "", Line: 1},
		},
	},
	{
		input: "+-/*",
		expectedTokens: []token.Token{
			{Type: token.PLUS, Literal: "+", Line: 1},
			{Type: token.MINUS, Literal: "-", Line: 1},
			{Type: token.SLASH, Literal: "/", Line: 1},
			{Type: token.ASTERISK, Literal: "*", Line: 1},
			{Type: token.EOF, Literal: "", Line: 1},
		},
	},
	{
		input: "!= == > < <= >=",
		expectedTokens: []token.Token{
			{Type: token.NOT_EQ, Literal: "!=", Line: 1},
			{Type: token.EQ, Literal: "==", Line: 1},
			{Type: token.GT, Literal: ">", Line: 1},
			{Type: token.LT, Literal: "<", Line: 1},
			{Type: token.LT_EQ, Literal: "<=", Line: 1},
			{Type: token.GT_EQ, Literal: ">=", Line: 1},
			{Type: token.EOF, Literal: "", Line: 1},
		},
	},
	{
		input: `true false 1 3.14 "hello" x`,
		expectedTokens: []token.Token{
			{Type: token.TRUE, Literal: "true", Line: 1},
			{Type: token.FALSE, Literal: "false", Line: 1},
			{Type: token.INT, Literal: "1", Line: 1},
			{Type: token.FLOAT, Literal: "3.14", Line: 1},
			{Type: token.STRING, Literal: "hello", Line: 1},
			{Type: token.IDENT, Literal: "x", Line: 1},
			{Type: token.EOF, Literal: "", Line: 1},
		},
	},
}

func TestNextToken(t *testing.T) {
	assert := assert.New(t)

	for _, testCase := range nextTokenTestCases {
		lexer := New(testCase.input)

		for _, expectedToken := range testCase.expectedTokens {
			actualToken := lexer.NextToken()
			assert.Equal(expectedToken, actualToken)
		}
	}
}
