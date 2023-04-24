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
		fmt.Println("Usage: web <filename.wb>")
		os.Exit(1)
	}

	filename := os.Args[1]
	if !isValidFilename(filename) {
		fmt.Println("Error: Invalid file extension. Please provide a .wb file.")
		os.Exit(1)
	}

	fileContent, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}

	l := lexer.NewLexer(string(fileContent))
	p := parser.NewParser(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		printParserErrors(p.Errors())
		os.Exit(1)
	}

	env := object.NewEnvironment()
	evaluated := evaluator.Eval(program, env)
	if evaluated != nil {
		fmt.Println(evaluated.Inspect())
	}
}

func isValidFilename(filename string) bool {
	return len(filename) > 3 && filename[len(filename)-3:] == ".wb"
}

func printParserErrors(errors []string) {
	fmt.Println("Parser errors:")
	for _, err := range errors {
		fmt.Println("\t", err)
	}
}
