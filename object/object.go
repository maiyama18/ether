package object

import "strconv"

type Type string

const (
	INTEGER = "INTEGER"
)

type Object interface {
	String() string
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
