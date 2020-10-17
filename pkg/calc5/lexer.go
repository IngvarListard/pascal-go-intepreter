package calc5

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

const (
	Program TokenTyp = iota
	VarT
	Colon // :
	Comma // ,
	Integer
	Real
	IntegerConst
	RealConst
	IntegerDiv // DIV
	FloatDiv   // /
	Number
	Plus
	Minus
	Mul
	Div
	Lparen // (
	Rparen // )
	Begin
	End
	Dot
	Semi
	Id
	Assign
	Procedure
	EOF

	NullRune rune = 0
)

var ReservedKeywords = map[string]*Token{
	"program":   {typ: Program, value: "program"},
	"var":       {typ: VarT, value: "var"},
	"integer":   {typ: Integer, value: "integer"},
	"real":      {typ: Real, value: "real"},
	"begin":     {typ: Begin, value: "begin"},
	"end":       {typ: End, value: "end"},
	"div":       {typ: IntegerDiv, value: "div"},
	"procedure": {typ: Procedure, value: "procedure"},
}

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

func (l *Lexer) skipComment() {
	for l.currentRune != '}' {
		l.next()
	}
	l.next()
}

func (l *Lexer) getNextToken() *Token {
	for l.currentRune != NullRune {
		switch r := l.currentRune; {
		case r == '{':
			l.next()
			l.skipComment()
			continue
		case r == ':' && l.peek() == '=':
			l.next()
			l.next()
			return &Token{typ: Assign, value: r}
		case r == ':':
			l.next()
			return &Token{typ: Colon, value: r}
		case r == ',':
			l.next()
			return &Token{typ: Comma, value: r}
		case unicode.IsLetter(r) || r == '_':
			return l.id()
		case unicode.IsSpace(r):
			l.skipWhitespace()
		case unicode.IsDigit(r):
			return l.readNumber()
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
			return &Token{typ: FloatDiv, value: r}
		case r == '(':
			l.next()
			return &Token{typ: Lparen, value: r}
		case r == ')':
			l.next()
			return &Token{typ: Rparen, value: r}
		case r == ';':
			l.next()
			return &Token{typ: Semi, value: r}
		case r == '.':
			l.next()
			return &Token{typ: Dot, value: r}
		default:
			panic(fmt.Sprintf("Unexpected character occurance: %s", string(r)))
		}
	}
	return &Token{typ: EOF, value: NullRune}
}

func (l *Lexer) readNumber() *Token {
	var numberBuf bytes.Buffer
	for unicode.IsDigit(l.currentRune) && l.currentRune != NullRune {
		numberBuf.WriteRune(l.currentRune)
		l.next()
	}

	if l.currentRune == '.' {
		numberBuf.WriteRune(l.currentRune)
		l.next()

		for unicode.IsDigit(l.currentRune) && l.currentRune != NullRune {
			numberBuf.WriteRune(l.currentRune)
			l.next()
		}

		realNumber, err := strconv.ParseFloat(numberBuf.String(), 64)
		if err != nil {
			panic(fmt.Sprintf(`real number parsing from string "%s" error: %v`, numberBuf.String(), err))
		}
		return &Token{typ: RealConst, value: realNumber}
	}

	number, _ := strconv.Atoi(numberBuf.String())
	return &Token{typ: IntegerConst, value: number}
}

func (l *Lexer) peek() rune {
	pos := l.pos + 1
	if pos > len(l.text)-1 {
		return NullRune
	}
	return l.text[pos]
}

func (l *Lexer) id() *Token {
	var result bytes.Buffer
	for l.currentRune != NullRune && (unicode.IsDigit(l.currentRune) || unicode.IsLetter(l.currentRune)) || l.currentRune == '_' {
		result.WriteRune(l.currentRune)
		l.next()
	}

	id := strings.ToLower(result.String())
	tok, ok := ReservedKeywords[id]
	if !ok {
		tok = &Token{typ: Id, value: id}
	}
	return tok
}
