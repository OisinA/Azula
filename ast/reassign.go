package ast

import (
	"github.com/OisinA/Azula/token"
	"bytes"
)

type ReassignStatement struct {
	Token token.Token
	Name *Identifier
	Value Expression
}

func (rs *ReassignStatement) statementNode() {}

func (rs *ReassignStatement) TokenLiteral() string {
	return rs.Token.Literal
}

func (rs *ReassignStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " = ")
	out.WriteString(rs.Value.String() + ";")

	return out.String()
}
