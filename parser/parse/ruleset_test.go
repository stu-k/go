package parse_test

import (
	"errors"
	"testing"

	errs "github.com/stu-k/go/parser/errors"
	"github.com/stu-k/go/parser/parse"
)

var ss = func(s ...string) []string {
	if len(s) == 0 {
		return nil
	}
	if len(s) == 1 && s[0] == "" {
		return nil
	}
	return s
}
var eq = func(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if b[i] != v {
			return false
		}
	}
	return true
}

func TestRuleset(t *testing.T) {
	type rulesettest struct {
		in   string
		rs   *parse.Ruleset
		want []string
		rest string
		err  error
	}

	rstests := make(map[*parse.Ruleset][]rulesettest)

	ruleset := parse.NewRuleset("alpha", parse.RuleAlpha)
	rstests[ruleset] = []rulesettest{
		{"abc", ruleset, ss("abc"), "", nil},
		{"ab.c", ruleset, ss("ab"), ".c", nil},
		{"abc.", ruleset, ss("abc"), ".", nil},
		{".", ruleset, ss(), "", errs.ErrBadMatch},
	}

	ruleset = parse.NewRuleset("num", parse.RuleNum)
	rstests[ruleset] = []rulesettest{
		{"123", ruleset, ss("123"), "", nil},
		{"12.3", ruleset, ss("12"), ".3", nil},
		{"123.", ruleset, ss("123"), ".", nil},
	}

	ruleset = parse.NewRuleset(
		"alphanum",
		parse.RuleAlpha,
		parse.RuleNum,
	)
	rstests[ruleset] = []rulesettest{
		{"a1", ruleset, ss("a", "1"), "", nil},
		{"abc123", ruleset, ss("abc", "123"), "", nil},
		{"a1.", ruleset, ss("a", "1"), ".", nil},
		{"a", ruleset, ss(), "", errs.ErrBadMatch},
		{"1", ruleset, ss(), "", errs.ErrBadMatch},
		{"a.1", ruleset, ss(), "", errs.ErrBadMatch},
	}

	ruleset = parse.NewRuleset(
		"kv(var var)",
		parse.RuleAlpha,
		parse.RuleAny.Chars(":"),
		parse.RuleNum,
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

	ruleset = parse.NewRuleset(
		"obj(kv(var var))",
		parse.RuleAny.Chars("{"),
		parse.RuleAlpha,
		parse.RuleAny.Chars(":"),
		parse.RuleAlpha,
		parse.RuleAny.Chars("}"),
	)
	rstests[ruleset] = []rulesettest{
		{"{a:x}", ruleset, ss("{", "a", ":", "x", "}"), "", nil},
		{"{abc:xyz}", ruleset, ss("{", "abc", ":", "xyz", "}"), "", nil},
		{".", ruleset, ss(), "", errs.ErrBadMatch},
		{"{a:x}...", ruleset, ss("{", "a", ":", "x", "}"), "...", nil},
	}

	ruleset = parse.NewRuleset(
		"obj(kv(_var_ var))",
		parse.RuleAny.Chars("{"),
		parse.RuleAny.Chars("_"),
		parse.RuleAlpha,
		parse.RuleAny.Chars("_"),
		parse.RuleAny.Chars(":"),
		parse.RuleAlpha,
		parse.RuleAny.Chars("}"),
	)
	rstests[ruleset] = []rulesettest{
		{"{_a_:x}", ruleset, ss("{", "_", "a", "_", ":", "x", "}"), "", nil},
		{"{_abc_:x}", ruleset, ss("{", "_", "abc", "_", ":", "x", "}"), "", nil},
	}

	for rs, tests := range rstests {
		t.Run(rs.Name(), func(t *testing.T) {
			for _, test := range tests {
				got, rest, err := test.rs.Parse(test.in)
				if !eq(got, test.want) {
					t.Errorf("for \"%v\" expected output \"%v\"; got \"%v\"", test.in, test.want, got)
				}

				if rest != test.rest {
					t.Errorf("for \"%v\" expected remainder \"%v\"; got \"%v\"", test.in, test.rest, rest)
				}
				if !errors.Is(err, test.err) {
					t.Errorf("for \"%v\" expected error \"%v\"; got \"%v\"", test.in, test.err, err)
				}

			}
		})
	}
}

func TestRulesetUntilFail(t *testing.T) {
	mk := func(s, r string) (*parse.RulesetUntilFail, error) {
		rfs, err := parse.NewRulesetFromStr(s, r)
		if err != nil {
			return nil, err
		}
		return rfs.UntilFail(), nil
	}

	type rulesettest struct {
		in   string
		rs   *parse.RulesetUntilFail
		want []string
		rest string
		err  error
	}

	rstests := make(map[*parse.RulesetUntilFail][]rulesettest)

	rs, err := mk("alpha 1", "ralpha, #1")
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

	rs, err = mk("alpha comma", "ralpha, #1 | c,, #1")
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

	rs, err = mk("str: num",
		"ralpha | c:, #1 | rnum",
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
				got, rest, err := test.rs.Parse(test.in)
				if !eq(got, test.want) {
					t.Errorf("for \"%v\" expected output \"%v\"; got \"%v\"", test.in, test.want, got)
				}

				if rest != test.rest {
					t.Errorf("for \"%v\" expected remainder \"%v\"; got \"%v\"", test.in, test.rest, rest)
				}
				if !errors.Is(err, test.err) {
					t.Errorf("for \"%v\" expected error \"%v\"; got \"%v\"", test.in, test.err, err)
				}

			}
		})
	}
}
