package syntax_test

import (
	"errors"
	"testing"
	"unicode"

	errs "github.com/stu-k/go/parser/errors"
	stx "github.com/stu-k/go/parser/syntax"
)

func TestRule(t *testing.T) {
	type testobj struct {
		in   string
		want string
		rest string
		err  error
	}

	rulemap := make(map[*stx.Rule][]testobj)

	rulemap[stx.RuleAlpha] = []testobj{
		{"abc", "abc", "", nil},
		{"abc.", "abc", ".", nil},
		{"ab.c", "ab", ".c", nil},
		{"a.bc", "a", ".bc", nil},
		{".abc", "", "", errs.ErrBadMatch},

		{".", "", "", errs.ErrBadMatch},
		{"", "", "", errs.ErrBadMatch},
	}

	rulemap[stx.RuleNum] = []testobj{
		{"123", "123", "", nil},
		{"123.", "123", ".", nil},
		{"12.3", "12", ".3", nil},
		{"1.23", "1", ".23", nil},
		{".123", "", "", errs.ErrBadMatch},

		{".", "", "", errs.ErrBadMatch},
		{"", "", "", errs.ErrBadMatch},
	}

	rulemap[stx.RuleAlpha.Count(3)] = []testobj{
		{"abc", "abc", "", nil},
		{"abcd", "abc", "d", nil},
		{"abc.", "abc", ".", nil},
		{"ab", "", "", errs.ErrBadMatch},

		{".", "", "", errs.ErrBadMatch},
		{"", "", "", errs.ErrBadMatch},
	}

	rulemap[stx.RuleAlpha.Capture(false)] = []testobj{
		{"a", "", "", nil},
		{"ab", "", "", nil},
		{"abc", "", "", nil},
		{"abc.", "", ".", nil},
		{"ab.c", "", ".c", nil},
		{"a.bc", "", ".bc", nil},
		{".abc", "", "", errs.ErrBadMatch},

		{".", "", "", errs.ErrBadMatch},
		{"", "", "", errs.ErrBadMatch},
	}

	rulemap[stx.RuleAny.Check(unicode.IsNumber)] = []testobj{
		{"1", "1", "", nil},
		{"12", "12", "", nil},
		{"123", "123", "", nil},
		{"1.", "1", ".", nil},
		{".1", "", "", errs.ErrBadMatch},

		{".", "", "", errs.ErrBadMatch},
		{"", "", "", errs.ErrBadMatch},
	}

	rulemap[stx.RuleAny.Chars("a1")] = []testobj{
		{"a", "a", "", nil},
		{"aa", "aa", "", nil},
		{"1", "1", "", nil},
		{"11", "11", "", nil},
		{"a1", "a1", "", nil},
		{"aa1", "aa1", "", nil},
		{"aa11", "aa11", "", nil},
		{"a1a1", "a1a1", "", nil},
		{"1a1a", "1a1a", "", nil},

		{"a1b2", "a1", "b2", nil},
		{".a1", "", "", errs.ErrBadMatch},

		{".", "", "", errs.ErrBadMatch},
		{"", "", "", errs.ErrBadMatch},
	}

	for rule, tests := range rulemap {
		for _, test := range tests {
			got, err := rule.Parse(test.in)
			if !errors.Is(err, test.err) {
				t.Fatalf("for \"%s\" expected error \"%v\"; got \"%v\"", test.in, test.err, err)
			}

			if err != nil {
				return
			}

			if got == nil {
				t.Fatalf("for \"%s\" expected output not to be nil", test.in)
			}

			if !eq(got.Strings(), ss(test.want)) {
				t.Errorf("for \"%s\" expected output \"%v\"; got \"%v\"", test.in, test.want, got)
			}

			if got.Rest() != test.rest {
				t.Errorf("for \"%s\" expected remainder \"%v\"; got \"%v\"", test.in, test.rest, got.Rest())
			}
		}
	}
}
