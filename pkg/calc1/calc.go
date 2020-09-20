package calc1

import (
	"errors"
	"fmt"
	"strconv"
	"unicode"
)

const (
	INTEGER tokenType = iota
	PLUS
	EOF
)

type tokenType int

type Token struct {
	typ   tokenType
	value int
}

func NewToken(typ tokenType, value int) *Token {
	return &Token{
		typ:   typ,
		value: value,
	}
}

func (t *Token) String() string {
	template := "Token{ typ: %v value: %v}"
	switch t.typ {
	case INTEGER:
		return fmt.Sprintf(template, "INTEGER", t.value)
	case PLUS:
		return fmt.Sprintf(template, "PLUS", t.value)
	case EOF:
		return fmt.Sprintf(template, "EOF", t.value)
	default:
		return fmt.Sprintf(template, "UNKNOWN", t.value)
	}
}

func NewInterpreter(text string) *Interpreter {
	return &Interpreter{
		text:         []rune(text),
		pos:          0,
		currentToken: nil,
	}
}

type Interpreter struct {
	text         []rune
	pos          int
	currentToken *Token
}

func (i *Interpreter) getNextToken() (*Token, error) {
	text := i.text

	if i.pos > len(text)-1 {
		return NewToken(EOF, 0), nil
	}

	currentChar := i.text[i.pos]

	if unicode.IsDigit(currentChar) {
		v, err := strconv.Atoi(string(currentChar))
		if err != nil {
			return nil, err
		}
		i.pos++

		return NewToken(INTEGER, v), nil
	}

	if currentChar == '+' {
		i.pos++
		return NewToken(PLUS, '+'), nil
	}

	return nil, errors.New("error parsing input")
}

func (i *Interpreter) eat(typ tokenType) error {
	if i.currentToken.typ == typ {
		t, err := i.getNextToken()
		if err != nil {
			return err
		}
		i.currentToken = t
		return nil
	}
	return errors.New("eat: got wrong type")
}

func (i *Interpreter) Eval() (int, error) {
	var err error
	i.currentToken, err = i.getNextToken()
	if err != nil {
		return 0, err
	}

	left := i.currentToken
	if err = i.eat(INTEGER); err != nil {
		return 0, err
	}

	_ = i.currentToken
	if err = i.eat(PLUS); err != nil {
		return 0, err
	}

	right := i.currentToken
	if err = i.eat(INTEGER); err != nil {
		return 0, err
	}

	return left.value + right.value, nil
}
