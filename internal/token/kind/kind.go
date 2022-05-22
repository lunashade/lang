package kind

// Kind is token's kind
type Kind int

const (
	Invalid Kind = iota
	Eof
	Identifier
	Keyword
	// Literal
	Integer
	String
	// Symbol
	Plus        // '+'
	Minus       // '-'
	Multiply    // '*'
	Divide      // '/'
	LeftParen   // '('
	RightParen  // ')'
	Assign      // '='
	LessThan    // '<'
	GreaterThan // '>'
)

const Symbols = "+-*/()=<>"

func SymbolKind(c rune) Kind {
	for i, r := range Symbols {
		if r == c {
			return Kind(i + int(Plus))
		}
	}
	return Invalid
}
