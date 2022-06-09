package kind

// Kind is token's kind
type Kind int

const (
	Invalid Kind = iota
	Eof
	Identifier
	// Keywords
	KwIf   // "if"
	KwThen // "then"
	KwElse // "else"
	// Literal
	Integer
	String
	// Symbol
	Plus        // '+'
	Minus       // '-'
	Multiply    // '*'
	Divide      // '/'
	Assign      // '='
	LeftParen   // '('
	RightParen  // ')'
	LeftBrace   // '{'
	RightBrace  // '}'
	LessThan    // '<'
	GreaterThan // '>'
	Semicolon   // ';'
	Not         // '!'
)

const Symbols = "+-*/=(){}<>;!"

func SymbolKind(c rune) Kind {
	for i, r := range Symbols {
		if r == c {
			return Kind(i + int(Plus))
		}
	}
	return Invalid
}

var Keywords = []string{
	"if", "then", "else",
}

func KeywordKind(s string) Kind {
	for i, kw := range Keywords {
		if kw == s {
			return Kind(i + int(KwIf))
		}
	}
	return Identifier
}
