package revim

import "unicode/utf8"

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
	for {
		ok, _ := re.s.process(make(map[*state]string), s)
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

func (s *state) process(m map[*state]string, str string) (bool, string) {
	if s.match {
		return true, str
	}
	if s.backtrack != nil {
		m[s.backtrack] = str
	}
	if bs, ok := m[s]; ok {
		str = bs
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
			return s.out.process(m, str[size:])
		}
		return false, str
	}
	ok, rest := s.out.process(m, str)
	if ok {
		return true, rest
	}
	if s.out1 == nil {
		return false, str
	}
	ok, rest = s.out1.process(m, str)
	if ok {
		return true, rest
	}
	return false, str
}
