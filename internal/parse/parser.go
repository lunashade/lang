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

func Parse(ch chan token.Token) ast.AST {
	p := &Parser{
		ch:     ch,
		tokens: make([]token.Token, 0),
	}
	ast, err := p.Root(0)
	if err != nil {
		panic(err)
	}
	return ast
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
	// TODO: will be like below
	_, ast, err := p.Sum(0)
	if !p.ateof {
		panic("not at eof")
	}
	if err != nil {
		panic(err)
	}
	return ast, nil
}

type NonTerminal func(int) (int, ast.AST, error)

func (p *Parser) Sum(pos int) (int, ast.AST, error) {
	var nx int
	var ast ast.AST
	var err error

	nx, ast, err = p.Add(pos)
	if err == nil {
		return nx, ast, nil
	}
	nx, ast, err = p.Sub(pos)
	if err == nil {
		return nx, ast, nil
	}
	nx, ast, err = p.Prod(pos)
	if err == nil {
		return nx, ast, nil
	}
	return pos, nil, err
}
func (p *Parser) Prod(pos int) (int, ast.AST, error) {
	var nx int
	var ast ast.AST
	var err error

	nx, ast, err = p.Mul(pos)
	if err == nil {
		return nx, ast, nil
	}
	nx, ast, err = p.Div(pos)
	if err == nil {
		return nx, ast, nil
	}
	nx, ast, err = p.Term(pos)
	if err == nil {
		return nx, ast, nil
	}
	return pos, nil, err
}

func (p *Parser) Add(pos int) (int, ast.AST, error) {
	// read LHS
	nx, lhs, err := p.Prod(pos)
	if err != nil {
		return pos, nil, err
	}
	// skip "+"
	nx, t := p.Consume(kind.Plus, nx)
	if t == nil {
		return pos, nil, errors.New("invalid token")
	}
	// read RHS
	nx, rhs, err := p.Sum(nx)

	// return
	return nx, &ast.Add{LHS: lhs, RHS: rhs}, nil
}
func (p *Parser) Sub(pos int) (int, ast.AST, error) {
	// read LHS
	nx, lhs, err := p.Prod(pos)
	if err != nil {
		return pos, nil, err
	}
	// skip "+"
	nx, t := p.Consume(kind.Minus, nx)
	if t == nil {
		return pos, nil, errors.New("invalid token")
	}
	// read RHS
	nx, rhs, err := p.Sum(nx)

	// return
	return nx, &ast.Sub{LHS: lhs, RHS: rhs}, nil
}

func (p *Parser) Mul(pos int) (int, ast.AST, error) {
	// read LHS
	nx, lhs, err := p.Term(pos)
	if err != nil {
		return pos, nil, err
	}
	// skip "+"
	nx, t := p.Consume(kind.Multiply, nx)
	if t == nil {
		return pos, nil, errors.New("invalid token")
	}
	// read RHS
	nx, rhs, err := p.Prod(nx)

	// return
	return nx, &ast.Mul{LHS: lhs, RHS: rhs}, nil
}

func (p *Parser) Div(pos int) (int, ast.AST, error) {
	// read LHS
	nx, lhs, err := p.Term(pos)
	if err != nil {
		return pos, nil, err
	}
	// skip "+"
	nx, t := p.Consume(kind.Divide, nx)
	if t == nil {
		return pos, nil, errors.New("invalid token")
	}
	// read RHS
	nx, rhs, err := p.Prod(nx)

	// return
	return nx, &ast.Div{LHS: lhs, RHS: rhs}, nil
}

func (p *Parser) Term(pos int) (int, ast.AST, error) {
	return p.Integer(pos)
}
func (p *Parser) Integer(pos int) (int, ast.AST, error) {
	nx, t := p.Consume(kind.Integer, pos)
	if t == nil {
		return pos, nil, errors.New("not an integer")
	}
	return nx, &ast.Int{Value: t.IntValue()}, nil
}
