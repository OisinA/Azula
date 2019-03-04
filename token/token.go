package token

const (
	
	ILLEGAL = "ILLEGAL"
	EOF = "EOF"

	IDENT = "IDENT" //identifier (x, y)
	INT = "INT" //integer

	ASSIGN = "="
	PLUS = "+"

	COMMA = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	FUNCTION = "FUNCTION"
	LET = "LET"

)

type TokenType string

type Token struct {
	Type TokenType
	Literal string
}
