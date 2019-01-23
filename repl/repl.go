package repl

import (
	"bufio"
	"fmt"
	"github.com/muiscript/ether/lexer"
	"github.com/muiscript/ether/token"
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
		lex := lexer.New(input)

		for tok := lex.NextToken(); tok.Type != token.EOF; tok = lex.NextToken() {
			fmt.Printf("%+v\n", tok)
		}
	}
}
