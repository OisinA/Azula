package evaluator

import (
	"azula/object"
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
}
