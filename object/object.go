package object

import "strconv"

type Object interface {
	String() string
}

type Integer struct {
	Value int
}

func (i *Integer) String() string {
	return strconv.Itoa(i.Value)
}
