package token

import (
	"github.com/lunashade/lang/internal/token/kind"
)

type Token struct {
	Kind kind.Kind
	sval string
}

func makeToken(kind kind.Kind, sval string) Token {
	return Token{
		Kind: kind,
		sval: sval,
	}
}

func (t Token) String() string {
	return t.sval
}
