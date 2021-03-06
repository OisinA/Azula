package object

import (
	"github.com/OisinA/Azula/ast"
	"bytes"
	"strings"
)

type Function struct {
	Name *ast.Identifier
	Parameters []*ast.TypedIdentifier
	Body *ast.BlockStatement
	ReturnType *ast.Identifier
	Env *Environment
}

func (f *Function) Type() ObjectType {
	return FUNCTION_OBJ
}

func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("func ")
	out.WriteString(f.Name.String())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n{")

	return out.String()
}
