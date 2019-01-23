package parser

import "fmt"

type ParserError struct {
	error
	line int
	msg  string
}

func NewParserError(line int, err error) *ParserError {
	switch err := err.(type) {
	case *ParserError:
		return err
	default:
		return &ParserError{line: line}
	}
}

func (pe *ParserError) Error() string {
	return fmt.Sprintf("line %d: %s", pe.line, pe.msg)
}
