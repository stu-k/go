package parse_test

import (
	"errors"
	"testing"

	errs "github.com/stu-k/go/parser/errors"
	"github.com/stu-k/go/parser/parse"
)

func TestAlpha(t *testing.T) {
	type testobj struct {
		in   string
		want string
		rest string
		err  error
	}

	rulemap := make(map[*parse.Rule][]testobj)

	rulemap[parse.Alpha] = []testobj{
		{"abc", "abc", "", nil},
		{"   abc", "abc", "", nil},
		{"   a   bc", "abc", "", nil},
		{"   a   b   c", "abc", "", nil},
		{"   a   b   c   ", "abc", "", nil},
		{"   a   ", "a", "", nil},
		{"   a   .", "a", ".", nil},
		{"   a   .  ", "a", ".  ", nil},
		{"ab.c", "ab", ".c", nil},
		{"abc.", "abc", ".", nil},

		{".", "", "", errs.ErrBadMatch},
		{"", "", "", errs.ErrBadMatch},
	}

	rulemap[parse.Alpha.Count(3)] = []testobj{
		{"abc", "abc", "", nil},
		{"abcd", "abc", "d", nil},
		{"abc.", "abc", ".", nil},

		{"a", "", "", errs.ErrBadMatch},
		{"ab", "", "", errs.ErrBadMatch},

		{"abc ", "abc", "", nil},
		{"   abc ", "abc", "", nil},
		{"   abcd", "abc", "d", nil},
		{"   a    b c  d", "abc", "d", nil},

		{".a", "", "", errs.ErrBadMatch},
		{"a.", "", "", errs.ErrBadMatch},
		{".ab", "", "", errs.ErrBadMatch},
		{"ab.", "", "", errs.ErrBadMatch},
		{"a.b", "", "", errs.ErrBadMatch},

		{".", "", "", errs.ErrBadMatch},
		{"", "", "", errs.ErrBadMatch},
	}

	rulemap[parse.Alpha.Capture(false)] = []testobj{
		{"a", "", "", nil},
		{"ab", "", "", nil},
		{"abc", "", "", nil},
		{"a.", "", ".", nil},
		{"a   .", "", ".", nil},

		{".", "", "", errs.ErrBadMatch},
	}

	for rule, tests := range rulemap {
		for _, test := range tests {
			got, rest, err := rule.Parse(test.in)
			if !errors.Is(err, test.err) {
				t.Errorf("for \"%s\" expected error \"%v\"; got \"%v\"", test.in, test.err, err)
			}

			if !eq(got, ss(test.want)) {
				t.Errorf("for \"%s\" expected output \"%v\"; got \"%v\"", test.in, test.want, got)
			}

			if rest != test.rest {
				t.Errorf("for \"%s\" expected remainder \"%v\"; got \"%v\"", test.in, test.rest, rest)
			}
		}
	}
}
