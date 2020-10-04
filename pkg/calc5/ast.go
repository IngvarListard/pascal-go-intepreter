package calc5

import (
	"fmt"
)

type TokenTyp int

type Token struct {
	typ   TokenTyp
	value interface{}
}

type Node interface {
	Token() *Token
	Value() (interface{}, error)
}

type BinOp struct {
	left  Node
	right Node
	op    *Token
}

func (b *BinOp) Token() *Token { return b.op }

func (b *BinOp) Value() (interface{}, error) {
	lVal, _ := b.left.Value()
	lInt := lVal.(int)
	rVal, _ := b.right.Value()
	rInt := rVal.(int)
	switch b.op.typ {
	case Mul:
		return lInt * rInt, nil
	case Div:
		return lInt / rInt, nil
	case Plus:
		return lInt + rInt, nil
	case Minus:
		return lInt - rInt, nil
	default:
		return nil, fmt.Errorf("unexpected token type %v", b.op.typ)
	}
}

type Num struct {
	token *Token
	value interface{}
}

func (n *Num) Token() *Token { return n.token }

func (n *Num) Value() (interface{}, error) { return n.value, nil }

type UnaryOp struct {
	expr Node
	op   *Token
}

func (u *UnaryOp) Token() *Token { return u.op }

func (u *UnaryOp) Value() (interface{}, error) {
	return u.expr.Value()
}
