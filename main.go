package main

import (
	"azula/repl"
	"fmt"
	"os"
)

func main() {
	fmt.Printf("\033[0;92mAzula V0.2\n")
	repl.Start(os.Stdin, os.Stdout)
}

func printParserErrors(errors []string) {
	fmt.Print("parser errors:\n")
	for _, msg := range errors {
		fmt.Print("\t" + msg + "\n")
	}
}
