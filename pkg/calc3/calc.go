package calc3

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
	Mul
	Div
)

type Parser struct {
	text        []rune
	pos         int
	currentRune rune
}

func (p *Parser) next() {
	p.pos++
	if p.pos >= len(p.text) {
		p.currentRune = -1
		return
	}
	p.currentRune = p.text[p.pos]
}

func (p *Parser) skipWhitespace() {
	for p.currentRune == ' ' || p.currentRune == '\n' || p.currentRune == '\t' {
		p.next()
	}
}

func (p *Parser) getNextLexeme() (Lexeme, error) {
	switch v := p.currentRune; {
	case unicode.IsDigit(v):
		l := p.readInt()
		p.skipWhitespace()
		return l, nil
	case v == '*':
		l := &term{
			token: Mul,
			value: '*',
		}
		p.next()
		p.skipWhitespace()
		return l, nil
	case v == '/':
		l := &term{
			token: Div,
			value: '/',
		}
		p.next()
		p.skipWhitespace()
		return l, nil
	case v == -1:
		return nil, nil
	default:
		return nil, fmt.Errorf("unexpected symbol occurance %v", p.currentRune)
	}
}

func (p *Parser) readInt() Lexeme {
	var numberBuf bytes.Buffer
	for unicode.IsDigit(p.currentRune) {
		numberBuf.WriteRune(p.currentRune)
		p.next()
	}
	number, _ := strconv.Atoi(numberBuf.String())
	return &term{token: Integer, value: number}
}

type Interpreter struct {
	*Parser
	currentLexeme Lexeme
}

func (i *Interpreter) consume(token Token) error {
	var err error
	if token == i.currentLexeme.Type() {
		i.currentLexeme, err = i.getNextLexeme()
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *Interpreter) factor() (interface{}, error) {
	l := i.currentLexeme
	if err := i.consume(l.Type()); err != nil {
		return nil, err
	}
	return l.Value(), nil
}

func (i *Interpreter) Expr() (int, error) {
	v, err := i.factor()
	if err != nil {
		return 0, err
	}

	result, ok := v.(int)
	if !ok {
		return 0, fmt.Errorf("expected int type got %T instead", v)
	}

	for i.currentLexeme != nil && (i.currentLexeme.Type() == Mul || i.currentLexeme.Type() == Div) {
		switch i.currentLexeme.Type() {
		case Mul:
			if err := i.consume(Mul); err != nil {
				return 0, err
			}
			val, err := i.factor()
			if err != nil {
				return 0, err
			}

			v, ok := val.(int)
			if !ok {
				return 0, fmt.Errorf("expected in type got %T isntead", val)
			}

			result *= v
		case Div:
			if err := i.consume(Div); err != nil {
				return 0, err
			}
			val, err := i.factor()
			if err != nil {
				return 0, err
			}

			v, ok := val.(int)
			if !ok {
				return 0, fmt.Errorf("expected in type got %T isntead", val)
			}

			result /= v
		}
	}

	return result, nil
}

func NewParser(text string) (*Parser, error) {
	if len(text) == 0 {
		return nil, fmt.Errorf("empty input")
	}
	runes := []rune(text)
	return &Parser{
		text:        runes,
		pos:         0,
		currentRune: runes[0],
	}, nil
}

func NewInterpreter(p *Parser) (*Interpreter, error) {
	l, err := p.getNextLexeme()
	if err != nil {
		return nil, err
	}
	return &Interpreter{
		Parser:        p,
		currentLexeme: l,
	}, nil
}

type Lexeme interface {
	Type() Token
	Value() interface{}
}

type term struct {
	token Token
	value interface{}
}

func (t *term) Value() interface{} {
	return t.value
}

func (t *term) Type() Token {
	return t.token
}
