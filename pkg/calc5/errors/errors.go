package errors

import (
	"bytes"
	"errors"
	calc5 "github.com/IngvarListard/pascal-go-intepreter/pkg/calc_last"
)

type errorCode string
type errorType string

const (
	UnexpectedToken errorCode = "Unexpected token"
	IDNotFound                = "ID not found"
	DuplicateID               = "Duplicate ID"

	LexerError    errorType = "LexerError"
	ParserError             = "ParserError"
	SemanticError           = "SemanticError"
)

type Error struct {
	*errorCode
	Token     *calc5.Token
	err       error
	typ       errorType
	callStack []string
}

func (e *Error) String() string {

	var buf bytes.Buffer
	buf.WriteString(string(e.typ))
	buf.WriteString(": ")
	if e.errorCode != nil {
		buf.WriteString(string(*e.errorCode))
		buf.WriteString(": ")
	}
	buf.WriteString(e.err.Error())

	for _, scopeDesc := range e.callStack {
		buf.WriteString(scopeDesc)
		buf.WriteString("\n")
	}

	return buf.String()
}

func (e *Error) Error() string {
	return e.err.Error()
}

func (e *Error) Through(scopeDescription string) *Error {
	e.callStack = append(e.callStack, scopeDescription)
	return e
}

func NewError(typ errorType, err, scopeDescription string, options ...Option) *Error {
	e := &Error{
		err:       errors.New(err),
		typ:       typ,
		callStack: []string{scopeDescription},
	}
	for _, option := range options {
		option(e)
	}
	return e
}

type Option func(*Error)

func Token(token *calc5.Token) Option {
	return func(e *Error) {
		e.Token = token
	}
}

func ErrorCode(ec *errorCode) Option {
	return func(e *Error) {
		e.errorCode = ec
	}
}

func NewLexerError(err, scopeDescription string, options ...Option) *Error {
	return NewError(LexerError, err, scopeDescription, options...)
}

func NewParserError(err, scopeDescription string, options ...Option) *Error {
	return NewError(ParserError, err, scopeDescription, options...)
}

func NewSemanticError(err, scopeDescription string, options ...Option) *Error {
	return NewError(SemanticError, err, scopeDescription, options...)
}
