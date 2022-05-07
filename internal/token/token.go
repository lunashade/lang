package token

// Kind is token's kind
type Kind int

const (
	Skip Kind = iota
	Eof
	Integer
	Punctuation
)

type Token struct {
	Kind Kind
	sval string
}

func makeToken(kind Kind, sval string) Token {
	return Token{
		Kind: kind,
		sval: sval,
	}
}

func (t Token) String() string {
	return t.sval
}
