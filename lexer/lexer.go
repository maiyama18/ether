package lexer

import (
	"github.com/muiscript/ether/token"
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
	lexer.consumeChar()

	return lexer
}

func (l *Lexer) NextToken() token.Token {
	l.skipSpaces()

	var tok token.Token
	switch l.ch {
	case '=':
		tok = token.Token{Type: token.ASSIGN, Literal: "=", Line: l.currentLine}
	case '+':
		tok = token.Token{Type: token.PLUS, Literal: "+", Line: l.currentLine}
	case '-':
		tok = token.Token{Type: token.MINUS, Literal: "-", Line: l.currentLine}
	case '*':
		tok = token.Token{Type: token.ASTER, Literal: "*", Line: l.currentLine}
	case '/':
		tok = token.Token{Type: token.SLASH, Literal: "/", Line: l.currentLine}
	case '(':
		tok = token.Token{Type: token.LPAREN, Literal: "(", Line: l.currentLine}
	case ')':
		tok = token.Token{Type: token.RPAREN, Literal: ")", Line: l.currentLine}
	case '{':
		tok = token.Token{Type: token.LBRACE, Literal: "{", Line: l.currentLine}
	case '}':
		tok = token.Token{Type: token.RBRACE, Literal: "}", Line: l.currentLine}
	case '|':
		tok = token.Token{Type: token.BAR, Literal: "|", Line: l.currentLine}
	case ',':
		tok = token.Token{Type: token.COMMA, Literal: ",", Line: l.currentLine}
	case ';':
		tok = token.Token{Type: token.SEMICOLON, Literal: ";", Line: l.currentLine}
	case 0:
		tok = token.Token{Type: token.EOF, Literal: "", Line: l.currentLine}
	default:
		if isLetter(l.ch) {
			literal := l.readName()
			tok = token.Token{Type: token.TypeByLiteral(literal), Literal: literal, Line: l.currentLine}
		} else if isDigit(l.ch) {
			literal := l.readInteger()
			tok = token.Token{Type: token.INTEGER, Literal: literal, Line: l.currentLine}
		} else {
			tok = token.Token{Type: token.ILLEGAL, Literal: string(l.ch), Line: l.currentLine}
		}
	}

	l.consumeChar()
	return tok
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

func (l *Lexer) peekChar() byte {
	if l.peekPosition >= len(l.input) {
		return 0
	}
	return l.input[l.peekPosition]
}

func (l *Lexer) readName() string {
	start := l.currentPosition
	for {
		if pC := l.peekChar(); !isLetter(pC) && !isDigit(pC) {
			break
		}
		l.consumeChar()
	}

	return l.input[start : l.currentPosition+1]
}

func (l *Lexer) readInteger() string {
	start := l.currentPosition
	for {
		if pc := l.peekChar(); !isDigit(pc) {
			break
		}
		l.consumeChar()
	}

	return l.input[start : l.currentPosition+1]
}

func (l *Lexer) skipSpaces() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\r' || l.ch == '\n' {
		l.consumeChar()
	}
}

func isDigit(c byte) bool {
	return '0' <= c && c <= '9'
}

func isLetter(c byte) bool {
	return ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z') || c == '_'
}
