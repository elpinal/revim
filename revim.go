package revim

import (
	"log"
	"unicode/utf8"
)

type Regexp struct {
	expr string
	s    *state
}

func Compile(expr string) *Regexp {
	s := parse([]byte(expr))
	return &Regexp{
		expr: expr,
		s:    s,
	}
}

func (re *Regexp) MatchString(str string) bool {
	s := str
	log.Println("MatchString", s)
	for {
		ok, rest := re.s.process(s)
		log.Println("MatchString", rest)
		if ok {
			return true
		}
		if len(s) == 0 {
			return false
		}
		s = s[1:]
	}
	return false
}

func (s *state) process(str string) (bool, string) {
	log.Printf("process: %#v %s", *s, str)
	if s.match {
		return true, str
	}
	if !s.split {
		if len(str) == 0 {
			return false, str
		}
		r, size := utf8.DecodeRuneInString(str)
		if s.r == r {
			if s.out == nil {
				return true, str[size:]
			}
			return s.out.process(str[size:])
		}
		return false, str
	}
	ok, rest := s.out.process(str)
	if ok {
		return true, rest
	}
	ok, rest = s.out1.process(str)
	if ok {
		return true, rest
	}
	return false, str
}
