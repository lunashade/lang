package token

import (
	"bufio"
	"errors"
	"io"
)

func Lex(r io.Reader) chan Token {
	l := &lexer{
		src: bufio.NewReader(r),
		buf: make([]rune, 0, 8),
		ch:  make(chan Token),
	}

	go l.run()
	return l.ch
}

type stateFn func(*lexer) stateFn
type lexer struct {
	src    *bufio.Reader
	peeked rune
	buf    []rune
	ch     chan Token
}

const eof rune = -1

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

func (l *lexer) emit(kind Kind) {
	tok := makeToken(kind, string(l.buf))
	l.ch <- tok
	l.buf = nil
}

func (l *lexer) run() {
	for state := lexSkip; state != nil; {
		state = state(l)
	}
	close(l.ch)
}

func lexSkip(l *lexer) stateFn {
	for {
		c := l.peek()
		if isDigit(c) {
			return lexNumber
		}
		if l.next() == eof {
			break
		}
	}
	l.emit(Eof)
	return nil
}

func lexNumber(l *lexer) stateFn {
	for {
		c := l.next()
		if !isDigit(c) {
			l.backup()
			break
		}
		l.buf = append(l.buf, c)
	}
	l.emit(Integer)
	return lexSkip
}

func isDigit(c rune) bool {
	return '0' <= c && c <= '9'
}
