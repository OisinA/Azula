package builtins

import (
	"azula/object"
)

func FunctionLength(args ...object.Object) object.Object {
	if len(args) != 1 {
		return &object.Error{Message: "wrong number of arguments"}
	}

	switch arg := args[0].(type) {
	case *object.String:
		return &object.Integer{Value: int64(len(arg.Value))}
	case *object.Array:
		return &object.Integer{Value: int64(len(arg.Elements))}
	default:
		return &object.Error{Message: "argument to len not supported"}
	}
}