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
	LET_INT = "LET_INT"

)

type TokenType string

type Token struct {
	Type TokenType
	Literal string
}

var keywords = map[string] TokenType {
	"func": FUNCTION,
	"int": LET_INT,
}

// LookupIdent checks the keywords table to see if identifier is a keyword
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
