package revim

type Regexp struct {
	expr string
	s    *state
}

func Compile(expr string) (*Regexp, error) {
	s, err := parse([]byte(expr))
	if err != nil {
		return nil, err
	}
	return &Regexp{
		expr: expr,
		s:    s,
	}, nil
}

func (re *Regexp) MatchString(str string) bool {
	s := []rune(str)
	for {
		ok, _ := re.s.process(s)
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

func (s *state) process(str []rune) (bool, []rune) {
	if s.match {
		return true, str
	}
	if s.and {
		ok, rest := s.out.process(str)
		if !ok {
			return false, str
		}
		ok, rest = s.out1.process(str)
		if ok {
			return true, rest
		}
		return false, str
	}
	if !s.split {
		if len(str) == 0 {
			return false, str
		}
		r := str[0]
		if s.r == r {
			if s.out == nil {
				return true, str[1:]
			}
			return s.out.process(str[1:])
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
