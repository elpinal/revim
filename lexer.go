package revim

import (
	"errors"
	"fmt"
	"unicode/utf8"
)

const eof = 0

type patternLex struct {
	line []byte
	peek rune
	err  error

	pattern frag

	off int // information for error messages
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
	x.off++
	if c == utf8.RuneError && size == 1 {
		x.err = errors.New("next: invalid utf8")
		return x.next()
	}
	return c
}

func (x *patternLex) Error(s string) {
	x.err = fmt.Errorf("parse error (offset: %d, peek: %q): %s", x.off, x.peek, s)
}

func parse(line []byte) (*state, error) {
	l := patternLex{line: line}
	yyParse(&l)
	if l.err != nil {
		return nil, l.err
	}
	f := l.pattern
	patch(f, *matchState)
	return f.start, nil
}
