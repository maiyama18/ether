package lexer

import (
	"github.com/muiscript/ether/token"
	"go/constant"
	"strings"
)

type Lexer struct {
	input           string
	currentPosition int
	peekPosition    int
	currentLine     int
	ch              byte
}

func NewLexer(input string) *Lexer {
	lexer := &Lexer{input: input}
	lexer.sanitizeInput()

	return lexer
}

func (l *Lexer) NextToken() token.Token {
	l.skipSpaces()
}

// replace \t and \r to white space.
// after executing this function, lexer should have
// only white space(' ') and new line('\n') as space characters.
func (l *Lexer) sanitizeInput() {
	sanitized := strings.Replace(l.input, "\t", " ", -1)
	sanitized = strings.Replace(sanitized, "\r", " ", -1)
	l.input = sanitized
}

func (l *Lexer) consumeChar() {
	if l.peekPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.peekPosition]
	}
	l.currentPosition = l.peekPosition
	l.peekPosition++
}

func (l *Lexer) skipSpaces() {
	for l.ch == ' ' || l.ch == '\n' {
		switch l.ch {
		case ' ':

		case '\n':

		}
	}
	if l.ch == ' ' || l.ch == '\t' || l.ch == '\r' {
		l.consumeChar()
	}
}
