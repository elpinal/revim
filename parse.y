%{

package revim

import (
	"log"
	"unicode/utf8"
)

%}

%union {
        char rune
        re Re
        tok token
}

%type	<re>	pattern branch concat piece atom

%token  <tok>   ALT AND LPAREN RPAREN PLUS QUESTION

%token	<char>	CHAR

%%

pattern:
	branch
        {
                $$ = $1
                if l, ok := yylex.(*patternLex); ok {
                        l.pattern = $$
                }
        }
|	pattern ALT branch
	{
                $$ = pattern($1, $3)
                if l, ok := yylex.(*patternLex); ok {
                        l.pattern = $$
                }
	}

branch:
	concat
|	branch AND concat
        {
                $$ = branch($1, $3)
        }

concat:
        piece
|	concat piece
        {
                $$ = concat($1, $2)
        }

piece:
        atom
|	atom '*'
        {
                $$ = multi($1)
        }
|	atom PLUS
        {
                $$ = plus($1)
        }
|	atom QUESTION
        {
                $$ = question($1)
        }

atom:
	CHAR
	{
		$$ = literal($1)
	}
|	LPAREN pattern RPAREN
        {
                $$ = $2
        }


%%

const eof = 0

type patternLex struct {
	line []byte
	peek rune

        pattern Re
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

func parse(line []byte) Re {
        l := patternLex{line: line}
        yyParse(&l)
        return l.pattern
}
