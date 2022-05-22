package parse

import (
	"errors"

	"github.com/lunashade/lang/internal/ast"
	"github.com/lunashade/lang/internal/token/kind"
)

// PEG
// Root <- Expr
// Expr <- Sum
// Sum <- Add / Sub / Prod
// [Add] <- Prod "+" Sum
// [Sub] <- Prod "-" Sum
// Prod <- Mul / Div / Term
// [Mul] <- Term "*" Prod
// [Div] <- Term "/" Prod
// Term <- ParenExpr / Integer
// [ParenExpr] <- "(" Expr ")"
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
	return p.Sum(pos)
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
	return p.Select(p.ParenExpr, p.Integer)(pos)
}

func (p *Parser) ParenExpr(pos int) (int, ast.AST, error) {
	return p.Concat(
		func(nodes []ast.AST) ast.AST {
			return nodes[1]
		},
		p.Skip(kind.LParen),
		p.Expr,
		p.Skip(kind.RParen),
	)(pos)
}

func (p *Parser) Integer(pos int) (int, ast.AST, error) {
	nx, t := p.Consume(kind.Integer, pos)
	if t == nil {
		return pos, nil, errors.New("not an integer token")
	}
	return nx, &ast.Int{Value: t.IntValue()}, nil
}
