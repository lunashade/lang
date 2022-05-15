package compile

import (
	"io"

	"github.com/lunashade/lang/internal/gen"
	"github.com/lunashade/lang/internal/parse"
	"github.com/lunashade/lang/internal/token"
)

func Run(r io.Reader, w io.Writer) {
	tokens := token.Lex(r)
	node, err := parse.Run(tokens)
	if err != nil {
		panic(err)
	}
	err = gen.Run(w, node)
	if err != nil {
		panic(err)
	}
}
