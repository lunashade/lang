package token

import (
	"github.com/lunashade/lang/internal/token/kind"
)

var INVALID = Token{Kind: kind.Invalid}

type Token struct {
	Kind kind.Kind
	Sval string
}

func makeToken(kind kind.Kind, sval string) Token {
	return Token{
		Kind: kind,
		Sval: sval,
	}
}
