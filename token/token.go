package token

type Token struct {
	Type    TokenType
	Literal string
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
	}

	return IDENT, false
}
