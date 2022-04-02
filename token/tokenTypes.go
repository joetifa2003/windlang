package token

type TokenType int

const (
	ILLEGAL TokenType = iota
	EOF

	// Identifiers + literals
	IDENT // add, foobar, x, y, ...
	INT   // 1343456
	STRING
	FLOAT

	// Operators
	ASSIGN
	PLUS
	MINUS
	BANG
	ASTERISK
	SLASH
	MODULO // %
	LT     // <
	GT     // >
	LT_EQ  // <=
	GT_EQ  // >=
	AND    // &&
	OR     // ||
	EQ     // ==
	NOT_EQ // !=
	PLUSPLUS
	MINUSMINUS

	// Delimiters
	COMMA
	COLON
	DOT
	SEMICOLON
	LPAREN
	RPAREN
	LBRACE
	RBRACE
	LBRACKET
	RBRACKET

	// Keywords
	FUNCTION
	LET
	TRUE
	FALSE
	IF
	ELSE
	RETURN
	FOR
	INCLUDE
	AS
	WHILE
	NIL
)

func (t *TokenType) String() string {
	switch *t {
	case ILLEGAL:
		return "ILLEGAL"
	case EOF:
		return "EOF"
	case IDENT:
		return "IDENT"
	case INT:
		return "INT"
	case ASSIGN:
		return "="
	case PLUS:
		return "+"
	case MINUS:
		return "-"
	case BANG:
		return "!"
	case ASTERISK:
		return "*"
	case SLASH:
		return "/"
	case LT:
		return "<"
	case GT:
		return ">"
	case EQ:
		return "=="
	case NOT_EQ:
		return "!="
	case COMMA:
		return ","
	case SEMICOLON:
		return ";"
	case LPAREN:
		return "("
	case RPAREN:
		return ")"
	case LBRACE:
		return "{"
	case RBRACE:
		return "}"
	case FUNCTION:
		return "FUNCTION"
	case LET:
		return "LET"
	case TRUE:
		return "TRUE"
	case FALSE:
		return "FALSE"
	case IF:
		return "IF"
	case ELSE:
		return "ELSE"
	case RETURN:
		return "RETURN"
	case STRING:
		return "STRING"
	case FOR:
		return "FOR"
	case PLUSPLUS:
		return "++"
	case MINUSMINUS:
		return "--"
	case INCLUDE:
		return "INCLUDE"
	case WHILE:
		return "WHILE"
	case FLOAT:
		return "FLOAT"
	case NIL:
		return "NIL"
	case AS:
		return "AS"
	default:
		return "UNKNOWN"
	}
}
