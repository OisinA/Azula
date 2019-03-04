package main

import (
	"fmt"
	"os"
	"azula/repl"
)

func main() {
	fmt.Printf("Azula V0.0\n")
	repl.Start(os.Stdin, os.Stdout)
}
