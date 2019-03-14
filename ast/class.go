package ast

import (
	"azula/token"
	"bytes"
	"strings"
)

type ClassLiteral struct {
	Token token.Token
	Name *Identifier
	Parameters []*TypedIdentifier
	Body *BlockStatement
}

func (cl *ClassLiteral) expressionNode() {}

func (cl *ClassLiteral) TokenLiteral() string {
	return cl.Token.Literal
}

func (cl *ClassLiteral) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range cl.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(cl.TokenLiteral())
	out.WriteString(cl.Name.Value)
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {")
	out.WriteString(cl.Body.String())
	out.WriteString("}")

	return out.String()
}
