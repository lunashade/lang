package parse

import (
	"fmt"
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
		{
			"1<=1",
			&ast.BinOp{
				Kind: ast.LessThanOrEqual,
				LHS:  &ast.Int{Value: 1},
				RHS:  &ast.Int{Value: 1},
			},
		},
		{
			"1>=1",
			&ast.BinOp{
				Kind: ast.GreaterThanOrEqual,
				LHS:  &ast.Int{Value: 1},
				RHS:  &ast.Int{Value: 1},
			},
		},
		{
			"1<1",
			&ast.BinOp{
				Kind: ast.LessThan,
				LHS:  &ast.Int{Value: 1},
				RHS:  &ast.Int{Value: 1},
			},
		},
		{
			"1>1",
			&ast.BinOp{
				Kind: ast.GreaterThan,
				LHS:  &ast.Int{Value: 1},
				RHS:  &ast.Int{Value: 1},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			input := fmt.Sprintf("main(){%s}", tt.input)
			ch := token.Lex(strings.NewReader(input))
			node, _ := Run(ch)
			root := node.(*ast.Root)
			fn := root.Nodes[0].(*ast.Function)
			body := fn.Body[0].(*ast.ExprStmt)
			got := body.Expr
			assert.DeepEqual(t, tt.want, got)
		})
	}
}

func TestParseFunc(t *testing.T) {
	tests := []struct {
		input string
		want  ast.AST
	}{
		{
			"main(){1; 5*5}",
			&ast.Function{
				Name: &ast.Ident{
					Name: "main",
				},
				Body: []ast.AST{
					&ast.Semi{
						Expr: &ast.Int{Value: 1},
					},
					&ast.ExprStmt{
						Expr: &ast.BinOp{
							Kind: ast.Mul,
							LHS:  &ast.Int{Value: 5},
							RHS:  &ast.Int{Value: 5},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			input := fmt.Sprintf("%s", tt.input)
			ch := token.Lex(strings.NewReader(input))
			node, _ := Run(ch)
			root := node.(*ast.Root)
			got := root.Nodes[0].(*ast.Function)
			assert.DeepEqual(t, tt.want, got)
		})
	}

}
