package builtins

import (
	"azula/object"
	"fmt"
)

func FunctionPrint(args ...object.Object) object.Object {
	if len(args) != 1 {
		return &object.Error{Message: "wrong number of arguments"}
	}
	fmt.Println(args[0].Inspect())
	return &object.Null{}
}