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
				{Kind: kind.Integer, Sval: "1"},
				{Kind: kind.Eof, Sval: ""},
			},
		},
		{
			"number with skip", "255\n\n",
			[]Token{
				{Kind: kind.Integer, Sval: "255"},
				{Kind: kind.Eof, Sval: ""},
			},
		},
		{
			"numbers", "255\t\n78\n361\n",
			[]Token{
				{Kind: kind.Integer, Sval: "255"},
				{Kind: kind.Integer, Sval: "78"},
				{Kind: kind.Integer, Sval: "361"},
				{Kind: kind.Eof, Sval: ""},
			},
		},
		// punct
		{
			"symbol", "+-*/=(){}<>;!",
			[]Token{
				{Kind: kind.Plus, Sval: "+"},
				{Kind: kind.Minus, Sval: "-"},
				{Kind: kind.Multiply, Sval: "*"},
				{Kind: kind.Divide, Sval: "/"},
				{Kind: kind.Assign, Sval: "="},
				{Kind: kind.LeftParen, Sval: "("},
				{Kind: kind.RightParen, Sval: ")"},
				{Kind: kind.LeftBrace, Sval: "{"},
				{Kind: kind.RightBrace, Sval: "}"},
				{Kind: kind.LessThan, Sval: "<"},
				{Kind: kind.GreaterThan, Sval: ">"},
				{Kind: kind.Semicolon, Sval: ";"},
				{Kind: kind.Not, Sval: "!"},
				{Kind: kind.Eof, Sval: ""},
			},
		},
		{
			"symbol with numbers", "255 + 78* 361",
			[]Token{
				{Kind: kind.Integer, Sval: "255"},
				{Kind: kind.Plus, Sval: "+"},
				{Kind: kind.Integer, Sval: "78"},
				{Kind: kind.Multiply, Sval: "*"},
				{Kind: kind.Integer, Sval: "361"},
				{Kind: kind.Eof, Sval: ""},
			},
		},
		{
			"paren", "1 + (2 * 3)",
			[]Token{
				{Kind: kind.Integer, Sval: "1"},
				{Kind: kind.Plus, Sval: "+"},
				{Kind: kind.LeftParen, Sval: "("},
				{Kind: kind.Integer, Sval: "2"},
				{Kind: kind.Multiply, Sval: "*"},
				{Kind: kind.Integer, Sval: "3"},
				{Kind: kind.RightParen, Sval: ")"},
				{Kind: kind.Eof, Sval: ""},
			},
		},
		{
			"identifier", "abc efg",
			[]Token{
				{Kind: kind.Identifier, Sval: "abc"},
				{Kind: kind.Identifier, Sval: "efg"},
				{Kind: kind.Eof, Sval: ""},
			},
		},
		{
			"simple function", "main(){1+1} ",
			[]Token{
				{Kind: kind.Identifier, Sval: "main"},
				{Kind: kind.LeftParen, Sval: "("},
				{Kind: kind.RightParen, Sval: ")"},
				{Kind: kind.LeftBrace, Sval: "{"},
				{Kind: kind.Integer, Sval: "1"},
				{Kind: kind.Plus, Sval: "+"},
				{Kind: kind.Integer, Sval: "1"},
				{Kind: kind.RightBrace, Sval: "}"},
				{Kind: kind.Eof, Sval: ""},
			},
		},
		{
			"keywords", "if then else ifs",
			[]Token{
				{Kind: kind.KwIf, Sval: "if"},
				{Kind: kind.KwThen, Sval: "then"},
				{Kind: kind.KwElse, Sval: "else"},
				{Kind: kind.Identifier, Sval: "ifs"},
				{Kind: kind.Eof, Sval: ""},
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
					if g.Kind != tok.Kind || g.Sval != tok.Sval {
						t.Errorf("(%d): want %v, got %v", i, tok, g)
					}
				}
			},
		)
	}

}
