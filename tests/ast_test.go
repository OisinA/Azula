package tests

import (
	"azula/token"
	"azula/ast"
	"testing"
)

func TestString(t *testing.T) {
	program := &ast.Program{
		Statements: []ast.Statement{
			&ast.LetStatement{
				Token: token.Token{Type: token.LET, Literal: "string"},
				Name: &ast.TypedIdentifier{
					Token: token.Token{Type: token.IDENT, Literal: "myVar"},
					Value: "myVar",
				},
				Value: &ast.TypedIdentifier{
					Token: token.Token{Type: token.IDENT, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}

	if program.String() != "string myVar = anotherVar;" {
		t.Errorf("program.String() wrong, got=%q", program.String())
	}
}
