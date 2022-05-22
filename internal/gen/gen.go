package gen

import (
	"errors"
	"fmt"
	"io"

	"github.com/lunashade/lang/internal/ast"
)

type Reg int

const InvalidReg Reg = -1

type Generator struct {
	w           io.Writer
	reg         Reg
	indentDepth int
}

func (g *Generator) nextReg() Reg {
	g.reg++
	return g.reg
}

func Run(w io.Writer, tree ast.AST) error {
	g := &Generator{w: w}
	root, ok := tree.(*ast.Root)
	if !ok {
		return errors.New("must be root")
	}
	return g.Root(root)
}

func (g *Generator) Root(root *ast.Root) error {
	var reg Reg
	var err error
	g.emitf("define dso_local i32 @main() #0 {")
	g.indentDepth++
	for _, expr := range root.Nodes {
		if reg, err = g.Expr(expr); err != nil {
			return err
		}
	}
	g.emitf("ret i32 %%%d", reg)
	g.indentDepth--
	g.emitf("}")
	return nil
}

func (g *Generator) Expr(expr ast.AST) (Reg, error) {
	switch nd := expr.(type) {
	case *ast.BinOp:
		return g.BinOp(nd)
	case *ast.Int:
		return g.Int(nd)
	default:
		return InvalidReg, errors.New("unknown expression")
	}

}

func (g *Generator) BinOp(op *ast.BinOp) (Reg, error) {
	lhs, err := g.Expr(op.LHS)
	if err != nil {
		return InvalidReg, err
	}
	rhs, err := g.Expr(op.RHS)
	if err != nil {
		return InvalidReg, err
	}

	res := g.nextReg()
	switch op.Type {
	case ast.Add:
		g.emitf("%%%d = add i32 %%%d, %%%d", res, lhs, rhs)
		return res, nil
	case ast.Sub:
		g.emitf("%%%d = sub i32 %%%d, %%%d", res, lhs, rhs)
		return res, nil
	case ast.Mul:
		g.emitf("%%%d = mul i32 %%%d, %%%d", res, lhs, rhs)
		return res, nil
	case ast.Div:
		g.emitf("%%%d = sdiv i32 %%%d, %%%d", res, lhs, rhs)
		return res, nil
	}
	return InvalidReg, errors.New("unknown operator")
}

func (g *Generator) Int(i *ast.Int) (Reg, error) {
	alloc := g.nextReg()
	result := g.nextReg()
	g.emitf("%%%d = alloca i32, align 4", alloc)
	g.emitf("store i32 %d, i32* %%%d", i.Value, alloc)
	g.emitf("%%%d = load i32, i32* %%%d, align 4", result, alloc)
	return result, nil
}

func (g *Generator) emitf(format string, a ...interface{}) error {
	for i := 0; i < g.indentDepth; i++ {
		fmt.Fprintf(g.w, "\t")
	}
	_, err := fmt.Fprintf(g.w, format, a...)
	fmt.Fprintf(g.w, "\n")
	return err
}
