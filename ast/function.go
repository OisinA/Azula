package ast

import (
	"github.com/OisinA/Azula/token"
	"bytes"
	"strings"
)

type FunctionLiteral struct {
	Token      token.Token
	Name       *Identifier
	Parameters []*TypedIdentifier
	Body       *BlockStatement
	ReturnType *Identifier
}

func (fl *FunctionLiteral) expressionNode() {}

func (fl *FunctionLiteral) TokenLiteral() string {
	return fl.Token.Literal
}

func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(fl.TokenLiteral())
	out.WriteString(fl.Name.Value)
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")
	out.WriteString(": " + fl.ReturnType.Value)
	out.WriteString(fl.Body.String())

	return out.String()
}
