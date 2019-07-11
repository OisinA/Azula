package main

import (
	"azula/compiler"
	"azula/lexer"
	"azula/parser"
	"azula/repl"
	"azula/vm"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func main() {
	if len(os.Args[1:]) > 0 {
		if os.Args[1] == "build" {
			out := os.Stdout
			dat, err := ioutil.ReadFile(os.Args[2])
			if err != nil {
				fmt.Println("error: couldn't find file " + os.Args[1])
				return
			}
			l := lexer.New(string(dat))
			p := parser.New(l)

			program := p.ParseProgram()

			if len(p.Errors()) != 0 {
				printParserErrors(p.Errors())
				return
			}

			t := compiler.NewTypechecker()
			_, err = t.Typecheck(program)
			if err != nil {
				fmt.Fprintf(out, "\033[0;31mType Error:\033[0;0m %s\n", err)
				return
			}

			comp := compiler.New()
			err = comp.Compile(program)
			if err != nil {
				fmt.Fprintf(out, "\033[0;31mCompile Error:\033[0;0m %s\n", err)
				return
			}

			err = ioutil.WriteFile(os.Args[2]+"_executable", comp.Bytecode().Instructions, 0644)
			if err != nil {
				fmt.Fprintf(out, "\033[0;31mFile Error:\033[0;0m %s\n", err)
			}
		} else if os.Args[1] == "run" {
			out := os.Stdout
			dat, err := ioutil.ReadFile(os.Args[2])
			if err != nil {
				fmt.Println("error: couldn't find file " + os.Args[1])
				return
			}
			l := lexer.New(string(dat))
			p := parser.New(l)

			program := p.ParseProgram()

			if len(p.Errors()) != 0 {
				printParserErrors(p.Errors())
				return
			}

			t := compiler.NewTypechecker()
			_, err = t.Typecheck(program)
			if err != nil {
				fmt.Fprintf(out, "\033[0;31mType Error:\033[0;0m %s\n", err)
				return
			}

			comp := compiler.New()
			err = comp.Compile(program)
			if err != nil {
				fmt.Fprintf(out, "\033[0;31mCompile Error:\033[0;0m %s\n", err)
				return
			}

			mac := vm.New(comp.Bytecode())
			err = mac.Run()
			if err != nil {
				fmt.Fprintf(out, "Executing failed :(\n%s\n", err)
				return
			}

			stackTop := mac.LastPoppedStackElem()
			io.WriteString(out, stackTop.Inspect())
			io.WriteString(out, "\n")
		} else if os.Args[1] == "interpret" {
			fmt.Printf("\033[0;92mAzula V0.1\n")
			repl.StartInterpreted(os.Stdin, os.Stdout)
		}
	} else {
		fmt.Printf("\033[0;92mAzula V0.2\n")
		repl.Start(os.Stdin, os.Stdout)
	}
}

func printParserErrors(errors []string) {
	fmt.Print("parser errors:\n")
	for _, msg := range errors {
		fmt.Print("\t" + msg + "\n")
	}
}
