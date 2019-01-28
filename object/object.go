package object

import (
	"bytes"
	"github.com/muiscript/ether/ast"
	"strconv"
	"strings"
)

type Type string

const (
	INTEGER          = "INTEGER"
	BOOLEAN          = "BOOLEAN"
	ARRAY            = "ARRAY"
	FUNCTION         = "FUNCTION"
	RETURN_VALUE     = "RETURN_VALUE"
	BUILTIN_FUNCTION = "BUILTIN_FUNCTION"
	NULL             = "NULL"
)

type Object interface {
	Type() Type
	String() string
}

type Integer struct {
	Value int
}

func (i *Integer) String() string { return strconv.Itoa(i.Value) }
func (i *Integer) Type() Type     { return INTEGER }

type Boolean struct {
	Value bool
}

func (b *Boolean) String() string { return strconv.FormatBool(b.Value) }
func (b *Boolean) Type() Type     { return BOOLEAN }

type Array struct {
	Elements []Object
}

func (a *Array) String() string {
	var elemStrs []string
	for _, elem := range a.Elements {
		elemStrs = append(elemStrs, elem.String())
	}

	return "[" + strings.Join(elemStrs, ", ") + "]"
}
func (a *Array) Type() Type { return ARRAY }

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
func (f *Function) Type() Type { return FUNCTION }

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) String() string { return "Return<" + rv.Value.String() + ">" }
func (rv *ReturnValue) Type() Type     { return RETURN_VALUE }

type BuiltinFunction struct {
	Fn func(args ...Object) (Object, error)
}

func (bf *BuiltinFunction) String() string { return "Builtin" }
func (bf *BuiltinFunction) Type() Type     { return BUILTIN_FUNCTION }

type Null struct{}

func (n *Null) String() string { return "null" }
func (n *Null) Type() Type     { return NULL }
