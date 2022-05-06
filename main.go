package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func main() {
	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		panic("cannot read stdin")
	}
	code := string(b)

	num, err := strconv.Atoi(strings.TrimSuffix(code, "\n"))
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		panic("unknown integer")
	}

	fmt.Printf("define dso_local i32 @main() #0 {\n")
	fmt.Printf("\tret i32 %d\n", num)
	fmt.Printf("}\n")
	return
}
