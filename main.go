package main

import (
	"fmt"
	"os"

	"github.com/lunashade/lang/internal/token"
	"github.com/lunashade/lang/internal/token/kind"
)

func main() {
	tokens := token.Lex(os.Stdin)

	fmt.Printf("define dso_local i32 @main() #0 {\n")
	for tok := range tokens {
		if tok.Kind == kind.Eof {
			break
		}
		if tok.Kind == kind.Integer {
			fmt.Printf("\tret i32 %s\n", tok)
		}
	}
	fmt.Printf("}\n")
	return
}
