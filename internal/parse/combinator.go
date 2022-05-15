package parse

import "github.com/lunashade/lang/internal/ast"

type NonTerminal func(int) (int, ast.AST, error)

func (p *Parser) Select(cands ...NonTerminal) NonTerminal {
	return func(pos int) (int, ast.AST, error) {
		var nx int
		var node ast.AST
		var err error
		for _, cand := range cands {
			nx, node, err = cand(pos)
			if err == nil {
				return nx, node, nil
			}
		}
		// TODO: wraps error
		return pos, nil, err
	}
}

type Merger func([]ast.AST) ast.AST

func (p *Parser) Concat(m Merger, cands ...NonTerminal) NonTerminal {
	return func(pos int) (int, ast.AST, error) {
		var nx int = pos
		var node ast.AST
		var err error

		nodes := make([]ast.AST, 0)
		for _, cand := range cands {
			nx, node, err = cand(nx)
			if err != nil {
				return pos, nil, err
			}
			nodes = append(nodes, node)
		}
		return nx, m(nodes), nil
	}
}
