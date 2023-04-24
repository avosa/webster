package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/avosa/webster/lexer"
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

	if len(p.Errors()) != 0 {
		printParserErrors(p.Errors())
		os.Exit(1)
	}

}

func printParserErrors(errors []string) {
	for _, errorMsg := range errors {
		fmt.Printf("\t%s\n", errorMsg)
	}
}
