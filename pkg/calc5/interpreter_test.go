package calc5

import (
	"reflect"
	"testing"
)

func TestInterpreter_interpret(t *testing.T) {
	type fields struct {
		parser *Parser
	}
	tests := []struct {
		name   string
		fields fields
		want   interface{}
	}{
		{
			name: "simple_sum",
			fields: fields{
				parser: &Parser{
					lexer: &Lexer{
						text:        []rune("2 + 3"),
						currentRune: '2',
						pos:         0,
					},
					currentToken: nil,
				},
			},
			want: 5,
		},
		{
			name: "complex_expr",
			fields: fields{
				parser: &Parser{
					lexer: &Lexer{
						text:        []rune("7 + 3 * (10 / (12 / (3 + 1) - 1)) / (2 + 3) - 5 - 3 + (8)"),
						currentRune: '7',
						pos:         0,
					},
					currentToken: nil,
				},
			},
			want: 10,
		},
		{
			name: "multi_bracket",
			fields: fields{
				parser: &Parser{
					lexer: &Lexer{
						text:        []rune("7 + (((3 + 2)))"),
						currentRune: '7',
						pos:         0,
					},
					currentToken: nil,
				},
			},
			want: 12,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Interpreter{
				parser: tt.fields.parser,
			}
			i.parser.currentToken = i.parser.lexer.getNextToken()
			if got := i.interpret(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("interpret() = %v, want %v", got, tt.want)
			}
		})
	}
}
