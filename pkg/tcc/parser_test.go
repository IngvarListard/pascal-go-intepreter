package tcc

import (
	"bytes"
	"reflect"
	"testing"
)

func TestParser_Parse(t *testing.T) {
	lexer := NewLexer(bytes.NewReader([]byte(" i = 3; ")))
	type fields struct {
		lexer Lexer
	}
	tests := []struct {
		name    string
		fields  fields
		want    *Node
		wantErr bool
	}{
		{
			name: "test1",
			fields: fields{
				lexer: *lexer,
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{
				lexer: tt.fields.lexer,
			}
			got, err := p.Parse()
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() got = %v, want %v", got, tt.want)
			}
		})
	}
}
