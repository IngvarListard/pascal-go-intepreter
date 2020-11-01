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
	errorCode
	Token     *calc5.Token
	err       error
	typ       errorType
	callStack []string
}

func (e *Error) String() string {

	var buf bytes.Buffer
	buf.WriteString(string(e.typ))
	buf.WriteString(": ")
	buf.WriteString(string(e.errorCode))
	buf.WriteString(": ")
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

func NewError(errorCode errorCode, token *calc5.Token, typ errorType, err, scopeDescription string) *Error {
	return &Error{
		errorCode: errorCode,
		Token:     token,
		err:       errors.New(err),
		typ:       typ,
		callStack: []string{scopeDescription},
	}
}

func NewLexerError(errorCode errorCode, token *calc5.Token, err, scopeDescription string) *Error {
	return NewError(errorCode, token, LexerError, err, scopeDescription)
}

func NewParserError(errorCode errorCode, token *calc5.Token, err, scopeDescription string) *Error {
	return NewError(errorCode, token, ParserError, err, scopeDescription)
}

func NewSemanticError(errorCode errorCode, token *calc5.Token, err, scopeDescription string) *Error {
	return NewError(errorCode, token, SemanticError, err, scopeDescription)
}
