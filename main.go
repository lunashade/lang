package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/lunashade/lang/internal/token"
)

func main() {
	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		panic("cannot read stdin")
	}
	code := string(b)
	tokens := token.Lex(strings.NewReader(code))

	fmt.Printf("define dso_local i32 @main() #0 {\n")
	for tok := range tokens {
		if tok.Kind == token.Eof {
			break
		}
		if tok.Kind == token.Integer {
			fmt.Printf("\tret i32 %s\n", tok)
		}
	}
	fmt.Printf("}\n")
	return
}
