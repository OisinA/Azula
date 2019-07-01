package token

const (
	// ILLEGAL - illegal token
	ILLEGAL = "ILLEGAL"
	// EOF - end of file
	EOF = "EOF"

	// IDENT - identifier (x, y, foo)
	IDENT = "IDENT"

	// DATATYPES

	// INT - integer literal
	INT = "INT"
	// STRING - string literal
	STRING = "STRING"
	// TRUE - true boolean value
	TRUE = "TRUE"
	// FALSE - false boolean value
	FALSE = "FALSE"

	// VOID - no return type
	VOID = "VOID"
	// CLASS - class type
	CLASS = "CLASS"

	// OPERATORS

	// ASSIGN - token to assign to identifier
	ASSIGN = "="
	// PLUS - token for addition or appending
	PLUS = "+"
	// MINUS - subtracting
	MINUS = "-"
	// BANG - inversion
	BANG = "!"
	// ASTERISK - multiplying
	ASTERISK = "*"
	// LT - less than
	LT = "<"
	// GT - greater than
	GT = ">"
	// SLASH - slash
	SLASH = "/"

	// CODE STRUCTURE

	// COMMA - separate terms
	COMMA = ","
	// SEMICOLON - end of statement
	SEMICOLON = ";"
	// LPAREN - left parenthesis
	LPAREN = "("
	// RPAREN - right parenthesis
	RPAREN = ")"
	// LBRACE - left brace
	LBRACE = "{"
	// RBRACE - right brace
	RBRACE = "}"
	// LBRACKET - left bracket
	LBRACKET = "["
	// RBRACKET - right bracket
	RBRACKET = "]"
	// COLON - used to indicate return type
	COLON = ":"
	// ACCESS - access identifier inside identifier
	ACCESS = "ACCESS"

	// FUNCTION - function definition
	FUNCTION = "FUNCTION"
	// LET - initialising variable
	LET = "LET"
	// RETURN - return statement
	RETURN = "RETURN"
	// FOR - for loop
	FOR = "FOR"
	// IN - looping through list and checking inclusion
	IN = "IN"
	// IMPORT - importing another file into current file
	IMPORT = "IMPORT"
	// IF - if statement
	IF = "IF"
	// ELSE - else statement
	ELSE = "ELSE"

	// EQ - equality
	EQ = "=="
	// NOTEQ - inequality
	NOTEQ = "!="
)

// Type is the type of the token
type Type string

// Token is each token in a program
type Token struct {
	Type       Type
	Literal    string
	LineNumber int
	ColNumber  int
}

var keywords = map[string]Type{
	"int":    LET,
	"bool":   LET,
	"string": LET,
	"array":  LET,

	"void": VOID,

	"true":  TRUE,
	"false": FALSE,

	"if":   IF,
	"else": ELSE,
	"for":  FOR,
	"in":   IN,

	"func":  FUNCTION,
	"class": CLASS,

	"import": IMPORT,
	"return": RETURN,
}

// LookupIdent checks the keywords table to see if identifier is a keyword
func LookupIdent(ident string) Type {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
