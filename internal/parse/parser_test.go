package parse

import (
	"strings"
	"testing"

	"gotest.tools/assert"

	"github.com/lunashade/lang/internal/ast"
	"github.com/lunashade/lang/internal/token"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  ast.AST
	}{
		{
			"1+1", "1+1",
			&ast.BinOp{
				Type: ast.Add,
				LHS:  &ast.Int{Value: 1},
				RHS:  &ast.Int{Value: 1},
			},
		},
		{
			"1-1", "1-1",
			&ast.BinOp{
				Type: ast.Sub,
				LHS:  &ast.Int{Value: 1},
				RHS:  &ast.Int{Value: 1},
			},
		},
		{
			"1*1", "1*1",
			&ast.BinOp{
				Type: ast.Mul,
				LHS:  &ast.Int{Value: 1},
				RHS:  &ast.Int{Value: 1},
			},
		},
		{
			"1/1", "1/1",
			&ast.BinOp{
				Type: ast.Div,
				LHS:  &ast.Int{Value: 1},
				RHS:  &ast.Int{Value: 1},
			},
		},
		{
			"1+1*1", "1+1*1",
			&ast.BinOp{
				Type: ast.Add,
				LHS:  &ast.Int{Value: 1},
				RHS: &ast.BinOp{
					Type: ast.Mul,
					LHS:  &ast.Int{Value: 1},
					RHS:  &ast.Int{Value: 1},
				},
			},
		},
		{
			"1+1+1", "1+1+1",
			&ast.BinOp{
				Type: ast.Add,
				LHS:  &ast.Int{Value: 1},
				RHS: &ast.BinOp{
					Type: ast.Add,
					LHS:  &ast.Int{Value: 1},
					RHS:  &ast.Int{Value: 1},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ch := token.Lex(strings.NewReader(tt.input))
			got, _ := Run(ch)
			assert.DeepEqual(t, tt.want, got)
		})
	}

}
