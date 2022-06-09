package parse

import (
	"errors"
	"reflect"

	"github.com/lunashade/lang/internal/ast"
	"github.com/lunashade/lang/internal/token/kind"
)

type NonTerminal func(int) (int, ast.AST, error)
type Cache map[Key]*Ret
type Key struct {
	F   uintptr
	Pos int
}
type Ret struct {
	Pos int
	Nd  ast.AST
}

func (p *Parser) CachedCall(f NonTerminal, pos int) (int, ast.AST, error) {
	ptr := uintptr(reflect.ValueOf(f).Pointer())
	key := Key{F: ptr, Pos: pos}
	if ret, ok := p.cache[key]; ok {
		return ret.Pos, ret.Nd, nil
	}
	nx, node, err := f(pos)
	if err != nil {
		return pos, nil, err
	}
	ret := Ret{
		Pos: nx,
		Nd:  node,
	}
	p.cache[key] = &ret

	return nx, node, nil
}

func (p *Parser) Skip(kind kind.Kind) NonTerminal {
	return func(pos int) (int, ast.AST, error) {
		nx, t := p.consume(kind, pos)
		if t == nil {
			return pos, nil, errors.New("invalid token")
		}
		return nx, nil, nil
	}
}

func (p *Parser) Select(cands ...NonTerminal) NonTerminal {
	return func(pos int) (int, ast.AST, error) {
		var nx int
		var node ast.AST
		var err error
		for _, cand := range cands {
			nx, node, err = p.CachedCall(cand, pos)
			if err == nil {
				return nx, node, nil
			}
		}
		// TODO: wraps error
		return pos, nil, err
	}
}

func (p *Parser) Optional(cand NonTerminal) NonTerminal {
	return func(pos int) (int, ast.AST, error) {
		nx, node, err := cand(pos)
		if err != nil {
			return pos, nil, nil
		}
		return nx, node, nil
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

// Repeat is parser combinator of "cand*"
func (p *Parser) Repeat(m Merger, cand NonTerminal) NonTerminal {
	return func(pos int) (int, ast.AST, error) {
		var nx int = pos
		var node ast.AST
		var err error

		nodes := make([]ast.AST, 0)
		for {
			nx, node, err = p.CachedCall(cand, nx)
			if err != nil {
				break
			}
			nodes = append(nodes, node)
		}
		return nx, m(nodes), nil
	}
}

// RepeatWithOptionalLast is parser combinator of "cand* last?"
func (p *Parser) RepeatWithOptionalLast(m Merger, cand NonTerminal, last NonTerminal) NonTerminal {
	return func(pos int) (int, ast.AST, error) {
		var nx int = pos
		var node ast.AST
		var err error

		nodes := make([]ast.AST, 0)
		for {
			nx, node, err = p.CachedCall(cand, nx)
			if err == nil {
				nodes = append(nodes, node)
				continue
			}
			// if cand fails, try last parser and break anyway
			nx, node, err = p.CachedCall(last, nx)
			if err == nil {
				nodes = append(nodes, node)
			}
			break
		}
		return nx, m(nodes), nil
	}
}
