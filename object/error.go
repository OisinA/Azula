package object

import (
	"fmt"
)

type Error struct {
	Message    string
	LineNumber int
	ColNumber  int
}

func (e *Error) Type() Type {
	return ErrorObj
}

func (e *Error) Inspect() string {
	return "ERROR: " + e.Message + "(" + fmt.Sprint(e.LineNumber) + ", " + fmt.Sprint(e.ColNumber) + ")"
}
