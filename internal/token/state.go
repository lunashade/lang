package token

import "github.com/lunashade/lang/internal/token/kind"

type stateFn func(*lexer) stateFn

func lexSkip(l *lexer) stateFn {
	for {
		c := l.peek()
		if isDigit(c) {
			return lexNumber
		}
		if isPunctuation(c) {
			return lexOp
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
			break
		}
		l.buf = append(l.buf, c)
	}
	l.backup()
	l.emit(kind.Integer)
	return lexSkip
}

// lexOp consume punctuation symbol with single character
// TODO: support multi character
func lexOp(l *lexer) stateFn {
	c := l.next()
	l.buf = append(l.buf, c)
	l.emit(kind.Punctuation)
	return lexSkip
}
