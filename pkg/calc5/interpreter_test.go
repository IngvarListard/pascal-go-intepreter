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
program Main;
   var b, x, y : real;
   var z : integer;

   procedure AlphaA(a : integer);
      var b : integer;

      procedure Beta(c : integer);
         var y : integer;

         procedure Gamma(c : integer);
            var x : integer;
         begin { Gamma }
            x := a + b + c + x + y + z;
         end;  { Gamma }

      begin { Beta }

      end;  { Beta }

   begin { AlphaA }

   end;  { AlphaA }

   procedure AlphaB(a : integer);
      var c : real;
   begin { AlphaB }
      c := a + b;
   end;  { AlphaB }

begin { Main }
end.  { Main }
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
