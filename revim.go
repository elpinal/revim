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

func (re *Regexp) FindStringIndex(str string) (loc []int) {
	s := []rune(str)
	var i int
	for {
		ok, rest := re.s.process(s)
		if ok {
			return []int{i, len(str) - len(rest)}
		}
		if len(s) == 0 {
			return nil
		}
		s = s[1:]
		i++
	}
	return nil
}

func (re *Regexp) FindAllStringIndex(str string) [][]int {
	result := make([][]int, 0, 4)
	s := []rune(str)
	var i int
	for {
		ok, rest := re.s.process(s)
		if ok {
			result = append(result, []int{i, len(str) - len(rest)})
			if len(s) != len(rest) {
				s = rest
				i = len(str) - len(rest)
				continue
			}
		}
		if len(s) == 0 {
			if len(result) == 0 {
				return nil
			}
			return result
		}
		s = s[1:]
		i++
	}
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
