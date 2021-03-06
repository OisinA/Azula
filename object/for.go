package object

import (
	"github.com/OisinA/Azula/ast"
	"bytes"
)

type For struct {
	Parameter *ast.Identifier
	Iterator *ast.Expression
	Body *ast.BlockStatement
	Env *Environment
}

func (f *For) Type() ObjectType {
	return FOR_OBJ
}

func (f *For) Inspect() string {
	var out bytes.Buffer

	out.WriteString("for(")
	out.WriteString(f.Parameter.String())
	out.WriteString(" in ")
	//out.WriteString(f.Iterator.String())
	out.WriteString(") {")
	out.WriteString(f.Body.String())
	out.WriteString(" }")

	return out.String()
}
