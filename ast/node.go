package ast

import "github.com/muiscript/ether/token"

type Node interface {
	Token() token.Token
}
