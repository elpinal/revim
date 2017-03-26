package revim

import "strings"

type Regexp struct {
	expr string
}

func Compile(expr string) *Regexp {
	re := parse(expr)
	return re
}

func parse(expr string) *Regexp {
	return &Regexp{
		expr: expr,
	}
}

func (re *Regexp) MatchString(s string) bool {
	return strings.Index(s, re.expr) >= 0
}
