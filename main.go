package main

import (
	"azula/repl"
	"fmt"
	"os"
)

func main() {
	fmt.Printf("Azula V0.0\n")
	repl.Start(os.Stdin, os.Stdout)
}
