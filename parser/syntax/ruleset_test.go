package syntax_test

import (
	"errors"
	"testing"

	errs "github.com/stu-k/go/parser/errors"
	stx "github.com/stu-k/go/parser/syntax"
)

func TestRuleset(t *testing.T) {
	type rulesettest struct {
		in   string
		rs   *stx.Ruleset
		want []string
		rest string
		err  error
	}

	rstests := make(map[*stx.Ruleset][]rulesettest)

	ruleset := stx.NewRuleset("alpha", stx.RuleAlpha)
	rstests[ruleset] = []rulesettest{
		{"abc", ruleset, ss("abc"), "", nil},
		{"ab.c", ruleset, ss("ab"), ".c", nil},
		{"abc.", ruleset, ss("abc"), ".", nil},
		{".", ruleset, ss(), "", errs.ErrBadMatch},
	}

	ruleset = stx.NewRuleset("num", stx.RuleNum)
	rstests[ruleset] = []rulesettest{
		{"123", ruleset, ss("123"), "", nil},
		{"12.3", ruleset, ss("12"), ".3", nil},
		{"123.", ruleset, ss("123"), ".", nil},
	}

	ruleset = stx.NewRuleset(
		"alphanum",
		stx.RuleAlpha,
		stx.RuleNum,
	)
	rstests[ruleset] = []rulesettest{
		{"a1", ruleset, ss("a", "1"), "", nil},
		{"abc123", ruleset, ss("abc", "123"), "", nil},
		{"a1.", ruleset, ss("a", "1"), ".", nil},
		{"a", ruleset, ss(), "", errs.ErrBadMatch},
		{"1", ruleset, ss(), "", errs.ErrBadMatch},
		{"a.1", ruleset, ss(), "", errs.ErrBadMatch},
	}

	ruleset = stx.NewRuleset(
		"kv(var var)",
		stx.RuleAlpha,
		stx.RuleAny.Chars(":"),
		stx.RuleNum,
	)
	rstests[ruleset] = []rulesettest{
		{"a:1", ruleset, ss("a", ":", "1"), "", nil},
		{"abc:123", ruleset, ss("abc", ":", "123"), "", nil},
		{"a:1:", ruleset, ss("a", ":", "1"), ":", nil},
		{".a:1", ruleset, ss(), "", errs.ErrBadMatch},
		{"a.:1", ruleset, ss(), "", errs.ErrBadMatch},
		{"a:.1", ruleset, ss(), "", errs.ErrBadMatch},
		{"a:1.", ruleset, ss("a", ":", "1"), ".", nil},
	}

	ruleset = stx.NewRuleset(
		"obj(kv(var var))",
		stx.RuleAny.Chars("{"),
		stx.RuleAlpha,
		stx.RuleAny.Chars(":"),
		stx.RuleAlpha,
		stx.RuleAny.Chars("}"),
	)
	rstests[ruleset] = []rulesettest{
		{"{a:x}", ruleset, ss("{", "a", ":", "x", "}"), "", nil},
		{"{abc:xyz}", ruleset, ss("{", "abc", ":", "xyz", "}"), "", nil},
		{".", ruleset, ss(), "", errs.ErrBadMatch},
		{"{a:x}...", ruleset, ss("{", "a", ":", "x", "}"), "...", nil},
	}

	ruleset = stx.NewRuleset(
		"obj(kv(_var_ var))",
		stx.RuleAny.Chars("{"),
		stx.RuleAny.Chars("_"),
		stx.RuleAlpha,
		stx.RuleAny.Chars("_"),
		stx.RuleAny.Chars(":"),
		stx.RuleAlpha,
		stx.RuleAny.Chars("}"),
	)
	rstests[ruleset] = []rulesettest{
		{"{_a_:x}", ruleset, ss("{", "_", "a", "_", ":", "x", "}"), "", nil},
		{"{_abc_:x}", ruleset, ss("{", "_", "abc", "_", ":", "x", "}"), "", nil},
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

func TestRulesetUntilFail(t *testing.T) {
	mk := func(s string, r ...string) (stx.Parsable, error) {
		rfs, err := stx.NewRulesetFromStrs(s, r...)
		if err != nil {
			return nil, err
		}
		return rfs.UntilFail(), nil
	}

	type rulesettest struct {
		in   string
		rs   stx.Parsable
		want []string
		rest string
		err  error
	}

	rstests := make(map[stx.Parsable][]rulesettest)

	rs, err := mk(
		"alpha 1",
		"ralpha, #1",
	)
	if err != nil {
		t.Fatalf("ruleset creation failed: %v", err)
	}
	rstests[rs] = []rulesettest{
		{"a", rs, ss("a"), "", nil},
		{"ab", rs, ss("a", "b"), "", nil},
		{"abc", rs, ss("a", "b", "c"), "", nil},

		{"a.", rs, ss("a"), ".", nil},
		{"ab.", rs, ss("a", "b"), ".", nil},

		{".", rs, ss(), "", errs.ErrBadMatch},
	}

	rs, err = mk(
		"alpha comma",
		"ralpha, #1", "c,, #1",
	)
	if err != nil {
		t.Fatalf("ruleset creation failed: %v", err)
	}
	rstests[rs] = []rulesettest{
		{"a,", rs, ss("a", ","), "", nil},
		{"a,b,", rs, ss("a", ",", "b", ","), "", nil},
		{"a,b", rs, ss("a", ","), "b", nil},
		{"a,", rs, ss("a", ","), "", nil},

		{".", rs, ss(), "", errs.ErrBadMatch},
	}

	rs, err = mk(
		"str: num",
		"ralpha", "c:, #1", "rnum",
	)
	if err != nil {
		t.Fatalf("ruleset creation failed: %v", err)
	}
	rstests[rs] = []rulesettest{
		{"abc:123", rs, ss("abc", ":", "123"), "", nil},
		{"a:1b:2", rs, ss("a", ":", "1", "b", ":", "2"), "", nil},
		{"a:1", rs, ss("a", ":", "1"), "", nil},
		{"a:1b:", rs, ss("a", ":", "1"), "b:", nil},

		{".", rs, ss(), "", errs.ErrBadMatch},
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

func TestRulesetOneOf(t *testing.T) {
	mk := func(s string, r ...string) (stx.Parsable, error) {
		rfs, err := stx.NewRulesetFromStrs(s, r...)
		if err != nil {
			return nil, err
		}
		return rfs.OneOf(), nil
	}

	type rulesettest struct {
		in   string
		want []string
		err  error
	}

	rstests := make(map[stx.Parsable][]rulesettest)

	rs1, err := mk(
		"alp",
		"ralpha",
	)
	if err != nil {
		t.Fatalf("ruleset creation failed: %v", err)
	}

	rs2, err := mk(
		"1a",
		"ca, #1",
	)
	if err != nil {
		t.Fatalf("ruleset creation failed: %v", err)
	}

	rs3, err := mk(
		"2alp",
		"ralpha, #2",
	)
	if err != nil {
		t.Fatalf("ruleset creation failed: %v", err)
	}

	rs4, err := mk(
		"2b",
		"cb, #2",
	)
	if err != nil {
		t.Fatalf("ruleset creation failed: %v", err)
	}

	rs := stx.NewRuleset("a0.", rs1, rs2, rs3, rs4).OneOf()
	rstests[rs] = []rulesettest{
		{"a", ss("alp", "1a"), nil},
		{"aa", ss("alp", "1a", "2alp"), nil},
		{"aaa", ss("alp", "1a", "2alp"), nil},

		{"b", ss("alp"), nil},
		{"bb", ss("alp", "2b", "2alp"), nil},
		{"bbb", ss("alp", "2b", "2alp"), nil},

		{"x", ss("alp"), nil},
		{"xy", ss("alp", "2alp"), nil},

		{"", ss(), errs.ErrBadMatch},
		{".", ss(), errs.ErrBadMatch},
	}

	rs1, err = mk(
		"var",
		"ralpha",
	)
	if err != nil {
		t.Fatalf("ruleset creation failed: %v", err)
	}

	rs2, err = mk(
		"str",
		"c', #1", "ralpha", "c', #1",
	)
	if err != nil {
		t.Fatalf("ruleset creation failed: %v", err)
	}

	rs = stx.NewRuleset("str / var", rs1, rs2).OneOf()
	rstests[rs] = []rulesettest{
		{"'str'", ss("str"), nil},
		{"var", ss("var"), nil},

		{"", ss(), errs.ErrBadMatch},
		{".", ss(), errs.ErrBadMatch},
	}

	for rs, tests := range rstests {
		t.Run(rs.Name(), func(t *testing.T) {
			for _, test := range tests {
				got, err := rs.Parse(test.in)
				if !errors.Is(err, test.err) {
					t.Fatalf("for \"%v\" expected error \"%v\"; got \"%v\"", test.in, test.err, err)
				}

				if err != nil {
					return
				}

				if got == nil {
					t.Fatalf("for \"%v\" expected output not to be nil", test.in)
				}

				for _, v := range test.want {
					_, ok := got.NameMap()[v]
					if !ok {
						t.Fatalf("expected key \"%v\" in nameMap %v", v, got.NameMap())
					}
				}
			}
		})
	}
}
