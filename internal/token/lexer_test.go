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
			"number", "1",
			[]Token{
				{Kind: kind.Integer, sval: "1"},
				{Kind: kind.Eof, sval: ""},
			},
		},
		{
			"number with skip", "255\n\n",
			[]Token{
				{Kind: kind.Integer, sval: "255"},
				{Kind: kind.Eof, sval: ""},
			},
		},
		{
			"numbers", "255\t\n78\n361\n",
			[]Token{
				{Kind: kind.Integer, sval: "255"},
				{Kind: kind.Integer, sval: "78"},
				{Kind: kind.Integer, sval: "361"},
				{Kind: kind.Eof, sval: ""},
			},
		},
		// punct
		{
			"symbol", "+-*/+",
			[]Token{
				{Kind: kind.Symbol, sval: "+"},
				{Kind: kind.Symbol, sval: "-"},
				{Kind: kind.Symbol, sval: "*"},
				{Kind: kind.Symbol, sval: "/"},
				{Kind: kind.Symbol, sval: "+"},
				{Kind: kind.Eof, sval: ""},
			},
		},
		{
			"symbol with numbers", "255 + 78* 361",
			[]Token{
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
		t.Run(
			tt.name,
			func(t *testing.T) {
				got := Lex(strings.NewReader(tt.input))
				for i, tok := range tt.want {
					g, ok := <-got
					if !ok {
						t.Fatalf("closed before finish")
					}
					if g.Kind != tok.Kind || g.sval != tok.sval {
						t.Errorf("(%d): want %v, got %v", i, tok, g)
					}
				}
			},
		)
	}

}
