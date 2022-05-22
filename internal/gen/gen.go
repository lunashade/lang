package gen

import (
	"errors"
	"io"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"github.com/lunashade/lang/internal/ast"
)

type Generator struct {
	m        *ir.Module
	curFunc  *ir.Func
	curBlock *ir.Block
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
		var val value.Value
		var err error
		g.curFunc = g.m.NewFunc("main", types.I32)
		g.curBlock = g.curFunc.NewBlock("")
		for _, expr := range nd.Nodes {
			val, err = g.expr(expr)
			if err != nil {
				return err
			}
		}
		g.curBlock.NewRet(val)
	}
	return nil
}

func (g *Generator) expr(node ast.AST) (value.Value, error) {
	switch nd := node.(type) {
	case *ast.Int:
		return constant.NewInt(types.I32, nd.Value), nil
	case *ast.BinOp:
		return g.binOp(nd)
	default:
		return nil, errors.New("unknown expr")
	}
}

func (g *Generator) binOp(node *ast.BinOp) (value.Value, error) {
	lhs, err := g.expr(node.LHS)
	if err != nil {
		return nil, err
	}
	rhs, err := g.expr(node.RHS)
	if err != nil {
		return nil, err
	}

	switch node.Type {
	case ast.Add:
		res := g.curBlock.NewAdd(lhs, rhs)
		return res, nil
	case ast.Sub:
		res := g.curBlock.NewSub(lhs, rhs)
		return res, nil
	case ast.Mul:
		res := g.curBlock.NewMul(lhs, rhs)
		return res, nil
	case ast.Div:
		res := g.curBlock.NewSDiv(lhs, rhs)
		return res, nil
	}
	return nil, errors.New("unknown operator")
}
