package lexer

import (
	"testing"

	"github.com/OisinA/Azula/token"
)

func TestNextToken(t *testing.T) {
	input := `int five = 5;
	int ten = 10;

	func add(int x, int y): int {
		return x + y;
	}
	
	int result = add(five, ten);

	!-/*5;
	5 < 10 < 5;

	if 5 < 10 {
		return true;
	} else {
		return false;
	}

	5 != 10;
	5 == 5;

	"foobar";
	"yes";

	[1, 2];
	for(i in x) {
		print(i);
	}

	i = 5;

	class TestClass(int x) {
		func get_x(): int {
			return x;
		}
	}

	TestClass c = TestClass(1);
	c.get_x();

	import "path/string.azl";
	`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LET, "int"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.LET, "int"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.FUNCTION, "func"},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.LET, "int"},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.LET, "int"},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.RETURN_TYPE, ":"},
		{token.LET, "int"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.LET, "int"},
		{token.IDENT, "result"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.BANG, "!"},
		{token.MINUS, "-"},
		{token.SLASH, "/"},
		{token.ASTERISK, "*"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.LT, "<"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.IF, "if"},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.INT, "5"},
		{token.NOT_EQ, "!="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.INT, "5"},
		{token.EQ, "=="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.STRING, "foobar"},
		{token.SEMICOLON, ";"},
		{token.STRING, "yes"},
		{token.SEMICOLON, ";"},
		{token.LBRACKET, "["},
		{token.INT, "1"},
		{token.COMMA, ","},
		{token.INT, "2"},
		{token.RBRACKET, "]"},
		{token.SEMICOLON, ";"},
		{token.FOR, "for"},
		{token.LPAREN, "("},
		{token.IDENT, "i"},
		{token.IN, "in"},
		{token.IDENT, "x"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "print"},
		{token.LPAREN, "("},
		{token.IDENT, "i"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.IDENT, "i"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.CLASS, "class"},
		{token.IDENT, "TestClass"},
		{token.LPAREN, "("},
		{token.LET, "int"},
		{token.IDENT, "x"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.FUNCTION, "func"},
		{token.IDENT, "get_x"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.RETURN_TYPE, ":"},
		{token.LET, "int"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.IDENT, "x"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.RBRACE, "}"},
		{token.IDENT, "TestClass"},
		{token.IDENT, "c"},
		{token.ASSIGN, "="},
		{token.IDENT, "TestClass"},
		{token.LPAREN, "("},
		{token.INT, "1"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "c"},
		{token.ACCESS, "."},
		{token.IDENT, "get_x"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.IMPORT, "import"},
		{token.STRING, "path/string.azl"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}
