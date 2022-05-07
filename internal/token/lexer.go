package token

import (
	"bufio"
	"errors"
	"io"

	"github.com/lunashade/lang/internal/token/kind"
)

const eof rune = -1

type lexer struct {
	src    *bufio.Reader
	peeked rune
	buf    []rune
	ch     chan Token
}

func Lex(r io.Reader) chan Token {
	l := &lexer{
		src: bufio.NewReader(r),
		buf: make([]rune, 0, 8),
		ch:  make(chan Token),
	}

	go l.run()
	return l.ch
}

func (l *lexer) run() {
	for state := lexSkip; state != nil; {
		state = state(l)
	}
	close(l.ch)
}

func (l *lexer) next() rune {
	c, _, err := l.src.ReadRune()
	l.peeked = c
	if err != nil {
		if errors.Is(err, io.EOF) {
			l.peeked = eof
			return eof
		}
		// TODO: handle err
		panic(err)
	}
	return c
}
func (l *lexer) backup() {
	if l.peeked == eof {
		return
	}
	if err := l.src.UnreadRune(); err != nil {
		// TODO: handle err
		panic(err)
	}
}

func (l *lexer) peek() rune {
	c := l.next()
	l.backup()
	return c
}

func (l *lexer) emit(kind kind.Kind) {
	tok := makeToken(kind, string(l.buf))
	l.ch <- tok
	l.buf = nil
}
