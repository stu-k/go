package parse_test

import (
	"errors"
	"testing"

	errs "github.com/stu-k/go/parser/errors"
	"github.com/stu-k/go/parser/parse"
)

func TestAlpha(t *testing.T) {
	rule := parse.Alpha
	t.Run("ignore space", func(t *testing.T) {
		tt := []struct {
			in   string
			want string
			rest string
			err  error
		}{
			{"abc", "abc", "", nil},
			{"   abc", "abc", "", nil},
			{"   a   bc", "abc", "", nil},
			{"   a   b   c", "abc", "", nil},
			{"   a   b   c   ", "abc", "", nil},
			{"ab.c", "ab", ".c", nil},
			{"abc.", "abc", ".", nil},

			{".", "", "", errs.ErrBadMatch},
			{"", "", "", errs.ErrBadMatch},
		}

		for _, test := range tt {
			got, rest, err := rule.Parse(test.in)
			if !errors.Is(err, test.err) {
				t.Errorf("expected error \"%v\"; got \"%v\"", test.err, err)
			}

			if got != test.want {
				t.Errorf("expected output \"%v\"; got \"%v\"", test.want, got)
			}

			if rest != test.rest {
				t.Errorf("expected remainder \"%v\"; got \"%v\"", test.rest, rest)
			}
		}
	})

	rule = parse.Alpha.WithCount(3)
	t.Run("count 3", func(t *testing.T) {
		tt := []struct {
			in   string
			want string
			rest string
			err  error
		}{
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

		for _, test := range tt {
			got, rest, err := rule.Parse(test.in)
			if !errors.Is(err, test.err) {
				t.Errorf("for \"%v\" expected error \"%v\"; got \"%v\"", test.in, test.err, err)
			}

			if got != test.want {
				t.Errorf("for \"%v\" expected output \"%v\"; got \"%v\"", test.in, test.want, got)
			}

			if rest != test.rest {
				t.Errorf("for \"%v\" expected remainder \"%v\"; got \"%v\"", test.in, test.rest, rest)
			}
		}
	})

	rule = parse.Alpha.WrapWith('_')
	t.Run("wrap with _", func(t *testing.T) {
		tt := []struct {
			in   string
			want string
			rest string
			err  error
		}{
			{"_a_", "_a_", "", nil},
			{"_ab_", "_ab_", "", nil},
			{"_abc_", "_abc_", "", nil},
			{"_abcdefg_", "_abcdefg_", "", nil},
			{"_abc", "", "", errs.ErrBadMatch},

			{"_a_bc", "_a_", "bc", nil},
			{"_a.bc", "", "", errs.ErrBadMatch},
		}

		for _, test := range tt {
			got, rest, err := rule.Parse(test.in)
			if !errors.Is(err, test.err) {
				t.Errorf("for \"%v\" expected error \"%v\"; got \"%v\"", test.in, test.err, err)
			}

			if got != test.want {
				t.Errorf("for \"%v\" expected output \"%v\"; got \"%v\"", test.in, test.want, got)
			}

			if rest != test.rest {
				t.Errorf("for \"%v\" expected remainder \"%v\"; got \"%v\"", test.in, test.rest, rest)
			}
		}
	})
}
