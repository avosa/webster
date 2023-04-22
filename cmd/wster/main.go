package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"webster/pkg/webster"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: wster <filename>.wb")
		os.Exit(1)
	}

	filename := os.Args[1]
	source, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	interpreter := webster.NewInterpreter(string(source))
	result := interpreter.Expr()
	fmt.Printf("Result: %f\n", result)
}
