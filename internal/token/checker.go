package token

import "github.com/lunashade/lang/internal/token/kind"


func isDigit(c rune) bool {
	return '0' <= c && c <= '9'
}

func isSymbol(c rune) bool {
	for _, r := range kind.Symbols {
		if c == r {
			return true
		}
	}
	return false
}
