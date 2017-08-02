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
			t.Fatalf("compiling (%d, %q) : %v", i, test.s, err)
		}
		ok := re.MatchString(test.s)
		if ok != test.ok {
			t.Errorf("MatchString(%d, %q) should be %v", i, test.s, test.ok)
		}
	}
}
