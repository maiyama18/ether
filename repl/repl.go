package repl

import (
	"bufio"
	"fmt"
	"github.com/muiscript/ether/lexer"
	"github.com/muiscript/ether/parser"
	"os"
)

const PROMPT = "~> "

func Start() {
	scanner := bufio.NewScanner(os.Stdin)

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

		fmt.Println(program.String())
	}
}
