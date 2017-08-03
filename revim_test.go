package revim

import "testing"

func TestMatchString(t *testing.T) {
	for i, test := range []struct {
		expr string
		s    string
		ok   bool
	}{
		{
			`aaa`, `aaa`, true,
		},
		{
			`aaa`, ` aaaaa bbbbb `, true,
		},
		{
			`aaa`, `aa`, false,
		},
		{
			`aaa\|bbb`, `aaaa`, true,
		},
		{
			`aaa\|bbb`, `bbbb`, true,
		},
		{
			`aaa\|bbb`, `bb`, false,
		},
		{
			`aaabbb\&aaa`, `aaabbb`, true,
		},
		{
			`aaabbb\&aaa`, `aaa`, false,
		},
		{
			`bbb\&aaa\|aaabbb\&aaa`, `aaabbb`, true,
		},
		{
			`aaa\\bbb`, `aaa\bbb`, true,
		},
		{
			`aba*`, `ab`, true,
		},
		{
			`aba*`, `aba`, true,
		},
		{
			`aba*`, `abaa`, true,
		},
		{
			`aba*`, `aaa`, false,
		},
		{
			`\(aba\)*`, ``, true,
		},
		{
			`\(aba\)*`, `aba`, true,
		},
		{
			`aba\(aba\)*`, `abaaba`, true,
		},
		{
			`aba\(aba\)*`, `abaaaa`, true,
		},
		{
			`\(aba\)*`, `abbaabc`, true,
		},
		{
			`\(aba\)\+`, `abaaabc`, true,
		},
		{
			`\(aba\)\+`, `abaabac`, true,
		},
		{
			`\(aba\)\+`, `aaaabc`, false,
		},
		{
			`\(aba\)\?`, ``, true,
		},
		{
			`\(aba\)\?`, `aba`, true,
		},
		{
			`\(aba\)\?`, `ccc`, true,
		},
	} {
		re, err := Compile(test.expr)
		if err != nil {
			t.Fatalf("compiling (%d, %q): %v", i, test.s, err)
		}
		ok := re.MatchString(test.s)
		if ok != test.ok {
			t.Errorf("MatchString(%d, %q) should be %v", i, test.s, test.ok)
		}
	}
}

func TestFindStringIndex(t *testing.T) {
	for i, test := range []struct {
		expr string
		s    string
		loc  []int
	}{
		{
			`aaa`, `aaa`, []int{0, 3},
		},
		{
			`aaa`, ` aaaaa bbbbb `, []int{1, 4},
		},
		{
			`aaa`, `aa`, nil,
		},
		{
			`aaa\|bbb`, `aaaa`, []int{0, 3},
		},
		{
			`aaa\|bbb`, `bbbb`, []int{0, 3},
		},
		{
			`aaa\|bbb`, `bb`, nil,
		},
		{
			`aaabbb\&aaa`, `aaabbb`, []int{0, 3},
		},
		{
			`aaabbb\&aaa`, `aaa`, nil,
		},
		{
			`bbb\&aaa\|aaabbb\&aaa`, `aaabbb`, []int{0, 3},
		},
		{
			`aaa\\bbb`, `aaa\bbb`, []int{0, 7},
		},
		{
			`aba*`, `ab`, []int{0, 2},
		},
		{
			`aba*`, `aba`, []int{0, 3},
		},
		{
			`aba*`, `abaa`, []int{0, 4},
		},
		{
			`aba*`, `aaa`, nil,
		},
		{
			`\(aba\)*`, ``, []int{0, 0},
		},
		{
			`\(aba\)*`, `aba`, []int{0, 3},
		},
		{
			`aba\(aba\)*`, `abaaba`, []int{0, 6},
		},
		{
			`aba\(aba\)*`, `abaaaa`, []int{0, 3},
		},
		{
			`\(aba\)*`, `abbaabc`, []int{0, 0},
		},
		{
			`\(aba\)\+`, `abaaabc`, []int{0, 3},
		},
		{
			`\(aba\)\+`, `abaabac`, []int{0, 6},
		},
		{
			`\(aba\)\+`, `aaaabc`, nil,
		},
		{
			`\(aba\)\?`, ``, []int{0, 0},
		},
		{
			`\(aba\)\?`, `aba`, []int{0, 3},
		},
		{
			`\(aba\)\?`, `ccc`, []int{0, 0},
		},
	} {
		re, err := Compile(test.expr)
		if err != nil {
			t.Fatalf("compiling (%d, %q): %v", i, test.s, err)
		}
		loc := re.FindStringIndex(test.s)
		if loc == nil || test.loc == nil {
			if loc != nil && test.loc == nil {
				t.Errorf("FindStringIndex(%d, %q): want nil, but got %v", i, test.s, loc)
			} else if loc == nil && test.loc != nil {
				t.Errorf("FindStringIndex(%d, %q): want %v, but got nil", i, test.s, test.loc)
			}
		} else if loc[0] != test.loc[0] || loc[1] != test.loc[1] {
			t.Errorf("FindStringIndex(%d, %q): want %v, but got %v", i, test.s, test.loc, loc)
		}
	}
}
