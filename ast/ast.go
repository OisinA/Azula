package ast

import (
	"bytes"

	"azula/token"
)

// Node is an AST Node
type Node interface {
	TokenLiteral() string
	String() string
}

// Statement is an AST Statement
type Statement interface {
	Node
	statementNode()
}

// Expression is an AST Expression
type Expression interface {
	Node
	expressionNode()
}

// Program is a collection of Statements
type Program struct {
	Statements []Statement
}

// TokenLiteral returns the literal of the first statement in a program
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode() {}

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

func (i *Identifier) String() string {
	return i.Value
}

type TypedIdentifier struct {
	Token      token.Token
	Value      string
	ReturnType Identifier
}

func (i *TypedIdentifier) expressionNode() {}

func (i *TypedIdentifier) TokenLiteral() string {
	return i.Token.Literal
}

func (i *TypedIdentifier) String() string {
	return i.Value
}
