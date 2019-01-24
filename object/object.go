package object

import (
	"bytes"
	"github.com/muiscript/ether/ast"
	"strconv"
	"strings"
)

type Type string

const (
	INTEGER  = "INTEGER"
	FUNCTION = "FUNCTION"
)

type Object interface {
	Type() Type
}

type Integer struct {
	Value int
}

func (i *Integer) String() string {
	return strconv.Itoa(i.Value)
}
func (i *Integer) Type() Type {
	return INTEGER
}

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) String() string {
	var paramStrs []string
	for _, param := range f.Parameters {
		paramStrs = append(paramStrs, param.String())
	}

	var out bytes.Buffer
	out.WriteString("|")
	out.WriteString(strings.Join(paramStrs, ", "))
	out.WriteString("| ")
	out.WriteString(f.Body.String())

	return out.String()
}
func (f *Function) Type() Type {
	return FUNCTION
}
