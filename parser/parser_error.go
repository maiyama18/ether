package parser

import "fmt"

type ParserError struct {
	error
	line int
	msg  string
}

func NewParserError(line int, msg string) *ParserError {
	return &ParserError{line: line, msg: msg}
}

func (pe *ParserError) Error() string {
	return fmt.Sprintf("line %d: %s", pe.line, pe.msg)
}
