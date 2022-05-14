package main

import (
	"os"

	"github.com/lunashade/lang/internal/compile"
)

func main() {
	compile.Run(os.Stdin, os.Stdout)
	return
}
