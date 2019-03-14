package object

import (
	"azula/ast"
	"bytes"
	"strings"
)

type Class struct {
	Name *ast.Identifier
	Parameters []*ast.TypedIdentifier
	Body *ast.BlockStatement
	Env *Environment
}

func (c *Class) Type() ObjectType {
	return CLASS_OBJ
}

func (c *Class) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range c.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("class ")
	out.WriteString(c.Name.String())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(c.Body.String())
	out.WriteString("\n}")

	return out.String()
}
