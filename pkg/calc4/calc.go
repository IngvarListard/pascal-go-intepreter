package calc4

import (
	"bytes"
	"fmt"
	"strconv"
	"unicode"
)

type Token int

func (t Token) Priority() int {
	switch t {
	case Plus, Minus:
		return 3
	case Mul, Div:
		return 2
	}
	return -1
}

const (
	EOF Token = iota
	Integer
	Plus
	Minus
	Mul
	Div
	Lparen
	Rparen

	Null rune = -1
)

type Parser struct {
	text        []rune
	pos         int
	currentRune rune
}

func (p *Parser) next() {
	p.pos++
	if p.pos >= len(p.text) {
		p.currentRune = Null
		return
	}
	p.currentRune = p.text[p.pos]
}

func (p *Parser) readInt() Lexeme {
	var numBuf bytes.Buffer
	for unicode.IsDigit(p.currentRune) {
		numBuf.WriteRune(p.currentRune)
		p.next()
	}
	number, _ := strconv.Atoi(numBuf.String())
	p.next()
	p.skipWhitespace()
	return &Const{token: Integer, value: number}
}

type brackets struct {
	lexeme Lexeme
}

func (b *brackets) Token() Token {
	panic("implement me")
}

func (b *brackets) Value() (interface{}, error) {
	panic("implement me")
}

//type binExpr struct {
//	left Lexeme
//	right Lexeme
//	op Token
//	token Token
//}
//
//func (b *binExpr) Token() Token {
//	return b.token
//}
//
//func (b *binExpr) Value() (interface{}, error) {
//	l, err := b.left.Value()
//	if err != nil {
//		return nil, err
//	}
//
//	r, err := b.right.Value()
//	if err != nil {
//		return nil, err
//	}
//
//	lInt, ok := l.(int)
//	if !ok {
//		return nil, fmt.Errorf("int expected got %T", l)
//	}
//
//	rInt, ok := l.(int)
//	if !ok {
//		return nil, fmt.Errorf("int expected got %T", l)
//	}
//
//	switch b.op {
//	case Plus:
//		return lInt + rInt, nil
//	case Minus:
//		return lInt - rInt, nil
//	case Mul:
//		return lInt * rInt, nil
//	case Div:
//		return lInt / rInt, nil
//	default:
//		return nil, fmt.Errorf("unexpected operator %v", b.op)
//	}
//}
type bracket struct {
	token  Token
	lexeme Lexeme
}

func (b *bracket) Token() Token {
	return b.token
}

func (b *bracket) Value() (interface{}, error) {
	return b.lexeme.Value()
}

func (p *Parser) getNextLexeme() (Lexeme, error) {
	switch r := p.currentRune; {
	case unicode.IsDigit(r):
		return p.readInt(), nil
	case r == '(':
		p.next()
		p.skipWhitespace()
		return &Const{token: Lparen, value: r}, nil
	case r == ')':
		p.next()
		p.skipWhitespace()
		return &Const{token: Rparen, value: r}, nil
	case r == '+':
		p.next()
		p.skipWhitespace()
		return &Const{token: Plus, value: r}, nil
	case r == '-':
		p.next()
		p.skipWhitespace()
		return &Const{token: Minus, value: r}, nil
	case r == '*':
		p.next()
		p.skipWhitespace()
		return &Const{token: Mul, value: r}, nil
	case r == '/':
		p.next()
		p.skipWhitespace()
		return &Const{token: Div, value: r}, nil
	case r == Null:
		p.next()
		p.skipWhitespace()
		return &Const{token: EOF, value: Null}, nil
	default:
		return nil, fmt.Errorf("unexpected symbol occurance: %s", string(r))
	}
}

func (p *Parser) skipWhitespace() {
	for p.currentRune == ' ' || p.currentRune == '\n' || p.currentRune == '\t' {
		p.next()
	}
}

type Interpreter struct {
	parser        *Parser
	currentLexeme Lexeme
}

func (i *Interpreter) Term() (int, error) {
	result, err := i.Expr()
	if err != nil {
		return 0, err
	}
	for i.currentLexeme.Token() == Plus || i.currentLexeme.Token() == Minus {
		switch i.currentLexeme.Token() {
		case Plus:
			err = i.consume(Plus)
			right, err := i.Expr()
			if err != nil {
				return 0, err
			}
			result += right
		case Minus:
			err = i.consume(Minus)
			right, err := i.Expr()
			if err != nil {
				return 0, err
			}
			result -= right
		}
	}
	return result, nil
}

func (i *Interpreter) consume(t Token) error {
	var err error
	if i.currentLexeme.Token() == t {
		i.currentLexeme, err = i.parser.getNextLexeme()
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *Interpreter) Expr() (int, error) {
	result, err := i.Factor()
	if err != nil {
		return 0, err
	}
	for i.currentLexeme.Token() == Mul || i.currentLexeme.Token() == Div {
		switch i.currentLexeme.Token() {
		case Mul:
			err = i.consume(Mul)
			v, err := i.Factor()
			if err != nil {
				return 0, err
			}
			result *= v
		case Div:
			err = i.consume(Div)
			v, err := i.Factor()
			if err != nil {
				return 0, err
			}
			result /= v
		}
	}
	return result, err
}

func (i *Interpreter) Factor() (int, error) {

	switch i.currentLexeme.Token() {
	case Lparen:
		err := i.consume(Lparen)
		if err != nil {
			return 0, err
		}
		r, err := i.Term()
		if err != nil {
			return 0, err
		}
		return r, nil
	case Integer:
		v, err := i.currentLexeme.Value()
		if err != nil {
			return 0, err
		}
		vv, ok := v.(int)
		if !ok {
			return 0, fmt.Errorf("factor: got unexpected type: expected int got %T", v)
		}
		i.currentLexeme, err = i.parser.getNextLexeme()
		return vv, err
	default:
		return 0, fmt.Errorf("unexpected token %v", i.currentLexeme.Token())
	}

}

type Lexeme interface {
	Token() Token
	Value() (interface{}, error)
}

type Const struct {
	token Token
	value interface{}
}

func (c *Const) Token() Token {
	return c.token
}

func (c *Const) Value() (interface{}, error) {
	return c.value, nil
}

func NewParser(text string) (*Parser, error) {
	if len(text) == 0 {
		return nil, fmt.Errorf("text is empty")
	}

	runes := []rune(text)
	parser := &Parser{
		text:        runes,
		pos:         0,
		currentRune: runes[0],
	}
	parser.skipWhitespace()
	return parser, nil
}

func NewInterpreter(parser *Parser) (*Interpreter, error) {
	l, err := parser.getNextLexeme()
	if err != nil {
		return nil, err
	}
	return &Interpreter{
		parser:        parser,
		currentLexeme: l,
	}, nil
}
