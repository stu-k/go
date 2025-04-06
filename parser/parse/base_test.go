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
}
