package parser

import "fmt"

type ParserError struct {
	error
	line int
	msg  string
}

func (pe *ParserError) Error() string {
	return fmt.Sprintf("line %d: %s", pe.line, pe.msg)
}
