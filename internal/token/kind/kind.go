package kind

// Kind is token's kind
type Kind int

const (
	Invalid Kind = iota
	Eof
	// Literal
	Integer
	// Symbol
	Plus     // '+'
	Minus    // '-'
	Multiply // '*'
	Divide   // '/'
	LParen   // '('
	RParen   // ')'
)

const Symbols = "+-*/()"

func SymbolKind(c rune) Kind {
	for i, r := range Symbols {
		if r == c {
			return Kind(i + int(Plus))
		}
	}
	return Invalid
}
