package revim

//go:generate goyacc -o parse.go parse.y

import "strings"

type reRange struct {
	left, right int
}

type Re interface {
	match(string) *reRange
}

type lit string

func (l lit) match(s string) *reRange {
	i := strings.Index(s, string(l))
	if i < 0 {
		return nil
	}
	return &reRange{
		left:  i,
		right: i + len(s),
	}
}

func literal(s string) Re {
	return lit(s)
}

type alt struct {
	re1, re2 Re
}

func (a alt) match(s string) *reRange {
	rr := a.re1.match(s)
	if rr != nil {
		return rr
	}
	return a.re2.match(s)
}

func pattern(re1, re2 Re) Re {
	return alt{
		re1: re1,
		re2: re2,
	}
}
