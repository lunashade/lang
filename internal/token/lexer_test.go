package token

import (
	"strings"
	"testing"

	"github.com/lunashade/lang/internal/token/kind"
)

func TestLex(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []Token
	}{
		// number
		{
			name:  "number",
			input: "1",
			want: []Token{
				{Kind: kind.Integer, sval: "1"},
				{Kind: kind.Eof, sval: ""},
			},
		},
		{
			name:  "number with skip",
			input: "255\n\n",
			want: []Token{
				{Kind: kind.Integer, sval: "255"},
				{Kind: kind.Eof, sval: ""},
			},
		},
		{
			name:  "numbers",
			input: "255\t\n78\n361\n",
			want: []Token{
				{Kind: kind.Integer, sval: "255"},
				{Kind: kind.Integer, sval: "78"},
				{Kind: kind.Integer, sval: "361"},
				{Kind: kind.Eof, sval: ""},
			},
		},
		// punct
		{
			name:  "punctuations",
			input: "+-*/+",
			want: []Token{
				{Kind: kind.Symbol, sval: "+"},
				{Kind: kind.Symbol, sval: "-"},
				{Kind: kind.Symbol, sval: "*"},
				{Kind: kind.Symbol, sval: "/"},
				{Kind: kind.Symbol, sval: "+"},
				{Kind: kind.Eof, sval: ""},
			},
		},
		{
			name:  "punctuations with numbers",
			input: "255 + 78* 361",
			want: []Token{
				{Kind: kind.Integer, sval: "255"},
				{Kind: kind.Symbol, sval: "+"},
				{Kind: kind.Integer, sval: "78"},
				{Kind: kind.Symbol, sval: "*"},
				{Kind: kind.Integer, sval: "361"},
				{Kind: kind.Eof, sval: ""},
			},
		},
	}
	for _, tt := range tests {
		got := Lex(strings.NewReader(tt.input))
		for i, tok := range tt.want {
			g, ok := <-got
			if !ok {
				t.Fatalf("%s: closed before finish", tt.name)
			}
			if g.Kind != tok.Kind || g.sval != tok.sval {
				t.Errorf("%s(%d): want %v, got %v", tt.name, i, tok, g)
			}
		}
	}

}
