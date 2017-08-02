%{

package revim

%}

%union {
        char rune
        tok token
        frag frag
}

%type	<frag>	pattern branch concat piece atom

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
