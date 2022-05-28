package token

import (
	"errors"

	"github.com/lunashade/lang/internal/token/kind"
)

type Stream struct {
	ch     chan Token
	tokens []Token
	ateof  bool
}

func NewStream(ch chan Token) *Stream { return &Stream{ch: ch} }

var ErrAtEof = errors.New("at eof")

// fetch reads next token from channel
func (p *Stream) fetch() error {
	if p.ateof {
		return ErrAtEof
	}
	tok := <-p.ch
	p.tokens = append(p.tokens, tok)
	if tok.Kind == kind.Eof {
		p.ateof = true
		return ErrAtEof
	}
	return nil
}

// look at the token at the position
func (p *Stream) Look(at int) *Token {
	if at < len(p.tokens) {
		return &(p.tokens[at])
	}
	for {
		err := p.fetch()
		if err != nil {
			break
		}
		if at < len(p.tokens) {
			return &(p.tokens[at])
		}
	}
	return nil
}

func (p *Stream) Complete() bool {
	p.fetch()
	return p.ateof
}
