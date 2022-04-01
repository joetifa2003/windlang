package lexer

import (
	"testing"
	"wind-vm-go/token"

	"github.com/stretchr/testify/assert"
)

var nextTokenTestCases = []struct {
	input          string
	expectedTokens []token.Token
}{
	{
		input: "(){},;",
		expectedTokens: []token.Token{
			{Type: token.LPAREN, Literal: "("},
			{Type: token.RPAREN, Literal: ")"},
			{Type: token.LBRACE, Literal: "{"},
			{Type: token.RBRACE, Literal: "}"},
			{Type: token.COMMA, Literal: ","},
			{Type: token.SEMICOLON, Literal: ";"},
			{Type: token.EOF, Literal: ""},
		},
	},
	{
		input: "+-/*",
		expectedTokens: []token.Token{
			{Type: token.PLUS, Literal: "+"},
			{Type: token.MINUS, Literal: "-"},
			{Type: token.SLASH, Literal: "/"},
			{Type: token.ASTERISK, Literal: "*"},
			{Type: token.EOF, Literal: ""},
		},
	},
	{
		input: "!= == > < <= >=",
		expectedTokens: []token.Token{
			{Type: token.NOT_EQ, Literal: "!="},
			{Type: token.EQ, Literal: "=="},
			{Type: token.GT, Literal: ">"},
			{Type: token.LT, Literal: "<"},
			{Type: token.LT_EQ, Literal: "<="},
			{Type: token.GT_EQ, Literal: ">="},
			{Type: token.EOF, Literal: ""},
		},
	},
	{
		input: `true false 1 3.14 "hello" x`,
		expectedTokens: []token.Token{
			{Type: token.TRUE, Literal: "true"},
			{Type: token.FALSE, Literal: "false"},
			{Type: token.INT, Literal: "1"},
			{Type: token.FLOAT, Literal: "3.14"},
			{Type: token.STRING, Literal: "hello"},
			{Type: token.IDENT, Literal: "x"},
			{Type: token.EOF, Literal: ""},
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
