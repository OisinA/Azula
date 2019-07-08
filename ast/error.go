package ast

import (
	"azula/token"
)

type ErrorLiteral struct {
	Token token.Token
	Value string
}

func (sl *ErrorLiteral) expressionNode() {}

func (sl *ErrorLiteral) TokenLiteral() string {
	return sl.Token.Literal
}

func (sl *ErrorLiteral) String() string {
	return sl.Token.Literal
}