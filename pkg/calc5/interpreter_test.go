package calc5

import (
	"fmt"
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
						text: []rune(`
BEGIN
    BEGIN
        nUmbeR := 2;
        a := number;
        b := 10 * a + 10 * number DIV 4;
        c := a - - b
    END;
    x := 11;
END.
`),
						currentRune: '\n',
						pos:         0,
					},
					currentToken: nil,
				},
			},
			want: 5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Interpreter{
				parser: tt.fields.parser,
			}
			i.GlobalScope = make(map[string]interface{})
			i.parser.currentToken = i.parser.lexer.getNextToken()
			if got := i.interpret(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("interpret() = %v, want %v", got, tt.want)
			}
			fmt.Println("GLOBAL SCOPE", i.GlobalScope)
		})
	}
}
