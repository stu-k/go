package syntax_test

import (
	"errors"
	"testing"

	errs "github.com/stu-k/go/parser/errors"
	stx "github.com/stu-k/go/parser/syntax"
)

func TestRulesetFromStrs(t *testing.T) {
	type rulesettest struct {
		in   string
		rs   *stx.Sequence
		want []string
		rest string
		err  error
	}

	rstests := make(map[*stx.Sequence][]rulesettest)

	ruleset, err := stx.NewSequenceFromStrs(
		"alpha",
		"ralpha")
	if err != nil {
		t.Fatalf("ruleset creation failed: %v", err)
	}

	rstests[ruleset] = []rulesettest{
		{"abc", ruleset, ss("abc"), "", nil},
		{"ab.c", ruleset, ss("ab"), ".c", nil},
		{"abc.", ruleset, ss("abc"), ".", nil},
		{".", ruleset, ss(), "", errs.ErrBadMatch},
	}

	ruleset, err = stx.NewSequenceFromStrs(
		"num",
		"rnum",
	)
	if err != nil {
		t.Fatalf("ruleset creation failed: %v", err)
	}

	rstests[ruleset] = []rulesettest{
		{"123", ruleset, ss("123"), "", nil},
		{"12.3", ruleset, ss("12"), ".3", nil},
		{"123.", ruleset, ss("123"), ".", nil},
	}

	ruleset, err = stx.NewSequenceFromStrs(
		"alphanum",
		"ralpha", "rnum",
	)
	if err != nil {
		t.Fatalf("ruleset creation failed: %v", err)
	}

	rstests[ruleset] = []rulesettest{
		{"a1", ruleset, ss("a", "1"), "", nil},
		{"abc123", ruleset, ss("abc", "123"), "", nil},
		{"a1.", ruleset, ss("a", "1"), ".", nil},
		{"a", ruleset, ss(), "", errs.ErrBadMatch},
		{"1", ruleset, ss(), "", errs.ErrBadMatch},
		{"a.1", ruleset, ss(), "", errs.ErrBadMatch},
	}

	ruleset, err = stx.NewSequenceFromStrs(
		"kv(var var)",
		"ralpha", "c:, #1", "rnum",
	)
	if err != nil {
		t.Fatalf("ruleset creation failed: %v", err)
	}

	rstests[ruleset] = []rulesettest{
		{"a:1", ruleset, ss("a", ":", "1"), "", nil},
		{"abc:123", ruleset, ss("abc", ":", "123"), "", nil},

		{"a::1", ruleset, nil, "", errs.ErrBadMatch},
		{":1", ruleset, nil, "", errs.ErrBadMatch},
		{"a1", ruleset, nil, "", errs.ErrBadMatch},
		{"a:", ruleset, nil, "", errs.ErrBadMatch},

		{".a:1", ruleset, nil, "", errs.ErrBadMatch},
		{"a.:1", ruleset, nil, "", errs.ErrBadMatch},
		{"a:.1", ruleset, nil, "", errs.ErrBadMatch},
		{"a:1.", ruleset, ss("a", ":", "1"), ".", nil},
	}

	ruleset, err = stx.NewSequenceFromStrs(
		"obj(kv(var var))",
		"c{, #1", "ralpha", "c:, #1", "ralpha", "c}, #1",
	)
	if err != nil {
		t.Fatalf("ruleset creation failed: %v", err)
	}

	rstests[ruleset] = []rulesettest{
		{"{a:x}", ruleset, ss("{", "a", ":", "x", "}"), "", nil},
		{"{abc:xyz}", ruleset, ss("{", "abc", ":", "xyz", "}"), "", nil},

		{".", ruleset, ss(), "", errs.ErrBadMatch},
		{".{a:x}", ruleset, ss(), "", errs.ErrBadMatch},
		{"{.a:x}", ruleset, ss(), "", errs.ErrBadMatch},
		{"{a.:x}", ruleset, ss(), "", errs.ErrBadMatch},
		{"{a:.x}", ruleset, ss(), "", errs.ErrBadMatch},
		{"{a:x.}", ruleset, ss(), "", errs.ErrBadMatch},
		{"{a:x}.", ruleset, ss("{", "a", ":", "x", "}"), ".", nil},
	}

	ruleset, err = stx.NewSequenceFromStrs(
		"obj(kv(_var_ var))",
		"c{, #1", "c_, #1", "ralpha", "c_, #1", "c:, #1", "ralpha", "c}, #1",
	)
	if err != nil {
		t.Fatalf("ruleset creation failed: %v", err)
	}

	rstests[ruleset] = []rulesettest{
		{"{_a_:x}", ruleset, ss("{", "_", "a", "_", ":", "x", "}"), "", nil},
		{"{_abc_:xyz}", ruleset, ss("{", "_", "abc", "_", ":", "xyz", "}"), "", nil},

		{".", ruleset, ss(), "", errs.ErrBadMatch},
	}

	ruleset, err = stx.NewSequenceFromStrs(
		"test special vals",
		"c,, #1", "rnum", "c|, #1",
	)
	if err != nil {
		t.Fatalf("ruleset creation failed: %v", err)
	}

	rstests[ruleset] = []rulesettest{
		{".", ruleset, ss(), "", errs.ErrBadMatch},

		{",1|", ruleset, ss(",", "1", "|"), "", nil},
	}

	ruleset, err = stx.NewSequenceFromStrs(
		"alpha comma",
		"ralpha, #1", "c,, #1", "ralpha, #1", "c,, #1", "ralpha, #1", "c,, #1",
	)
	if err != nil {
		t.Fatalf("ruleset creation failed: %v", err)
	}

	rstests[ruleset] = []rulesettest{
		{".", ruleset, ss(), "", errs.ErrBadMatch},
		{"a,b,c,", ruleset, ss("a", ",", "b", ",", "c", ","), "", nil},
	}

	ruleset, err = stx.NewSequenceFromStrs(
		"alpha3",
		"ralpha, #1", "ralpha, #1", "ralpha, #1",
	)
	if err != nil {
		t.Fatalf("ruleset creation failed: %v", err)
	}

	rstests[ruleset] = []rulesettest{
		{".", ruleset, ss(), "", errs.ErrBadMatch},
		{"abc", ruleset, ss("a", "b", "c"), "", nil},
	}

	ruleset, err = stx.NewSequenceFromStrs(
		"capture",
		"ralpha", "c:, g0", "rnum",
	)
	if err != nil {
		t.Fatalf("ruleset creation failed: %v", err)
	}

	rstests[ruleset] = []rulesettest{
		{"a:1", ruleset, ss("a", "1"), "", nil},
		{"abc:123", ruleset, ss("abc", "123"), "", nil},
		{"a::1", ruleset, ss("a", "1"), "", nil},

		{".", ruleset, ss(), "", errs.ErrBadMatch},
	}

	for rs, tests := range rstests {
		t.Run(rs.Name(), func(t *testing.T) {
			for _, test := range tests {
				got, err := test.rs.Parse(test.in)
				if !errors.Is(err, test.err) {
					t.Fatalf("for \"%v\" expected error \"%v\"; got \"%v\"", test.in, test.err, err)
				}

				if err != nil {
					return
				}

				if got == nil {
					t.Fatalf("for \"%v\" expected output not to be nil", test.in)
				}

				if !eq(got.Strings(), test.want) {
					t.Errorf("for \"%v\" expected output \"%v\"; got \"%v\"", test.in, test.want, got)
				}

				if got.Rest() != test.rest {
					t.Errorf("for \"%v\" expected remainder \"%v\"; got \"%v\"", test.in, test.rest, got.Rest())
				}
			}
		})
	}
}
