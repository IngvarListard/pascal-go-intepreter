package calc4

import (
	"bytes"
	"fmt"
	"strconv"
	"unicode"
)

type Token int

const (
	EOF Token = iota
	Integer
	Plus
	Minus
	Mul
	Div

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

func (p *Parser) getNextLexeme() (Lexeme, error) {
	switch r := p.currentRune; {
	case unicode.IsDigit(r):
		return p.readInt(), nil
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
	v, ok := i.currentLexeme.Value().(int)
	if !ok {
		return 0, fmt.Errorf("factor: got unexpected type: expected int got %T", i.currentLexeme.Value())
	}
	var err error
	i.currentLexeme, err = i.parser.getNextLexeme()
	return v, err
}

type Lexeme interface {
	Token() Token
	Value() interface{}
}

type Const struct {
	token Token
	value interface{}
}

func (c *Const) Token() Token {
	return c.token
}

func (c *Const) Value() interface{} {
	return c.value
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
