package parse

import (
	"strings"
	"testing"

	"gotest.tools/assert"

	"github.com/lunashade/lang/internal/ast"
	"github.com/lunashade/lang/internal/token"
)

func TestParseExpr(t *testing.T) {
	tests := []struct {
		input string
		want  ast.AST
	}{
		{
			"1+1",
			&ast.BinOp{
				Kind: ast.Add,
				LHS:  &ast.Int{Value: 1},
				RHS:  &ast.Int{Value: 1},
			},
		},
		{
			"1-1",
			&ast.BinOp{
				Kind: ast.Sub,
				LHS:  &ast.Int{Value: 1},
				RHS:  &ast.Int{Value: 1},
			},
		},
		{
			"1*1",
			&ast.BinOp{
				Kind: ast.Mul,
				LHS:  &ast.Int{Value: 1},
				RHS:  &ast.Int{Value: 1},
			},
		},
		{
			"1/1",
			&ast.BinOp{
				Kind: ast.Div,
				LHS:  &ast.Int{Value: 1},
				RHS:  &ast.Int{Value: 1},
			},
		},
		{
			"1+1*1",
			&ast.BinOp{
				Kind: ast.Add,
				LHS:  &ast.Int{Value: 1},
				RHS: &ast.BinOp{
					Kind: ast.Mul,
					LHS:  &ast.Int{Value: 1},
					RHS:  &ast.Int{Value: 1},
				},
			},
		},
		{
			"1+1+1",
			&ast.BinOp{
				Kind: ast.Add,
				LHS:  &ast.Int{Value: 1},
				RHS: &ast.BinOp{
					Kind: ast.Add,
					LHS:  &ast.Int{Value: 1},
					RHS:  &ast.Int{Value: 1},
				},
			},
		},
		{
			"2*(3+4)",
			&ast.BinOp{
				Kind: ast.Mul,
				LHS:  &ast.Int{Value: 2},
				RHS: &ast.BinOp{
					Kind: ast.Add,
					LHS:  &ast.Int{Value: 3},
					RHS:  &ast.Int{Value: 4},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			ch := token.Lex(strings.NewReader(tt.input))
			node, _ := Run(ch)
			root := node.(*ast.Root)
			got := root.Nodes[0]
			assert.DeepEqual(t, tt.want, got)
		})
	}

}
