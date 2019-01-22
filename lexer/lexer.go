package lexer

import (
	"github.com/muiscript/ether/token"
	"strings"
)

type Lexer struct {
	input           string
	currentPosition int
	peekPosition    int
	currentLine     int
	ch              byte
}

func New(input string) *Lexer {
	lexer := &Lexer{input: input, currentLine: 1}
	lexer.sanitizeInput()

	return lexer
}

func (l *Lexer) NextToken() token.Token {
	l.skipSpaces()

	var tok token.Token
	switch l.ch {
	// TODO: lex one char tokens
	case 0:
		tok = token.Token{Type: token.EOF, Literal: "", Line: l.currentLine}
	default:
		tok = token.Token{Type: token.ILLEGAL, Literal: string(l.ch), Line: l.currentLine}
	}

	l.consumeChar()
	return tok
}

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
		if l.ch == '\n' {
			l.currentLine++
		}
	}
	l.currentPosition = l.peekPosition
	l.peekPosition++
}

func (l *Lexer) skipSpaces() {
	for l.ch == ' ' || l.ch == '\n' {
		l.consumeChar()
	}
}
