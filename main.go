package main

import (
	"fmt"
	"github.com/muiscript/ether/evaluator"
	"github.com/muiscript/ether/lexer"
	"github.com/muiscript/ether/object"
	"github.com/muiscript/ether/parser"
	"github.com/muiscript/ether/repl"
	"io/ioutil"
	"os"
)

const USAGE = `
usage: ether [FILE_PATH]
`

func main() {
	switch len(os.Args) {
	case 1:
		repl.Start()
	case 2:
		os.Exit(interpret(os.Args[1]))
	default:
		fmt.Fprintf(os.Stderr, USAGE)
		os.Exit(1)
	}
}

func interpret(filename string) int {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		return 1
	}

	l := lexer.New(string(bytes))
	p := parser.New(l)

	program, err := p.ParseProgram()
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		return 2
	}

	env := object.NewEnvironment()
	_, err = evaluator.Eval(program, env)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		return 3
	}

	return 0
}

