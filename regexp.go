package revim

import "log"

//go:generate goyacc -o parse.go parse.y

type state struct {
	split bool
	match bool

	r    rune
	out  *state
	out1 *state
	//lastList int

	backtrack *state
}

var matchState = &state{match: true, out: &state{}}

type frag struct {
	start *state
	out   []*state
}

func literal(r rune) frag {
	s := state{r: r, out: &state{}}
	return frag{
		start: &s,
		out:   []*state{(&s).out},
	}
}

func pattern(f1, f2 frag) frag {
	s := &state{
		split: true,
		out:   f1.start,
		out1:  f2.start,
	}
	return frag{
		start: s,
		out:   append(f1.out, f2.out...),
	}
}

func patch(f1 frag, s state) {
	for i := range f1.out {
		*f1.out[i] = s
	}
}

func branch(f1, f2 frag) frag {
	s := state{
		split: true, // FIXME: In fact, it's false.
		out:   f2.start,
	}
	patch(f1, s)
	f1.start.backtrack = s.out
	return frag{
		start: f1.start,
		out:   f2.out,
	}
}

func concat(f1, f2 frag) frag {
	for i := range f1.out {
		log.Println("concat", f1, f2)
		*f1.out[i] = *f2.start
	}
	return frag{
		start: f1.start,
		out:   f2.out,
	}
}

func multi(f frag) frag {
	s := state{
		split: true,
		out:   f.start,
		out1:  &state{},
	}
	for i := range f.out {
		*f.out[i] = s
	}
	return frag{
		start: &s,
		out:   []*state{(&s).out1},
	}
}

func plus(f frag) frag {
	s := state{
		split: true,
		out:   f.start,
		out1:  &state{},
	}
	for i := range f.out {
		*f.out[i] = s
	}
	return frag{
		start: f.start,
		out:   []*state{(&s).out1},
	}
}

func question(f frag) frag {
	s := state{
		split: true,
		out:   f.start,
		out1:  &state{},
	}
	return frag{
		start: &s,
		out:   append(f.out, s.out1),
	}
}
