package token

func isDigit(c rune) bool {
	return '0' <= c && c <= '9'
}

func isSymbol(c rune) bool {
	puncts := "+-*/"
	for _, r := range puncts {
		if c == r {
			return true
		}
	}
	return false
}
