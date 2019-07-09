package builtins

import (
	"azula/object"
	"bufio"
	"fmt"
	"os"
	"strings"
)

func FunctionInput(args ...object.Object) object.Object {
	if len(args) == 1 {
		fmt.Print(args[0].Inspect())
	}

	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	if err != nil {
		return &object.Error{Message: "could not read input"}
	}
	return &object.String{Value: strings.TrimSpace(string(text))}
}
