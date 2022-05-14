package parse

import (
	"errors"

	"github.com/lunashade/lang/internal/ast"
	"github.com/lunashade/lang/internal/token"
	"github.com/lunashade/lang/internal/token/kind"
)

type Parser struct {
	ch     chan token.Token
	tokens []token.Token
	ateof  bool
}

// next reads next token from channel
func (p *Parser) next() {
	if p.ateof {
		// TODO: return error if needs
		return
	}
	tok := <-p.ch
	if tok.Kind == kind.Eof {
		p.ateof = true
	}
	p.tokens = append(p.tokens, tok)

}
func (p *Parser) look(at int) *token.Token {
	if at < len(p.tokens) {
		return &(p.tokens[at])
	}
	for !p.ateof {
		p.next()
		if at < len(p.tokens) {
			return &(p.tokens[at])
		}
	}
	return nil
}

func (p *Parser) Consume(kind kind.Kind, at int) (int, *token.Token) {
	t := p.look(at)
	if t == nil || t.Kind != kind {
		return at, nil
	}
	return at + 1, p.look(at)
}

func (p *Parser) Skip(kind kind.Kind) NonTerminal {
	return func(pos int) (int, ast.AST, error) {
		nx, t := p.Consume(kind, pos)
		if t == nil {
			return pos, nil, errors.New("invalid token")
		}
		return nx, nil, nil
	}
}

func Parse(ch chan token.Token) ast.AST {
	p := &Parser{
		ch:     ch,
		tokens: make([]token.Token, 0),
	}
	node, err := p.Root(0)
	if err != nil {
		panic(err)
	}
	return node
}

// PEG
// Root <- Add
// Sum <- Add / Sub / Prod
// [Add] <- Prod "+" Sum
// [Sub] <- Prod "-" Sum
// Prod <- Mul / Div / Term
// [Mul] <- Term "*" Prod
// [Div] <- Term "/" Prod
// Term <- Integer
func (p *Parser) Root(pos int) (ast.AST, error) {
	_, node, err := p.Sum(0)
	if !p.ateof {
		panic("not at eof")
	}
	if err != nil {
		panic(err)
	}
	return node, nil
}

type NonTerminal func(int) (int, ast.AST, error)

func (p *Parser) Select(cands ...NonTerminal) NonTerminal {
	return func(pos int) (int, ast.AST, error) {
		var nx int
		var node ast.AST
		var err error
		for _, cand := range cands {
			nx, node, err = cand(pos)
			if err == nil {
				return nx, node, nil
			}
		}
		// TODO: wraps error
		return pos, nil, err
	}
}

type Merger func([]ast.AST) ast.AST

func (p *Parser) Concat(m Merger, cands ...NonTerminal) NonTerminal {
	return func(pos int) (int, ast.AST, error) {
		var nx int = pos
		var node ast.AST
		var err error

		nodes := make([]ast.AST, 0)
		for _, cand := range cands {
			nx, node, err = cand(nx)
			if err != nil {
				return pos, nil, err
			}
			nodes = append(nodes, node)
		}
		return nx, m(nodes), nil
	}

}

func (p *Parser) Sum(pos int) (int, ast.AST, error) {
	return p.Select(p.Add, p.Sub, p.Prod)(pos)
}
func (p *Parser) Prod(pos int) (int, ast.AST, error) {
	return p.Select(p.Mul, p.Div, p.Term)(pos)
}

func (p *Parser) Add(pos int) (int, ast.AST, error) {
	return p.Concat(
		func(nodes []ast.AST) ast.AST {
			return &ast.Add{LHS: nodes[0], RHS: nodes[2]}
		},
		p.Prod,
		p.Skip(kind.Plus),
		p.Sum,
	)(pos)
}

func (p *Parser) Sub(pos int) (int, ast.AST, error) {
	return p.Concat(
		func(nodes []ast.AST) ast.AST {
			return &ast.Sub{LHS: nodes[0], RHS: nodes[2]}
		},
		p.Prod,
		p.Skip(kind.Minus),
		p.Sum,
	)(pos)
}

func (p *Parser) Mul(pos int) (int, ast.AST, error) {
	return p.Concat(
		func(nodes []ast.AST) ast.AST {
			return &ast.Mul{LHS: nodes[0], RHS: nodes[2]}
		},
		p.Term,
		p.Skip(kind.Multiply),
		p.Prod,
	)(pos)
}

func (p *Parser) Div(pos int) (int, ast.AST, error) {
	return p.Concat(
		func(nodes []ast.AST) ast.AST {
			return &ast.Div{LHS: nodes[0], RHS: nodes[2]}
		},
		p.Term,
		p.Skip(kind.Divide),
		p.Prod,
	)(pos)
}

func (p *Parser) Term(pos int) (int, ast.AST, error) {
	return p.Select(p.Integer)(pos)
}
func (p *Parser) Integer(pos int) (int, ast.AST, error) {
	nx, t := p.Consume(kind.Integer, pos)
	if t == nil {
		return pos, nil, errors.New("not an integer")
	}
	return nx, &ast.Int{Value: t.IntValue()}, nil
}
