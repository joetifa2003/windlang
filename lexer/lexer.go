package lexer

import (
	"wind-vm-go/token"
)

type Lexer struct {
	input        []rune
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           rune // current char under examination
}

func New(input string) *Lexer {
	l := Lexer{input: []rune(input)}
	l.readChar()
	return &l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.Token{Type: token.EQ, Literal: "=="}
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case '&':
		if l.peekChar() == '&' {
			l.readChar()

			tok = token.Token{Type: token.AND, Literal: "&&"}
		} else {
			tok = token.Token{Type: token.ILLEGAL, Literal: string(l.ch) + string(l.peekChar())}
		}
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '%':
		tok = newToken(token.MODULO, l.ch)
	case '+':
		if l.peekChar() == '+' {
			l.readChar()
			tok = token.Token{Type: token.PLUSPLUS, Literal: "++"}
		} else {
			tok = newToken(token.PLUS, l.ch)
		}
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '!':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.Token{Type: token.NOT_EQ, Literal: "!="}
		} else {
			tok = newToken(token.BANG, l.ch)
		}
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '/':
		if l.peekChar() == '/' {
			l.readChar()

			// Skip comment
			for l.ch != '\n' && l.ch != 0 {
				l.readChar()
			}

			tok = l.NextToken()
		} else {
			tok = newToken(token.SLASH, l.ch)
		}
	case '<':
		if l.peekChar() == '=' {
			l.readChar()

			tok = token.Token{Type: token.LT_EQ, Literal: "<="}
		} else {
			tok = newToken(token.LT, l.ch)
		}
	case '>':
		if l.peekChar() == '=' {
			l.readChar()

			tok = token.Token{Type: token.GT_EQ, Literal: ">="}
		} else {
			tok = newToken(token.GT, l.ch)
		}
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)

			return tok
		} else if isDigit(l.ch) {
			tok.Literal, tok.Type = l.readNumber()

			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()

	return tok
}

func (l *Lexer) readIdentifier() string {
	position := l.position

	for isLetter(l.ch) {
		l.readChar()
	}

	return string(l.input[position:l.position])
}

func (l *Lexer) readString() string {
	position := l.position + 1

	for {
		l.readChar()

		if l.ch == '"' || l.ch == 0 {
			break
		}
	}

	return string(l.input[position:l.position])
}

func (l *Lexer) readNumber() (string, token.TokenType) {
	position := l.position

	dotCount := 0
	for isDigit(l.ch) || l.ch == '.' {
		l.readChar()

		if l.ch == '.' {
			dotCount++
		}
	}

	if dotCount == 1 {
		return string(l.input[position:l.position]), token.FLOAT
	} else {
		return string(l.input[position:l.position]), token.INT
	}
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) peekChar() rune {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func newToken(tokenType token.TokenType, ch rune) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func isLetter(ch rune) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}
