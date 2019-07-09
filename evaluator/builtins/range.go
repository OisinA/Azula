package builtins

import (
	"azula/object"
)

func FunctionRange(args ...object.Object) object.Object {
	lower := int64(0)
	higher := int64(0)
	if len(args) == 1 {
		hi, ok := args[0].(*object.Integer)
		if !ok {
			return &object.Error{Message: ("can't get range of non-int " + args[0].Inspect())}
		}
		higher = hi.Value
	} else if len(args) == 2 {
		low, ok := args[0].(*object.Integer)
		hi, ok := args[1].(*object.Integer)
		if !ok {
			return &object.Error{Message: ("can't get range of non-int " + args[1].Inspect())}
		}
		lower = low.Value
		higher = hi.Value
	} else {
		return &object.Error{Message: "wrong number of arguments"}
	}
	array := &object.Array{ElementType: "int", Elements: []object.Object{}}
	for i := lower; i < higher; i++ {
		array.Elements = append(array.Elements, &object.Integer{Value: int64(i)})
	}
	return array
}
