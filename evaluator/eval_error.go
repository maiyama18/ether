package evaluator

import "fmt"

type EvalError struct {
	error
	line int
	msg  string
}

func (ee *EvalError) Error() string {
	return fmt.Sprintf("line %d: %s", ee.line, ee.msg)
}
