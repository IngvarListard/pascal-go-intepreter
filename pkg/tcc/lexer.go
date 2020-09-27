package tcc

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
	"unicode"
)

type Token int

const (
	Num Token = iota + 1
	Id
	If
	Else
	While
	Do
	Lbra
	Rbra
	Lpar
	Rpar
	Plus
	Minus
	Less
	Equal
	Semicolon
	EOF
)

var Symbols = map[string]Token{
	"{": Lbra,
	"}": Rbra,
	"=": Equal,
	";": Semicolon,
	"(": Lpar,
	")": Rpar,
	"+": Plus,
	"-": Minus,
	"<": Less,
}

var Words = map[string]Token{
	"if":    If,
	"else":  Else,
	"do":    Do,
	"while": While,
}

type Lexer struct {
	reader      *bufio.Reader
	currentRune rune

	value int
	tok   Token
}

func (l *Lexer) nextRune() (err error) {
	l.currentRune, _, err = l.reader.ReadRune()
	return
}

func (l *Lexer) nextToken() error {

	have := func(r rune) bool { _, ok := Symbols[string(r)]; return ok }

	for l.tok == -1 {
		switch r := l.currentRune; {
		case r == 0:
			l.tok = EOF
		case unicode.IsSpace(r):
			if err := l.nextRune(); err != nil {
				return err
			}
		case have(r):
			l.tok = Symbols[string(r)]
			if err := l.nextRune(); err != nil {
				return nil
			}
		case unicode.IsDigit(r):
			var bufNum bytes.Buffer
			for unicode.IsDigit(l.currentRune) {
				_, err := bufNum.WriteRune(l.currentRune)
				if err != nil {
					return err
				}
			}
			number, err := strconv.Atoi(bufNum.String())
			if err != nil {
				return err
			}
			l.value = number
			l.tok = Num
		case unicode.IsLetter(r):
			var identBuf bytes.Buffer
			for unicode.IsLetter(l.currentRune) {
				identBuf.WriteRune(unicode.ToLower(l.currentRune))
				if err := l.nextRune(); err != nil {
					return err
				}
			}
			ident := identBuf.String()
			if _, ok := Words[ident]; ok {
				l.tok = Words[ident]
			} else if len(ident) == 1 {
				l.tok = Id
				l.value = int([]rune(ident)[0] - 'a')
			} else {
				return fmt.Errorf("unknown identifier %s", ident)
			}
		default:
			return fmt.Errorf("unexpected symbol %s", string(r))
		}
	}
	return nil
}

func NewLexer(r io.Reader) *Lexer {
	return &Lexer{
		currentRune: ' ',
		reader:      bufio.NewReader(r),
	}
}
