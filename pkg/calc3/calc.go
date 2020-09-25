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
	Plus
	Minus
)

type Interpreter struct {
	text         []rune
	pos          int
	currentRune  rune
	currentToken Token
}

func (i *Interpreter) next() {
	i.pos++
	if i.pos >= len(i.text) {
		i.currentToken = EOF
		i.currentRune = -1
		return
	}
	i.currentRune = i.text[i.pos]
}

func (i *Interpreter) skipWhitespace() {
	for i.currentRune == ' ' || i.currentRune == '\n' || i.currentRune == '\t' {
		i.next()
	}
}

func (i *Interpreter) getNextLexeme() (Lexeme, error) {
	switch v := i.currentRune; {
	case unicode.IsDigit(v):
		l := i.readInt()
		i.skipWhitespace()
		return l, nil
	case v == '+':
		l := &term{
			token: Plus,
			value: '+',
		}
		i.next()
		i.skipWhitespace()
		return l, nil
	case v == '-':
		l := &term{
			token: Minus,
			value: '-',
		}
		i.next()
		i.skipWhitespace()
		return l, nil
	default:
		return nil, fmt.Errorf("unexpected symbol occurance %v", i.currentRune)
	}
}

func (i *Interpreter) Expr() (int, error) {
	i.next()
	l, err := i.getNextLexeme()
	if err != nil {
		return 0, err
	}
	lv, ok := l.Value().(int)
	if !ok {
		return 0, fmt.Errorf("left value should be integer, got %T instead", l.Value())
	}

	result := lv

	for i.currentRune == '+' || i.currentRune == '-' {
		op, err := i.getNextLexeme()
		if err != nil {
			return 0, err
		}

		r, err := i.getNextLexeme()
		if err != nil {
			return 0, err
		}

		rv, ok := r.Value().(int)
		if !ok {
			return 0, fmt.Errorf("right value should be integer, go %T instead", r.Value())
		}

		switch op.Type() {
		case Plus:
			result += rv
		case Minus:
			result -= rv
		}
	}

	if i.currentToken != EOF {
		return 0, fmt.Errorf("unexpected symbol occrance, %v expected EOF", i.currentRune)
	}
	return result, nil
}

func NewInterpreter(text string) *Interpreter {
	return &Interpreter{
		text: []rune(text),
		pos:  -1,
	}
}

func (i *Interpreter) readInt() Lexeme {
	var numberBuf bytes.Buffer
	for unicode.IsDigit(i.currentRune) {
		numberBuf.WriteRune(i.currentRune)
		i.next()
	}
	number, _ := strconv.Atoi(numberBuf.String())
	return &term{token: Integer, value: number}
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
