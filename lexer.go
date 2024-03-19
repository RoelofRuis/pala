package pala

import (
	"io"
	"unicode"
)

type Lexer interface {
	nextToken() token
}

type tokenType int

const (
	tokenLiteral tokenType = iota
	tokenNewline
	tokenVariable
	tokenLBracket
	tokenRBracket
	tokenLParen
	tokenRParen
	tokenComment
	tokenEOF
	tokenInvalid
)

type token struct {
	tpe   tokenType
	line  int
	value string
}

type basicLexer struct {
	scanner  io.RuneScanner
	next     func(l *basicLexer) token
	currCh   rune
	currLine int
}

func NewLexer(scanner io.RuneScanner) Lexer {
	lexer := &basicLexer{scanner: scanner, next: readLine}
	lexer.readChar()
	return lexer
}

func (l *basicLexer) nextToken() token {
	l.skipWhitespace()

	return l.next(l)
}

func (l *basicLexer) readChar() {
	ch, _, err := l.scanner.ReadRune()
	if err != nil {
		l.currCh = 0
	} else {
		l.currCh = ch
	}
	if l.currCh == '\n' {
		l.currLine++
	}
}

func (l *basicLexer) skipWhitespace() {
	for unicode.IsSpace(l.currCh) && !isLineEnd(l.currCh) {
		l.readChar()
	}
}

func (l *basicLexer) scanLine() string {
	var result []rune
	for !isLineEnd(l.currCh) {
		result = append(result, l.currCh)
		l.readChar()
	}
	return string(result)
}

func (l *basicLexer) scanWord() string {
	var result []rune
	for unicode.IsGraphic(l.currCh) && !unicode.IsSpace(l.currCh) && l.currCh != '[' && l.currCh != ']' && !isLineEnd(l.currCh) {
		result = append(result, l.currCh)
		l.readChar()
	}
	return string(result)
}

func readLine(l *basicLexer) token {
	switch {
	case l.currCh == '(':
		l.readChar()
		return l.makeToken(tokenLParen, "(")
	case l.currCh == ')':
		l.readChar()
		return l.makeToken(tokenRParen, ")")
	case l.currCh == '[':
		l.readChar()
		return l.makeToken(tokenLBracket, "[")
	case l.currCh == ']':
		l.readChar()
		return l.makeToken(tokenRBracket, "]")
	case l.currCh == '#':
		return l.makeToken(tokenComment, l.scanLine())
	case l.currCh == '\n':
		l.readChar()
		return l.makeToken(tokenNewline, "\n")
	case l.currCh == '$':
		return l.makeToken(tokenVariable, l.scanWord())
	case unicode.IsGraphic(l.currCh):
		return l.makeToken(tokenLiteral, l.scanWord())
	case l.currCh == 0:
		return l.makeToken(tokenEOF, "")
	default:
	}

	l.readChar()
	return l.makeToken(tokenInvalid, string(l.currCh))
}

func (l *basicLexer) makeToken(tpe tokenType, Value string) token {
	return token{tpe: tpe, value: Value, line: l.currLine}
}

func isLineEnd(c rune) bool {
	return c == '\n' || c == 0
}
