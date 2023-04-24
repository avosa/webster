package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/avosa/webster/evaluator"
	"github.com/avosa/webster/lexer"
	"github.com/avosa/webster/object"
	"github.com/avosa/webster/parser"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: web <filename>")
		os.Exit(1)
	}

	filename := os.Args[1]
	code, err := ioutil.ReadFile(filename)

	if err != nil {
		fmt.Printf("Error reading file: %s\n", err)
		os.Exit(1)
	}

	l := lexer.New(string(code))
	p := parser.New(l)

	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		printParserErrors(p.Errors())
		os.Exit(1)
	}

	env := object.NewEnvironment()
	evaluated := evaluator.Eval(program, env)

	if evaluated != nil {
		fmt.Println(evaluated.Inspect())
	}
}

func printParserErrors(errors []string) {
	for _, errorMsg := range errors {
		fmt.Printf("\t%s\n", errorMsg)
	}
}
