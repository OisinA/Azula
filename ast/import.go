package ast

import (
	"github.com/OisinA/Azula/token"
	"bytes"
)

type ImportStatement struct {
	Token token.Token
	Value Expression
}

func (is *ImportStatement) statementNode() {}

func (is *ImportStatement) TokenLiteral() string {
	return is.Token.Literal
}

func (is *ImportStatement) String() string {
	var out bytes.Buffer

	out.WriteString(is.TokenLiteral() + " ")
	out.WriteString(is.Value.String())
	out.WriteString(";")

	return out.String()
}
