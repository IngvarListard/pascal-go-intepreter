package calc4

import "testing"

func TestInterpreter_Expr(t *testing.T) {
	type fields struct {
		parser        *Parser
		currentLexeme Lexeme
	}
	tests := []struct {
		name    string
		fields  fields
		want    int
		wantErr bool
	}{
		{
			name: "simple_sum",
			fields: fields{
				parser: &Parser{
					text:        []rune("3 + 2"),
					pos:         0,
					currentRune: '3',
				},
				currentLexeme: nil,
			},
			want:    5,
			wantErr: false,
		},
		{
			name: "simple_sub",
			fields: fields{
				parser: &Parser{
					text:        []rune("3 - 2"),
					pos:         0,
					currentRune: '3',
				},
				currentLexeme: nil,
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "simple_mul",
			fields: fields{
				parser: &Parser{
					text:        []rune("3 * 2"),
					pos:         0,
					currentRune: '3',
				},
				currentLexeme: nil,
			},
			want:    6,
			wantErr: false,
		},
		{
			name: "simple_div",
			fields: fields{
				parser: &Parser{
					text:        []rune("6 / 3"),
					pos:         0,
					currentRune: '6',
				},
				currentLexeme: nil,
			},
			want:    2,
			wantErr: false,
		},
		{
			name: "complex",
			fields: fields{
				parser: &Parser{
					text:        []rune(" 3 + 2 * 6 / 4 "),
					pos:         0,
					currentRune: ' ',
				},
				currentLexeme: nil,
			},
			want:    6,
			wantErr: false,
		},
		{
			name: "simple_brackets",
			fields: fields{
				parser: &Parser{
					text:        []rune(" 2 * (2 + 5) "),
					pos:         0,
					currentRune: ' ',
				},
				currentLexeme: nil,
			},
			want:    14,
			wantErr: false,
		},
		{
			name: "brackets_formula",
			fields: fields{
				parser: &Parser{
					text:        []rune("2 * ((2 + 5) + ((3 * 8) * (2 + 5))"),
					pos:         0,
					currentRune: '2',
				},
				currentLexeme: nil,
			},
			want:    350,
			wantErr: false,
		},
		{
			name: "brackets_formula_3",
			fields: fields{
				parser: &Parser{
					text:        []rune("(8 + 10) / 2"),
					pos:         0,
					currentRune: '(',
				},
				currentLexeme: nil,
			},
			want:    9,
			wantErr: false,
		},
		//{
		//	name: "brackets_formula_2",
		//	fields: fields{
		//		parser: &Parser{
		//			text:        []rune("7 + 3 * (10 / (12 / (3 + 1) - 1))"),
		//			pos:         0,
		//			currentRune: '7',
		//		},
		//		currentLexeme: nil,
		//	},
		//	want:    22,
		//	wantErr: false,
		//},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Interpreter{
				parser:        tt.fields.parser,
				currentLexeme: tt.fields.currentLexeme,
			}
			var err error
			i.parser.skipWhitespace()
			i.currentLexeme, err = i.parser.getNextLexeme()
			if err != nil {
				t.Fatalf("%v", err)
				return
			}
			got, err := i.Term()
			if i.parser.pos != -1 {
				t.Errorf("IS NOT EOF")
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("Expr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Expr() got = %v, want %v", got, tt.want)
			}
		})
	}
}
