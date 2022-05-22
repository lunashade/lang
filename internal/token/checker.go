package token

import "github.com/lunashade/lang/internal/token/kind"


func isDigit(c rune) bool {
	return '0' <= c && c <= '9'
}

// check if rune of identifier
func isIdent(c rune) bool {
	islower := 'a' <= c && c <= 'z'
	isupper := 'A' <= c && c <= 'Z'
	return islower || isupper
}

func isSymbol(c rune) bool {
	for _, r := range kind.Symbols {
		if c == r {
			return true
		}
	}
	return false
}
