package parse

import (
	"errors"
	"strconv"

	"github.com/lunashade/lang/internal/ast"
	"github.com/lunashade/lang/internal/token/kind"
)

// PEG
// Root <- Function*
// Function <- ident "(" ")" Block
// === statements ===
// Stmt <- Stmt2 / ExprStmt
// Stmt2 <- Semi
// [Semi] <- Expr ";"
// [ExprStmt] <- Expr
// === expressions ===
// Expr <- Assign / Sum
// [Block] <- "{" Stmt2* ExprStmt?  "}"
// [Assign] <- ident "=" Sum
// Sum <- Add / Sub / Prod
// [Add] <- Prod "+" Sum
// [Sub] <- Prod "-" Sum
// Prod <- Mul / Div / Primary
// [Mul] <- Primary "*" Prod
// [Div] <- Primary "/" Prod
// Primary <- Block / ParenExpr / int / ident
// [ParenExpr] <- "(" Expr ")"

// Root parses root node
// PEG: Root <- Function*
func (p *Parser) Root(pos int) (ast.AST, error) {
	_, node, err := p.Repeat(
		func(nodes []ast.AST) ast.AST {
			return &ast.Root{Nodes: nodes}
		},
		p.Select(p.Function),
	)(0)
	if err != nil {
		return nil, err
	}
	if !p.stream.Complete() {
		return nil, errors.New("not at eof")
	}
	return node, nil
}

// Function parses function node
// PEG: Function <- ident "(" ")" Block
func (p *Parser) Function(pos int) (int, ast.AST, error) {
	return p.Concat(
		func(nodes []ast.AST) ast.AST {
			body := nodes[3].(*ast.Block)
			return &ast.Function{
				Name: nodes[0],
				Body: body.Stmts,
			}
		},
		p.Identifier,
		p.Skip(kind.LeftParen),
		p.Skip(kind.RightParen),
		p.Block,
	)(pos)
}

func (p *Parser) Block(pos int) (int, ast.AST, error) {
	return p.Concat(
		func(nodes []ast.AST) ast.AST {
			return nodes[1]
		},
		p.Skip(kind.LeftBrace),
		p.RepeatWithOptionalLast(
			func(nodes []ast.AST) ast.AST {
				return &ast.Block{
					Stmts: nodes,
				}
			},
			p.Stmt2, // no ExprStmt
			p.ExprStmt,
		),
		p.Skip(kind.RightBrace),
	)(pos)
}

func (p *Parser) Stmt(pos int) (int, ast.AST, error) {
	return p.Select(p.Stmt2, p.ExprStmt)(pos)
}
func (p *Parser) Stmt2(pos int) (int, ast.AST, error) {
	return p.Select(p.Semi)(pos)
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
	nx, t := p.consume(kind.Integer, pos)
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
	nx, t := p.consume(kind.Identifier, pos)
	if t == nil {
		return pos, nil, errors.New("not an identifier token")
	}
	return nx, &ast.Ident{Name: t.Sval}, nil
}
