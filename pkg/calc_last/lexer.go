package calc5

import (
	"bytes"
	"fmt"
	"strconv"
	"unicode"
)

const (
	Number TokenTyp = iota
	Plus
	Minus
	Mul
	Div
	Lparen // (
	Rparen // )
	EOF

	NullRune rune = 0
)

// Lexer or Tokenizer
type Lexer struct {
	text        []rune
	currentRune rune
	pos         int
}

func (l *Lexer) next() {
	l.pos++
	if l.pos >= len(l.text) {
		l.currentRune = NullRune
		return
	}
	l.currentRune = l.text[l.pos]
}

func (l *Lexer) skipWhitespace() {
	for unicode.IsSpace(l.currentRune) {
		l.next()
	}
}

func (l *Lexer) getNextToken() *Token {
	for l.currentRune != NullRune {
		switch r := l.currentRune; {
		case unicode.IsSpace(r):
			l.skipWhitespace()
		case unicode.IsDigit(r):
			return l.readInt()
		case r == '+':
			l.next()
			return &Token{typ: Plus, value: r}
		case r == '-':
			l.next()
			return &Token{typ: Minus, value: r}
		case r == '*':
			l.next()
			return &Token{typ: Mul, value: r}
		case r == '/':
			l.next()
			return &Token{typ: Div, value: r}
		case r == '(':
			l.next()
			return &Token{typ: Lparen, value: r}
		case r == ')':
			l.next()
			return &Token{typ: Rparen, value: r}
		default:
			panic(fmt.Sprintf("Unexpected character occurance: %s", string(r)))
		}
	}
	return &Token{typ: EOF, value: NullRune}
}

func (l *Lexer) readInt() *Token {
	var numberBuf bytes.Buffer
	for unicode.IsDigit(l.currentRune) {
		numberBuf.WriteRune(l.currentRune)
		l.next()
	}
	number, _ := strconv.Atoi(numberBuf.String())
	return &Token{typ: Number, value: number}
}
