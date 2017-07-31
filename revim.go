package revim

type Regexp struct {
	expr string
	Re
}

func Compile(expr string) *Regexp {
	re := parse([]byte(expr))
	return &Regexp{
		expr: expr,
		Re:   re,
	}
}

func (re *Regexp) MatchString(s string) bool {
	rr := re.match(s)
	return rr != nil
}
