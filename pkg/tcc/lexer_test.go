package tcc

import (
	"bufio"
	"bytes"
	"testing"
)

func TestLexer_getNextToken(t *testing.T) {
	type fields struct {
		reader      *bufio.Reader
		currentRune rune
		value       int
		tok         Token
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "test1",
			fields: fields{
				reader:      bufio.NewReader(bytes.NewReader([]byte("a = 3"))),
				currentRune: ' ',
				value:       0,
				tok:         0,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Lexer{
				reader:      tt.fields.reader,
				currentRune: tt.fields.currentRune,
				value:       tt.fields.value,
				tok:         tt.fields.tok,
			}
			if err := l.getNextToken(); (err != nil) != tt.wantErr {
				t.Errorf("getNextToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
