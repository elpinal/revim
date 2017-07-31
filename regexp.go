package revim

//go:generate goyacc -o parse.go parse.y

import (
	"strings"
	"unicode/utf8"
)

type reRange struct {
	left, right int
}

type Re interface {
	match(string) *reRange
}

type lit rune

func (l lit) match(s string) *reRange {
	i := strings.IndexRune(s, rune(l))
	if i < 0 {
		return nil
	}
	return &reRange{
		left:  i,
		right: i + utf8.RuneLen(rune(l)),
	}
}

func literal(r rune) Re {
	return lit(r)
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

type and struct {
	re1, re2 Re
}

func (a and) match(s string) *reRange {
	rr1 := a.re1.match(s)
	if rr1 == nil {
		return nil
	}
	rr2 := a.re2.match(s)
	if rr2 == nil {
		return nil
	}
	if rr1.left != rr2.left {
		return nil
	}
	return rr2
}

func branch(re1, re2 Re) Re {
	return and{
		re1: re1,
		re2: re2,
	}
}

type con struct {
	re1, re2 Re
}

func (c con) match(s string) *reRange {
	rr1 := c.re1.match(s)
	if rr1 == nil {
		return nil
	}
	rr2 := c.re2.match(s[rr1.right:])
	if rr2 == nil {
		return nil
	}
	if rr2.left != 0 {
		return nil
	}
	return &reRange{
		left:  rr1.left,
		right: rr1.right + rr2.right,
	}
}

func concat(re1, re2 Re) Re {
	return con{
		re1: re1,
		re2: re2,
	}
}

type mul struct {
	re Re
}

func (m mul) match(s string) *reRange {
	rr := m.re.match(s)
	if rr == nil {
		return &reRange{
			left:  0,
			right: 0,
		}
	}
	off := rr.right
	left := rr.left
	right := rr.right
	for {
		rr := m.re.match(s[off:])
		if rr == nil {
			return &reRange{
				left:  left,
				right: right,
			}
		}
		off += rr.right
		right = off
	}
}

func multi(re Re) Re {
	return mul{
		re: re,
	}
}

type pl struct {
	re Re
}

func (p pl) match(s string) *reRange {
	rr := p.re.match(s)
	if rr == nil {
		return nil
	}
	m := mul{re: p.re}
	rr0 := m.match(s[rr.right:])
	if rr0 == nil {
		return rr
	}
	return &reRange{
		left:  rr.left,
		right: rr.right + rr0.right,
	}
}

func plus(re Re) Re {
	return pl{
		re: re,
	}
}
