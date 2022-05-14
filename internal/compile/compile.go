package compile

import (
	"fmt"
	"io"

	"github.com/lunashade/lang/internal/token"
	"github.com/lunashade/lang/internal/token/kind"
)

func Run(r io.Reader, w io.Writer) {
	tokens := token.Lex(r)

	fmt.Fprintf(w, "define dso_local i32 @main() #0 {\n")
	for tok := range tokens {
		if tok.Kind == kind.Eof {
			break
		}
		if tok.Kind == kind.Integer {
			fmt.Fprintf(w, "\tret i32 %s\n", tok)
		}
	}
	fmt.Fprintf(w, "}\n")
}
