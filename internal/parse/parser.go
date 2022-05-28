package parse

import (
	"github.com/lunashade/lang/internal/ast"
	"github.com/lunashade/lang/internal/token"
	"github.com/lunashade/lang/internal/token/kind"
)

type Parser struct {
	stream *token.Stream
}

func Run(ch chan token.Token) (ast.AST, error) {
	p := &Parser{stream: token.NewStream(ch)}
	node, err := p.Root(0)
	if err != nil {
		return nil, err
	}
	return node, nil
}

func (p *Parser) consume(kind kind.Kind, at int) (int, *token.Token) {
	t := p.stream.Look(at)
	if t == nil || t.Kind != kind {
		return at, nil
	}
	return at + 1, t
}
