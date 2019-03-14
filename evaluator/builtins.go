package evaluator

import (
	"github.com/OisinA/Azula/object"
	"fmt"
	"bufio"
	"os"
	"strings"
	"strconv"
)

var builtins = map[string]*object.Builtin {
	"len": &object.Builtin {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%q, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			default:
				return newError("argument to 'len' not supported, got %s", args[0].Type())
			}
		},
	},
	"input": &object.Builtin {
		Fn: func(args ...object.Object) object.Object {
			if len(args) > 1 {
				return newError("wrong number of arguments. got=%q, want <= 1", len(args))
			}
			if len(args) == 1 {
				fmt.Print(args[0].Inspect())
			}

			reader := bufio.NewReader(os.Stdin)
			text, err := reader.ReadString('\n')
			if err != nil {
				return newError("error reading in input")
			}
			return &object.String{Value: strings.TrimSpace(string(text))}
		},
	},
	"toInt": &object.Builtin {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%q, want 1", len(args))
			}
			i, err := strconv.Atoi(args[0].Inspect())
			if err != nil {
				return newError("couldn't convert '%s' to int", args[0].Inspect())
			}
			return &object.Integer{Value: int64(i)}
		},
	},
	"print": &object.Builtin {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%q, want 1", len(args))
			}
			fmt.Println(args[0].Inspect())
			return NULL
		},
	},
	"range": &object.Builtin {
		Fn: func(args ...object.Object) object.Object {
			lower := int64(0)
			higher := int64(0)
			if len(args) == 1 {
				hi, ok := args[0].(*object.Integer)
				if !ok {
					return newError("can't get range of non-int " + args[0].Inspect())
				}
				higher = hi.Value
			} else if len(args) == 2 {
				low, ok := args[0].(*object.Integer)
				hi, ok := args[1].(*object.Integer)
				if !ok {
					return newError("can't get range of a non-int")
				}
				lower = low.Value
				higher = hi.Value
			} else {
				return newError("wrong number of arguments. got=%q, want 1/2", len(args))
			}
			array := &object.Array{ElementType: "int", Elements: []object.Object{}}
			for i := lower; i < higher; i++ {
				array.Elements = append(array.Elements, &object.Integer{Value: int64(i)})
			}
			return array
		},
	},
	"string_to_list": &object.Builtin {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%q", len(args))
			}

			s, ok := args[0].(*object.String)
			if !ok {
				return newError("cannot convert %v to string", args[0])
			}
			array := &object.Array{ElementType: "string", Elements: []object.Object{}}
			for _, c := range s.Value {
				array.Elements = append(array.Elements, &object.String{Value: string(c)})
			}
			return array
		},
	},
}
