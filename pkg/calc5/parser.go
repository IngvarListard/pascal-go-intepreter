package calc5

import "fmt"

type Parser struct {
	lexer        *Lexer
	currentToken *Token
}

func (p *Parser) term() Node {
	node := p.factor()

Loop:
	for {
		token := p.currentToken
		switch t := p.currentToken.typ; t {
		case Mul:
			p.consume(t)
		case IntegerDiv:
			p.consume(t)
		case FloatDiv:
			p.consume(t)
		default:
			break Loop
		}
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
	case IntegerConst:
		p.consume(IntegerConst)
		return &Num{token: token, value: token.value}
	case RealConst:
		p.consume(RealConst)
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
	p.consume(Program)
	varNode := p.variable()
	progName := varNode.Token().value
	p.consume(Semi)

	blockNode := p.block()
	programNode := &program{
		name:  progName.(string),
		block: blockNode.(*block),
	}
	p.consume(Dot)

	return programNode
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
			decs = append(decs, varDecl...)
			p.consume(Semi)
		}
	}

	for p.currentToken.typ == Procedure {
		p.consume(Procedure)
		procName := p.currentToken.value
		p.consume(Id)
		p.consume(Semi)
		blockNode := p.block()
		procDecl := &procDecl{
			procName: procName.(string),
			block:    blockNode.(*block),
		}
		decs = append(decs, procDecl)
		p.consume(Semi)
	}
	return decs
}

func (p *Parser) variableDeclaration() []Node {
	varNodes := []Node{&Var{
		token: p.currentToken,
		value: p.currentToken.value,
	}}
	p.consume(Id)

	for p.currentToken.typ == Comma {
		p.consume(Comma)
		varNodes = append(varNodes, &Var{
			token: p.currentToken,
			value: p.currentToken.value,
		})
		p.consume(Id)
	}
	p.consume(Colon)
	typNode := p.typeSpec()

	varDecls := make([]Node, len(varNodes))
	for i, node := range varNodes {
		varDecls[i] = &varDecl{
			varNode:  node,
			typeNode: typNode,
		}
	}
	return varDecls
}

func (p *Parser) typeSpec() Node {
	token := p.currentToken

	switch typ := p.currentToken.typ; typ {
	case Integer:
		p.consume(typ)
	case Real:
		p.consume(typ)
	default:
		panic(fmt.Sprintf("unexpected token %v", typ))
	}

	return &typeNode{
		token: token,
		value: token.value,
	}
}
