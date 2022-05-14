package kind

// Kind is token's kind
type Kind int

const (
	Invalid Kind = iota
	Eof
	// Literal
	Integer
	// Symbol
	Plus
	Minus
	Multiply
	Divide
)
