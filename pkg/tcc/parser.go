package tcc

import (
	"fmt"
)

type Kind int

const (
	Const Kind = iota
	Var
	Add
	Sub
	Lt
	Set
	If1
	If2
	While1
	Do1
	Empty
	Seq
	Expr
	Prog
)

type Node struct {
	kind  Kind
	value int
	op1   *Node
	op2   *Node
	op3   interface{}
}

type Parser struct {
	lexer Lexer
}

func (p *Parser) term() (*Node, error) {
	switch p.lexer.tok {
	case Id:
		n := &Node{kind: Var, value: p.lexer.value}
		if err := p.lexer.nextToken(); err != nil {
			return nil, err
		}
		return n, nil
	case Num:
		n := &Node{kind: Const, value: p.lexer.value}
		if err := p.lexer.nextToken(); err != nil {
			return nil, err
		}
		return n, nil
	default:
		return p.parenExpr()
	}
}

func (p *Parser) sum() (*Node, error) {
	n, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.lexer.tok == Plus || p.lexer.tok == Minus {
		var kind Kind
		if p.lexer.tok == Plus {
			kind = Add
		} else {
			kind = Sub
		}
		op2, err := p.term()
		if err != nil {
			return nil, err
		}
		n = &Node{kind: kind, op1: n, op2: op2}
	}
	return n, nil
}

func (p *Parser) test() (*Node, error) {
	n, err := p.sum()
	if err != nil {
		return nil, err
	}

	if p.lexer.tok == Less {
		err := p.lexer.nextToken()
		if err != nil {
			return nil, err
		}

		op2, err := p.sum()
		if err != nil {
			return nil, err
		}

		n = &Node{kind: Lt, op1: n, op2: op2}
	}
	return n, nil
}

func (p *Parser) expr() (*Node, error) {
	if p.lexer.tok != Id {
		return p.test()
	}

	n, err := p.test()
	if err != nil {
		return nil, err
	}

	if n.kind == Var || p.lexer.tok == Equal {
		err := p.lexer.nextToken()
		if err != nil {
			return nil, err
		}

		op2, err := p.expr()
		if err != nil {
			return nil, err
		}
		n = &Node{kind: Set, op1: n, op2: op2}
	}
	return n, nil
}

func (p *Parser) parenExpr() (*Node, error) {
	if p.lexer.tok != Lpar {
		return nil, fmt.Errorf("expected '(' got %v", p.lexer.value)
	}

	if err := p.lexer.nextToken(); err != nil {
		return nil, err
	}

	n, err := p.expr()
	if err != nil {
		return nil, err
	}

	if p.lexer.tok != Rpar {
		return nil, fmt.Errorf("expected ')' got %v", p.lexer.value)
	}

	if err := p.lexer.nextToken(); err != nil {
		return nil, err
	}
	return n, nil
}

func (p *Parser) statement() (n *Node, err error) {
	switch p.lexer.tok {
	case If:
		n = &Node{kind: If1}
		if err = p.lexer.nextToken(); err != nil {
			return nil, err
		}
		if n.op1, err = p.parenExpr(); err != nil {
			return nil, err
		}
		if n.op2, err = p.statement(); err != nil {
			return nil, err
		}

		if p.lexer.tok == Else {
			n.kind = If2
			if err = p.lexer.nextToken(); err != nil {
				return nil, err
			}
			if n.op3, err = p.statement(); err != nil {
				return nil, err
			}
		}
	case While:
		n = &Node{kind: While1}
		if err = p.lexer.nextToken(); err != nil {
			return nil, err
		}
		if n.op1, err = p.parenExpr(); err != nil {
			return nil, err
		}
		if n.op2, err = p.statement(); err != nil {
			return nil, err
		}
	case Do:
		n = &Node{kind: Do1}
		if err = p.lexer.nextToken(); err != nil {
			return nil, err
		}
		if n.op1, err = p.statement(); err != nil {
			return nil, err
		}
		if p.lexer.tok != While {
			return nil, fmt.Errorf("'while' expected")
		}
		if err = p.lexer.nextToken(); err != nil {
			return nil, err
		}
		if n.op2, err = p.parenExpr(); err != nil {
			return n, err
		}
		if p.lexer.tok != Semicolon {
			return nil, fmt.Errorf("expected semicolon, got %v", p.lexer.value)
		}
	case Semicolon:
		n = &Node{kind: Empty}
		if err = p.lexer.nextToken(); err != nil {
			return nil, err
		}
	case Lbra:
		n = &Node{kind: Empty}
		if err = p.lexer.nextToken(); err != nil {
			return nil, err
		}

		for p.lexer.tok != Rbra {
			st, err := p.statement()
			if err != nil {
				return nil, err
			}
			n = &Node{kind: Seq, op1: n, op2: st}
			if err = p.lexer.nextToken(); err != nil {
				return nil, err
			}
		}
	default:
		e, err := p.expr()
		if err != nil {
			return nil, err
		}
		n = &Node{kind: Expr, op1: e}
		if p.lexer.tok != Semicolon {
			return nil, fmt.Errorf("expected ';' got %v", p.lexer.value)
		}
		if err = p.lexer.nextToken(); err != nil {
			return nil, err
		}
	}
	return n, nil
}

func (p *Parser) Parse() (*Node, error) {
	if err := p.lexer.nextToken(); err != nil {
		return nil, err
	}

	st, err := p.statement()
	if err != nil {
		return nil, err
	}
	node := &Node{kind: Prog, op1: st}
	if p.lexer.tok != EOF {
		panic("invalid statement syntax")
	}
	return node, nil
}
