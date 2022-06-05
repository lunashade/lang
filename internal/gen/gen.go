package gen

import (
	"errors"
	"io"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"github.com/lunashade/lang/internal/ast"
)

type Generator struct {
	m          *ir.Module
	funcStack  Stack[ir.Func]
	blockStack Stack[ir.Block]
}

func Run(w io.Writer, tree ast.AST) error {
	g := &Generator{
		m: ir.NewModule(),
	}
	if err := g.walk(tree); err != nil {
		return err
	}
	g.m.WriteTo(w)
	return nil
}

func (g *Generator) walk(node ast.AST) error {
	switch nd := node.(type) {
	case *ast.Root:
		for _, node := range nd.Nodes {
			err := g.walk(node)
			if err != nil {
				return err
			}
		}
		return nil
	case *ast.Function:
		name := nd.Name.(*ast.Ident)
		ty := types.I32
		fn := g.m.NewFunc(name.Name, ty)
		g.funcStack.Push(fn)

		blk := g.funcStack.Top().NewBlock("")
		g.blockStack.Push(blk)
		var val value.Value
		for _, node := range nd.Body {
			var err error
			stmt := node.(ast.Stmt)
			val, err = g.stmt(stmt)
			if err != nil {
				return err
			}
		}
		if val != nil {
			g.blockStack.Top().NewRet(val)
		} else {
			g.blockStack.Top().NewRet(constant.NewInt(ty, 0))
		}
		g.blockStack.Pop()
		g.funcStack.Pop()
	}
	return nil
}

func (g *Generator) stmt(node ast.Stmt) (value.Value, error) {
	switch nd := node.(type) {
	case *ast.ExprStmt:
		expr := nd.Expr.(ast.Expr)
		return g.expr(expr)
	case *ast.Semi:
		expr := nd.Expr.(ast.Expr)
		_, err := g.expr(expr)
		return nil, err
	default:
		return nil, errors.New("unknown statement")
	}
}

func (g *Generator) expr(node ast.Expr) (value.Value, error) {
	switch nd := node.(type) {
	case *ast.Int:
		return constant.NewInt(types.I32, nd.Value), nil
	case *ast.BinOp:
		return g.binOp(nd)
	case *ast.Block:
		var val value.Value
		var err error
		for _, n := range nd.Stmts {
			stmt := n.(ast.Stmt)
			val, err = g.stmt(stmt)
			if err != nil {
				return nil, err
			}
		}
		return val, nil
	default:
		return nil, errors.New("unknown expr")
	}
}

func (g *Generator) binOp(node *ast.BinOp) (value.Value, error) {
	// TODO: remove type assertion
	// LHS, RHS must be expr so solve this in parse section
	lhsNode := node.LHS.(ast.Expr)
	lhs, err := g.expr(lhsNode)
	if err != nil {
		return nil, err
	}
	rhsNode := node.RHS.(ast.Expr)
	rhs, err := g.expr(rhsNode)
	if err != nil {
		return nil, err
	}

	switch node.Kind {
	case ast.Add:
		res := g.blockStack.Top().NewAdd(lhs, rhs)
		return res, nil
	case ast.Sub:
		res := g.blockStack.Top().NewSub(lhs, rhs)
		return res, nil
	case ast.Mul:
		res := g.blockStack.Top().NewMul(lhs, rhs)
		return res, nil
	case ast.Div:
		res := g.blockStack.Top().NewSDiv(lhs, rhs)
		return res, nil
	case ast.Equal:
		b := g.blockStack.Top().NewICmp(enum.IPredEQ, lhs, rhs)
		res := g.blockStack.Top().NewZExt(b, types.I32)
		return res, nil
	case ast.NotEqual:
		b := g.blockStack.Top().NewICmp(enum.IPredNE, lhs, rhs)
		res := g.blockStack.Top().NewZExt(b, types.I32)
		return res, nil
	case ast.LessThan:
		b := g.blockStack.Top().NewICmp(enum.IPredSLT, lhs, rhs)
		res := g.blockStack.Top().NewZExt(b, types.I32)
		return res, nil
	case ast.GreaterThan:
		b := g.blockStack.Top().NewICmp(enum.IPredSGT, lhs, rhs)
		res := g.blockStack.Top().NewZExt(b, types.I32)
		return res, nil
	case ast.LessThanOrEqual:
		b := g.blockStack.Top().NewICmp(enum.IPredSLE, lhs, rhs)
		res := g.blockStack.Top().NewZExt(b, types.I32)
		return res, nil
	case ast.GreaterThanOrEqual:
		b := g.blockStack.Top().NewICmp(enum.IPredSGE, lhs, rhs)
		res := g.blockStack.Top().NewZExt(b, types.I32)
		return res, nil
	}
	return nil, errors.New("unknown operator")
}
