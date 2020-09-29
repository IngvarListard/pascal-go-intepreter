package calc5

import "fmt"

type Parser struct {
	lexer        *Lexer
	currentToken *Token
}

func (p *Parser) term() Node {
	node := p.expr()

	for p.currentToken.typ == Plus || p.currentToken.typ == Minus {
		token := p.currentToken
		p.consume(p.currentToken.typ)
		node = &BinOp{left: node, right: p.expr(), op: token}
	}
	return node
}

func (p *Parser) expr() Node {
	node := p.factor()

	for p.currentToken.typ == Mul || p.currentToken.typ == Div {
		token := p.currentToken
		p.consume(p.currentToken.typ)
		node = &BinOp{left: node, right: p.factor(), op: token}
	}

	return node
}

func (p *Parser) factor() Node {
	token := p.currentToken
	if token.typ == Number {
		p.consume(Number)
		return &Num{token: token, value: token.value}
	} else if token.typ == Lparen {
		p.consume(Lparen)
		node := p.expr()
		p.consume(Rparen)
		return node
	}
	return nil
}

func (p *Parser) consume(typ TokenTyp) {
	if p.currentToken.typ == typ {
		p.currentToken = p.lexer.getNextToken()
		return
	}
	panic(fmt.Sprintf("Got unexpected op type: %v", typ))
}

func (p *Parser) parse() Node {
	return p.expr()
}
