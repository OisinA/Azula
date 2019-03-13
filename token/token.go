package token

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	IDENT = "IDENT" //identifier (x, y)
	INT   = "INT"   //integer
	VOID = "VOID"

	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"

	LT = "<"
	GT = ">"

	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	RETURN_TYPE = ":"

	FUNCTION = "FUNCTION"
	LET  = "LET"
	RETURN   = "RETURN"
	FOR = "FOR"
	IN = "IN"

	CLASS = "CLASS"

	STRING = "STRING"

	TRUE  = "TRUE"
	FALSE = "FALSE"
	IF    = "IF"
	ELSE  = "ELSE"

	LBRACKET = "["
	RBRACKET = "]"

	EQ     = "=="
	NOT_EQ = "!="
)

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

var keywords = map[string]TokenType{
	"func":   FUNCTION,
	"int":    LET,
	"bool":   LET,
	"string": LET,
	"array":  LET,
	"return": RETURN,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"for":    FOR,
	"in":     IN,
	"void":   VOID,
	"class":  CLASS,
}

// LookupIdent checks the keywords table to see if identifier is a keyword
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
