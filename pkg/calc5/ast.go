package calc5

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
	token *Token
}

func (b *BinOp) Token() *Token { return b.token }

func (b *BinOp) Value() (interface{}, error) {
	panic("not implemented")
}

type Num struct {
	token *Token
	value interface{}
}

func (n *Num) Token() *Token { return n.token }

func (n *Num) Value() (interface{}, error) { return n.value, nil }
