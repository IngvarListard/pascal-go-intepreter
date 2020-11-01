package calc5

import (
	"bytes"
	"fmt"
	"github.com/IngvarListard/pascal-go-intepreter/pkg/calc5/errors"
	"strconv"
	"strings"
	"unicode"
)

const (
	// single character token types
	Colon    TokenTyp = iota // ":"
	Comma                    // ","
	FloatDiv                 // "/"
	Plus                     // "+"
	Minus                    // "="
	Mul                      // "*"
	Lparen                   // "("
	Rparen                   // ")"
	Dot                      // "."
	Semi                     // ";"
	// reserved words
	Begin      // "BEGIN"
	Program    // "PROGRAM"
	VarT       // "VAR"
	Integer    // "INTEGER"
	Real       // "REAL"
	IntegerDiv // "DIV"
	Procedure  // "PROCEDURE"
	End        // "END"
	// misc
	Id           // "ID"
	IntegerConst // "INTEGER_CONST"
	RealConst    // "REAL_CONST"
	Assign       // ":="
	EOF          // "EOF"

	NullRune rune = 0
)

var TokenTypes = [...]string{
	// single character token types
	Colon:    ":",
	Comma:    ",",
	FloatDiv: "/",
	Plus:     "+",
	Minus:    "=",
	Mul:      "*",
	Lparen:   "(",
	Rparen:   ")",
	Dot:      ".",
	Semi:     ";",
	// reserved words
	Program:    "PROGRAM",
	VarT:       "VAR",
	Integer:    "INTEGER",
	Real:       "REAL",
	IntegerDiv: "DIV",
	Procedure:  "PROCEDURE",
	Begin:      "BEGIN",
	End:        "END",
	// misc
	Id:           "ID",
	IntegerConst: "INTEGER_CONST",
	RealConst:    "REAL_CONST",
	Assign:       ":=",
	EOF:          "EOF",
}

var ReservedKeywords = map[string]*Token{
	"program":   {typ: Program, value: "program"},
	"var":       {typ: VarT, value: "var"},
	"integer":   {typ: Integer, value: "integer"},
	"real":      {typ: Real, value: "real"},
	"div":       {typ: IntegerDiv, value: "div"},
	"procedure": {typ: Procedure, value: "procedure"},
	"begin":     {typ: Begin, value: "begin"},
	"end":       {typ: End, value: "end"},
}

// TODO
//var ReservedKeywords = buildReservedKeywords()
//
//func buildReservedKeywords() map[string]*Token {
//	startIndex := Program
//	endIndex := End
//
//	reservedKeywords := make(map[string]*Token, startIndex-endIndex)
//	for i := startIndex; i <= endIndex; i++ {
//		val := TokenTypes[i]
//		reservedKeywords[strings.ToLower(val)] = &Token{typ: i, value: strings.ToLower(val)}
//	}
//	return reservedKeywords
//}

// Lexer or Tokenizer
type Lexer struct {
	text        []rune
	currentRune rune
	pos         int
	lineno      int
	column      int
}

func (l *Lexer) next() {
	if l.currentRune == '\n' {
		l.lineno++
		l.column = 0
	}

	l.pos++
	if l.pos >= len(l.text) {
		l.currentRune = NullRune
		return
	}
	l.currentRune = l.text[l.pos]
	l.column++
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
			l.panic(fmt.Sprintf("Unexpected character occurance: %s", string(r)), "getNextToken")
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
			l.panic(fmt.Sprintf(`real number parsing from string "%s" error: %v`, numberBuf.String(), err), "readNumber")
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

func (l *Lexer) panic(err, context string) {
	msg := fmt.Sprintf("Lexer error on %s: line: %v column: %v: %s", string(l.currentRune), l.lineno, l.column, err)
	panic(errors.NewLexerError(msg, context).String())
}
