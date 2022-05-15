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

func Run(ch chan token.Token) (ast.AST, error) {
	p := &Parser{ch: ch}
	node, err := p.Root(0)
	if err != nil {
		return nil, err
	}
	return node, nil
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

// look at the token at the position
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
	return at + 1, t
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
	_, node, err := p.Expr(0)
	if err != nil {
		return nil, err
	}
	if !p.ateof {
		return nil, errors.New("not at eof")
	}
	return &ast.Root{Nodes: []ast.AST{node}}, nil
}

func (p *Parser) Expr(pos int) (int, ast.AST, error) {
	nx, node, err := p.Sum(pos)
	if err != nil {
		return pos, nil, err
	}
	return nx, &ast.Expr{Node: node}, nil
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
			return &ast.BinOp{
				Type: ast.Add, LHS: nodes[0], RHS: nodes[2]}
		},
		p.Prod,
		p.Skip(kind.Plus),
		p.Sum,
	)(pos)
}

func (p *Parser) Sub(pos int) (int, ast.AST, error) {
	return p.Concat(
		func(nodes []ast.AST) ast.AST {
			return &ast.BinOp{
				Type: ast.Sub, LHS: nodes[0], RHS: nodes[2]}
		},
		p.Prod,
		p.Skip(kind.Minus),
		p.Sum,
	)(pos)
}

func (p *Parser) Mul(pos int) (int, ast.AST, error) {
	return p.Concat(
		func(nodes []ast.AST) ast.AST {
			return &ast.BinOp{
				Type: ast.Mul, LHS: nodes[0], RHS: nodes[2]}
		},
		p.Term,
		p.Skip(kind.Multiply),
		p.Prod,
	)(pos)
}

func (p *Parser) Div(pos int) (int, ast.AST, error) {
	return p.Concat(
		func(nodes []ast.AST) ast.AST {
			return &ast.BinOp{
				Type: ast.Div, LHS: nodes[0], RHS: nodes[2]}
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
		return pos, nil, errors.New("not an integer token")
	}
	return nx, &ast.Int{Value: t.IntValue()}, nil
}
