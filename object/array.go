package object

import (
	"bytes"
	"strings"
)

type Array struct {
	ElementType string
	Elements []Object
}

func (ao *Array) Type() Type {
	return ArrayObj
}

func (ao *Array) Inspect() string {
	var out bytes.Buffer

	elements := []string{}
	for _, e := range ao.Elements {
		elements = append(elements, e.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}
