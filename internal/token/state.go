package token

import "github.com/lunashade/lang/internal/token/kind"

type stateFn func(*lexer) stateFn

func lexSkip(l *lexer) stateFn {
	for {
		c := l.peek()
		if isDigit(c) {
			return lexNumber
		}
		if isSymbol(c) {
			return lexSymbol
		}
		if isIdent(c) {
			return lexIdent
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
func lexIdent(l *lexer) stateFn {
	for {
		c := l.next()
		if !isIdent(c) {
			break
		}
		l.buf = append(l.buf, c)
	}
	l.backup()
	l.emit(kind.Identifier)
	return lexSkip
}

// lexSymbol consume punctuation symbol with single character
func lexSymbol(l *lexer) stateFn {
	c := l.next()
	l.buf = append(l.buf, c)
	k := kind.SymbolKind(c)
	if k == kind.Invalid {
		panic("unknown symbol")
	}
	l.emit(k)
	return lexSkip
}
