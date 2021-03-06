package main

import (
	"github.com/OisinA/Azula/repl"
	"fmt"
	"os"
	"io/ioutil"
	"github.com/OisinA/Azula/lexer"
	"github.com/OisinA/Azula/parser"
	"github.com/OisinA/Azula/evaluator"
	"github.com/OisinA/Azula/object"
)

func main() {
	if len(os.Args[1:]) > 0 {
		dat, err := ioutil.ReadFile(os.Args[1])
		if err != nil {
			fmt.Println("error: couldn't find file " + os.Args[1])
			return
		}
		env := object.NewEnvironment()
		l := lexer.New(string(dat))
		p := parser.New(l)

		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			printParserErrors(p.Errors())
			return
		}

		evaluated := evaluator.Eval(program, env)
		evaluated, ok := evaluated.(*object.Error)

		if ok {
			fmt.Println(evaluated.Inspect())
		}
	} else {
		fmt.Printf("Azula V0.0\n")
		repl.Start(os.Stdin, os.Stdout)
	}
}

func printParserErrors(errors []string) {
	fmt.Print("parser errors:\n")
	for _, msg := range errors {
		fmt.Print("\t"+msg+"\n")
	}
}
