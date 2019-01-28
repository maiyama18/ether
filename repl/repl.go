package repl

import (
	"bufio"
	"fmt"
	"github.com/muiscript/ether/evaluator"
	"github.com/muiscript/ether/lexer"
	"github.com/muiscript/ether/object"
	"github.com/muiscript/ether/parser"
	"os"
)

const PROMPT = "~> "

func Start() {
	scanner := bufio.NewScanner(os.Stdin)

	env := object.NewEnvironment()

	for {
		fmt.Print(PROMPT)

		if ok := scanner.Scan(); !ok {
			continue
		}
		input := scanner.Text()

		l := lexer.New(input)
		p := parser.New(l)

		program, err := p.ParseProgram()
		if err != nil {
			fmt.Println(err)
			continue
		}

		evaluated, err := evaluator.Eval(program, env)
		if err != nil {
			fmt.Println(err)
			continue
		}

		if evaluated != nil {
			fmt.Println(evaluated)
		}
	}
}
