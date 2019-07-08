package builtins

import (
	"azula/object"
	"fmt"
)

var typeMap = map[object.Type]string{
	object.IntegerObj: "int",
	object.BooleanObj: "bool",
	object.StringObj:  "string",
	object.ArrayObj:   "array",
	object.ErrorObj:   "error",
}

func FunctionAppend(args ...object.Object) object.Object {
	if len(args) < 2 {
		return &object.Error{Message: "wrong number of arguments"}
	}
	iterator := args[0]
	switch i := iterator.(type) {
	case *object.Array:
		return &object.Array{ElementType: i.ElementType, Elements: append(i.Elements, args[1])}
	case *object.Hash:
		key := args[1]

		hashKey, ok := key.(object.Hashable)
		if !ok {
			return &object.Error{Message: "unusable as hash key"}
		}

		if i.KeyType != typeMap[key.Type()] {
			return &object.Error{Message: "cannot assign to incorrect map type"}
		}

		value := args[2]

		hashed := hashKey.HashKey()
		i.Pairs[hashed] = object.HashPair{Key: key, Value: value}
		return i
	}
	fmt.Println(args[0].Inspect())
	return &object.Null{}
}