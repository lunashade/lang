package kind

// Kind is token's kind
type Kind int

const (
	Eof Kind = iota
	Integer
	// Symbol
	Plus
	Minus
	Multiply
	Divide
)
