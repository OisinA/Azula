package repl

import (
	"azula/compiler"
	"azula/evaluator"
	"azula/lexer"
	"azula/object"
	"azula/parser"
	"bufio"
	"fmt"
	"io"
)

const PROMPT = "\033[0;36m>>\033[0;0m "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	vars := make(map[string]compiler.Type)
	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		if line == "exit" {
			fmt.Println("Goodbye :)")
			break
		}
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		t := compiler.NewTypecheckerFromVars(vars)
		_, err := t.Typecheck(program)
		if err != nil {
			fmt.Fprintf(out, "\033[0;31mType Error:\033[0;0m %s\n", err)
			continue
		}

		comp := compiler.New()
		_, err = comp.Compile(program)
		if err != nil {
			fmt.Fprintf(out, "%s", err)
		}
		fmt.Println(comp.Module)
	}
}

func StartInterpreted(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()
	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()
		if line == "exit" {
			fmt.Println("Goodbye :)")
			break
		}
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program, env)

		evaluated, ok := evaluated.(*object.Error)
		if ok {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, "parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
