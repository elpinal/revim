package revim

import (
	"log"
	"unicode/utf8"
)

const eof = 0

type patternLex struct {
	line []byte
	peek rune

	pattern frag
}

func (x *patternLex) Lex(yylval *yySymType) int {
	for {
		c := x.next()
		switch c {
		case eof:
			return eof
		case '\\':
			return x.escape(yylval)
		case '*':
			return int(c)
		default:
			yylval.char = c
			return CHAR
		}
	}
}

func (x *patternLex) escape(yylval *yySymType) int {
	c := x.next()
	switch c {
	case '|':
		return ALT
	case '&':
		return AND
	case '(':
		return LPAREN
	case ')':
		return RPAREN
	case '+':
		return PLUS
	case '?':
		return QUESTION
	case '\\':
		yylval.char = '\\'
		return CHAR
	}
	if c != eof {
		x.peek = c
	}
	return '\\'
}

func (x *patternLex) next() rune {
	if x.peek != eof {
		r := x.peek
		x.peek = eof
		return r
	}
	if len(x.line) == 0 {
		return eof
	}
	c, size := utf8.DecodeRune(x.line)
	x.line = x.line[size:]
	if c == utf8.RuneError && size == 1 {
		log.Print("invalid utf8")
		return x.next()
	}
	return c
}

func (x *patternLex) Error(s string) {
	log.Printf("parse error (peek: %d): %s", x.peek, s)
}

func parse(line []byte) *state {
	l := patternLex{line: line}
	yyParse(&l)
	f := l.pattern
	patch(f, *matchState)
	return f.start
}
