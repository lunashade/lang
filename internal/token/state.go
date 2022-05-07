package token

import "github.com/lunashade/lang/internal/token/kind"

type stateFn func(*lexer) stateFn

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
	l.emit(kind.Eof)
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
	l.emit(kind.Integer)
	return lexSkip
}
