package kind

// Kind is token's kind
type Kind int

const (
	Skip Kind = iota
	Eof
	Integer
	Punctuation
)
