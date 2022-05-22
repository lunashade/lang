package compile

import (
	"bytes"
	"strings"
	"testing"
)

const sample string = `
1 + 1
`

func TestCompile(t *testing.T) {
	var buf bytes.Buffer
	Run(strings.NewReader(sample), &buf)  // check it won't panic
}
