package parse

import (
	"errors"
	"strconv"

	"github.com/lunashade/lang/internal/ast"
	"github.com/lunashade/lang/internal/token/kind"
)

// PEG
// Root <- Function*
// Function <- ident "(" ")" "{" Stmt* "}"
// === statements ===
// Stmt <- Semi / ExprStmt
// [Semi] <- Expr ";"
// [ExprStmt] <- Expr
// === expressions ===
// Expr <- Assign / Sum
// [Assign] <- ident "=" Sum
// Sum <- Add / Sub / Prod
// [Add] <- Prod "+" Sum
// [Sub] <- Prod "-" Sum
// Prod <- Mul / Div / Term
// [Mul] <- Term "*" Prod
// [Div] <- Term "/" Prod
// Term <- ParenExpr / int / ident
// [ParenExpr] <- "(" Expr ")"
func (p *Parser) Root(pos int) (ast.AST, error) {
	_, node, err := p.Function(0)
	if err != nil {
		return nil, err
	}
	p.next()
	if !p.ateof {
		return nil, errors.New("not at eof")
	}
	return &ast.Root{Nodes: []ast.AST{node}}, nil
}


func (p *Parser) Function(pos int) (int, ast.AST, error) {
	return p.Concat(
		func(nodes []ast.AST) ast.AST {
			return &ast.Function{
				Name: nodes[0],
				Body: nodes[4],
			}
		},
		p.Identifier,
		p.Skip(kind.LeftParen),
		p.Skip(kind.RightParen),
		p.Skip(kind.LeftBrace),
		p.Stmt,
		p.Skip(kind.RightBrace),
	)(pos)
}

func (p *Parser) Stmt(pos int) (int, ast.AST, error) {
	return p.Select(p.Semi, p.ExprStmt)(pos)
}

func (p *Parser) ExprStmt(pos int) (int, ast.AST, error) {
	return p.Concat(
		func(nodes []ast.AST) ast.AST {
			return &ast.ExprStmt{
				Expr: nodes[0],
			}
		},
		p.Expr,
	)(pos)
}

func (p *Parser) Semi(pos int) (int, ast.AST, error) {
	return p.Concat(
		func(nodes []ast.AST) ast.AST {
			return &ast.Semi{
				Expr: nodes[0],
			}
		},
		p.Expr,
		p.Skip(kind.Semicolon),
	)(pos)
}

func (p *Parser) Expr(pos int) (int, ast.AST, error) {
	return p.Select(p.Assign, p.Sum)(pos)
}

func (p *Parser) Assign(pos int) (int, ast.AST, error) {
	return p.Concat(
		func(nodes []ast.AST) ast.AST {
			return &ast.BinOp{
				Kind: ast.Assign, LHS: nodes[0], RHS: nodes[2],
			}
		},
		p.Identifier,
		p.Skip(kind.Assign),
		p.Sum,
	)(pos)
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
				Kind: ast.Add, LHS: nodes[0], RHS: nodes[2]}
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
				Kind: ast.Sub, LHS: nodes[0], RHS: nodes[2]}
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
				Kind: ast.Mul, LHS: nodes[0], RHS: nodes[2]}
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
				Kind: ast.Div, LHS: nodes[0], RHS: nodes[2]}
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
		p.Skip(kind.LeftParen),
		p.Expr,
		p.Skip(kind.RightParen),
	)(pos)
}

func (p *Parser) Integer(pos int) (int, ast.AST, error) {
	nx, t := p.Consume(kind.Integer, pos)
	if t == nil {
		return pos, nil, errors.New("not an integer token")
	}
	val, err := strconv.Atoi(t.Sval)
	if err != nil {
		return pos, nil, err
	}
	return nx, &ast.Int{Value: int64(val)}, nil
}

func (p *Parser) Identifier(pos int) (int, ast.AST, error) {
	nx, t := p.Consume(kind.Identifier, pos)
	if t == nil {
		return pos, nil, errors.New("not an identifier token")
	}
	return nx, &ast.Ident{Name: t.Sval}, nil
}
