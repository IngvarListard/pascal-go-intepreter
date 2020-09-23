package calc21

import (
	"bytes"
	"fmt"
	"strconv"
	"unicode"
)

const (
	WrongChar = -1

	EOF tokenType = iota
	Integer
	Plus
	Minus
	Mul
	Div
)

var op = map[int]tokenType{
	'+': Plus,
	'-': Minus,
}

type tokenType int

type Token struct {
	typ tokenType
	val int
}

func (t *Token) Value() int {
	return t.val
}

type Interpreter struct {
	pos          int
	text         []rune
	currentToken *Token
	currentChar  rune
}

func NewInterpreter(text string) *Interpreter {
	textRunes := []rune(text)
	return &Interpreter{
		pos:          0,
		text:         textRunes,
		currentToken: nil,
		currentChar:  textRunes[0],
	}
}

func (i *Interpreter) next() {
	i.pos++
	if i.pos > len(i.text)-1 {
		i.currentChar = WrongChar
	} else {
		i.currentChar = i.text[i.pos]
	}
}

func (i *Interpreter) skipWhitespaces() {
	for unicode.IsSpace(i.currentChar) && i.currentChar != WrongChar {
		i.next()
	}
}

func (i *Interpreter) readInteger() *Token {
	var number bytes.Buffer

	for {
		if i.currentChar != WrongChar && unicode.IsDigit(i.currentChar) {
			number.WriteRune(i.currentChar)
			i.next()
			continue
		}
		break
	}
	val, err := strconv.Atoi(number.String())
	if err != nil {
		panic(fmt.Sprintf("can't convert number %s to int", number.String()))
	}
	return &Token{
		typ: Integer,
		val: val,
	}
}

func (i *Interpreter) getNextToken() *Token {
	for i.currentChar != WrongChar {
		i.skipWhitespaces()

		switch {
		case unicode.IsDigit(i.currentChar):
			return i.readInteger()
		case i.currentChar == '+':
			i.next()
			return &Token{
				typ: Plus,
				val: '+',
			}
		case i.currentChar == '-':
			i.next()
			return &Token{
				typ: Minus,
				val: '-',
			}
		default:
			panic(fmt.Sprintf("unknown type char %T %v", i.currentChar, i.currentChar))
		}
	}
	return nil
}

func (i *Interpreter) consume(tokenTyp tokenType) {
	if i.currentToken.typ == tokenTyp {
		i.currentToken = i.getNextToken()
	} else {
		panic(fmt.Sprintf("error consume token expected %v got %v", tokenTyp, i.currentToken.typ))
	}
}

func (i *Interpreter) Eval() int {
	i.currentToken = i.getNextToken()

	left := i.currentToken
	i.consume(Integer)

	op := i.currentToken

	switch op.typ {
	case Plus:
		i.consume(Plus)
	case Minus:
		i.consume(Minus)
	case Mul:
		i.consume(Mul)
	case Div:
		i.consume(Div)
	default:
		panic("unknown operator")
	}

	right := i.currentToken
	i.consume(Integer)

	switch op.typ {
	case Plus:
		return left.val + right.val
	case Minus:
		return left.val - right.val
	case Mul:
		return left.val * right.val
	case Div:
		return left.val / right.val
	default:
		panic("unknown operand")
	}
}
