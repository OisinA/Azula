package ast

import (
	"azula/token"
	"bytes"
)

type ForLiteral struct {
	Token token.Token
	Parameter *Identifier
	Iterator Expression
	Body *BlockStatement
}

func (fl *ForLiteral) expressionNode() {}

func (fl *ForLiteral) TokenLiteral() string {
	return fl.Token.Literal
}

func (fl *ForLiteral) String() string {
	var out bytes.Buffer

	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(fl.Parameter.String())
	out.WriteString("in")
	out.WriteString(fl.Iterator.String())
	out.WriteString(") {")
	out.WriteString(fl.Body.String())
	out.WriteString("}")

	return out.String()
}
