package calc5

import "fmt"

type Parser struct {
	lexer        *Lexer
	currentToken *Token
}

func (p *Parser) term() Node {
	node := p.factor()

	for p.currentToken.typ == Mul || p.currentToken.typ == Div {
		token := p.currentToken
		p.consume(p.currentToken.typ)
		node = &BinOp{left: node, right: p.factor(), op: token}
	}
	return node
}

func (p *Parser) expr() Node {
	node := p.term()

	for p.currentToken.typ == Plus || p.currentToken.typ == Minus {
		token := p.currentToken
		p.consume(p.currentToken.typ)
		node = &BinOp{left: node, right: p.term(), op: token}
	}

	return node
}

func (p *Parser) factor() Node {
	token := p.currentToken

	switch token.typ {
	case Number:
		p.consume(Number)
		return &Num{token: token, value: token.value}
	case Lparen:
		p.consume(Lparen)
		node := p.expr()
		p.consume(Rparen)
		return node
	case Plus:
		p.consume(Plus)
		return &UnaryOp{expr: p.factor(), op: token}
	case Minus:
		p.consume(Minus)
		return &UnaryOp{expr: p.factor(), op: token}
	default:
		return p.variable()
	}
}

func (p *Parser) consume(typ TokenTyp) {
	if p.currentToken.typ == typ {
		p.currentToken = p.lexer.getNextToken()
		return
	}
	panic(fmt.Sprintf("Got unexpected op type: %v", typ))
}

func (p *Parser) parse() Node {
	node := p.program()
	if p.currentToken.typ != EOF {
		panic("not eof")
	}
	return node
}

func (p *Parser) program() Node {
	node := p.compoundStatement()
	p.consume(Dot)
	return node
}

func (p *Parser) compoundStatement() Node {
	p.consume(Begin)
	nodes := p.statementList()
	p.consume(End)

	root := &Compound{children: make([]Node, len(nodes))}
	for i, node := range nodes {
		root.children[i] = node
	}
	return root
}

func (p *Parser) statementList() []Node {
	node := p.statement()
	results := []Node{node}
	for p.currentToken.typ == Semi {
		p.consume(Semi)
		results = append(results, p.statement())
	}

	if p.currentToken.typ == Id {
		panic("unexpected token Id")
	}

	return results
}

func (p *Parser) statement() Node {
	switch p.currentToken.typ {
	case Begin:
		return p.compoundStatement()
	case Id:
		return p.assignmentStatement()
	default:
		return p.empty()
	}
}

func (p *Parser) assignmentStatement() Node {
	left := p.variable()
	token := p.currentToken
	p.consume(Assign)

	right := p.expr()

	return &assign{left: left, right: right, op: token}
}

func (p *Parser) variable() Node {
	node := &Var{token: p.currentToken}
	p.consume(Id)
	return node
}

func (p *Parser) empty() Node {
	return &NoOp{}
}

func (p *Parser) block() Node {
	declarationNodes := p.declarations()
	compoundStatementNode := p.compoundStatement()
	return &block{
		declarations:      declarationNodes,
		compoundStatement: compoundStatementNode.(*Compound),
	}
}

func (p *Parser) declarations() []Node {
	var decs []Node
	if p.currentToken.typ == VarT {
		p.consume(VarT)
		for p.currentToken.typ == Id {
			varDecl := p.variableDeclaration()
			decs = append(decs, varDecl)
			p.consume(Semi)
		}
	}
	return decs
}

func (p *Parser) variableDeclaration() Node {
	varNodes := []Node{&Var{
		token: p.currentToken,
		value: p.currentToken.value,
	}}

	for p.currentToken.typ == Comma {
		p.consume(Comma)
		varNodes = append(varNodes, &Var{
			token: p.currentToken,
			value: p.currentToken.value,
		})
		p.consume(Id)
	}

	p.consume(Comma)
	typNode := p.typeSpec()

	varDecls := make([]*varDecl, len())
	for _, node := range varNodes {
		varDecl{
			varNode:  node,
			typeNode: typNode,
		}
	}
}
