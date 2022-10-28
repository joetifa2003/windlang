package token

type Token struct {
	Type    TokenType
	Literal string
	Line    int
}

func LookupIdent(ident string) TokenType {
	if tok, ok := isKeyword(ident); ok {
		return tok
	}

	return IDENT
}

func isKeyword(ident string) (TokenType, bool) {
	switch ident {
	case "fn":
		return FUNCTION, true
	case "let":
		return LET, true
	case "const":
		return CONST, true
	case "true":
		return TRUE, true
	case "false":
		return FALSE, true
	case "if":
		return IF, true
	case "else":
		return ELSE, true
	case "return":
		return RETURN, true
	case "for":
		return FOR, true
	case "include":
		return INCLUDE, true
	case "while":
		return WHILE, true
	case "nil":
		return NIL, true
	case "as":
		return AS, true
	case "break":
		return BREAK, true
	case "continue":
		return CONTINUE, true
	case "echo":
		return ECHO, true
	}

	return IDENT, false
}
