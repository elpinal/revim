package revim

import "testing"

func TestMatchString(t *testing.T) {
	for _, test := range []struct {
		expr string
		s    string
		ok   bool
	}{
		{
			"aaa", "aaa", true,
		},
		{
			"aaa", " aaaaa bbbbb ", true,
		},
		{
			"aaa", "aa", false,
		},
		{
			`aaa\|bbb`, "aaaa", true,
		},
		{
			`aaa\|bbb`, "bbbb", true,
		},
		{
			`aaa\|bbb`, "bb", false,
		},
		{
			`aaabbb\&aaa`, "aaabbb", true,
		},
		{
			`aaabbb\&aaa`, "aaa", false,
		},
		{
			`bbb\&aaa\|aaabbb\&aaa`, "aaabbb", true,
		},
		{
			`aaa\\bbb`, `aaa\bbb`, true,
		},
	} {
		re := Compile(test.expr)
		ok := re.MatchString(test.s)
		if ok != test.ok {
			t.Errorf("MatchString(%q) should be %v", test.s, test.ok)
		}
	}
}
