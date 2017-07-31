%{

package revim

import (
	"bytes"
	"log"
	"unicode/utf8"
)

%}

%union {
        str string
        re Re
}

%type	<re>	pattern branch

%token '\\' '|'

%token	<str>	STRING

%%

pattern:
	branch
        {
                $$ = $1
                if l, ok := yylex.(*patternLex); ok {
                        l.pattern = $$
                }
        }
|	pattern '\\' '|' branch
	{
                $$ = pattern($1, $4)
                if l, ok := yylex.(*patternLex); ok {
                        l.pattern = $$
                }
	}

branch:
	STRING
	{
		$$ = literal($1)
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
                case '\\', '|':
                        return int(c)
		default:
			return x.str(c, yylval)
		}
	}
}

func (x *patternLex) str(c rune, yylval *yySymType) int {
	add := func(b *bytes.Buffer, c rune) {
		if _, err := b.WriteRune(c); err != nil {
			log.Fatalf("WriteRune: %s", err)
		}
	}
	var b bytes.Buffer
	add(&b, c)
	L: for {
		c = x.next()
		switch c {
		case eof, '\\', '|':
			break L
		default:
			add(&b, c)
		}
	}
	if c != eof {
		x.peek = c
	}
	yylval.str = b.String()
	return STRING
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
