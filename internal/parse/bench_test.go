package parse

import (
	"strings"
	"testing"

	"github.com/lunashade/lang/internal/token"
)

var code = `
main() {
	((((((1))))))
}
`

func BenchmarkParseExpr(b *testing.B) {
	b.StartTimer()
	Run(token.Lex(strings.NewReader(code)))
	b.StopTimer()
}
