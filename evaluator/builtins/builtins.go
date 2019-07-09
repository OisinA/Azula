package builtins

import (
	"azula/object"
)

var Builtins = map[string]*object.Builtin{
	"print":  &object.Builtin{FunctionPrint},
	"len":    &object.Builtin{FunctionLength},
	"append": &object.Builtin{FunctionAppend},
	"input":  &object.Builtin{FunctionInput},
	"range":  &object.Builtin{FunctionRange},
}
