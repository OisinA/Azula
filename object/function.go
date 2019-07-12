package object

import (
	"azula/ast"
	"azula/code"
	"bytes"
	"fmt"
	"strings"
)

type Function struct {
	Name       *ast.Identifier
	Parameters []*ast.TypedIdentifier
	Body       *ast.BlockStatement
	ReturnType *ast.Identifier
	Env        *Environment
}

func (f *Function) Type() Type {
	return FunctionObj
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

type CompiledFunction struct {
	Instructions  code.Instructions
	NumLocals     int
	NumParameters int
}

func (cf *CompiledFunction) Type() Type {
	return CompiledFunctionObj
}

func (cf *CompiledFunction) Inspect() string {
	return fmt.Sprintf("CompiledFunction(%p)", cf)
}
