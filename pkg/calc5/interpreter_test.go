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
PROGRAM Part12;
VAR
   a : INTEGER;

PROCEDURE P1;
VAR
   a : REAL;
   k : INTEGER;

   PROCEDURE P2;
   VAR
      a, z : INTEGER;
   BEGIN {P2}
      z := 777;
   END;  {P2}

BEGIN {P1}

END;  {P1}

BEGIN {Part12}
   a := 10;
END.  {Part12}
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
			fmt.Println("SYMBOLS", i.Symbols.symbols)
		})
	}
}
