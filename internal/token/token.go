package token

import (
	"strconv"

	"github.com/lunashade/lang/internal/token/kind"
)

var INVALID = Token{Kind: kind.Invalid}

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

func (t Token) IntValue() int {
	if t.Kind != kind.Integer {
		return 0
	}
	val, err := strconv.Atoi(t.sval)
	if err != nil {
		panic(err)
	}
	return val
}
