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
	ss := func(s ...string) []string { return s }
	eq := func(a, b []string) bool {
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

	type rulesettest struct {
		in   string
		rs   *parse.Ruleset
		want []string
		rest string
		err  error
	}

	rstests := make(map[*parse.Ruleset][]rulesettest)

	ruleset := parse.NewRuleset("alpha", parse.Alpha)
	rstests[ruleset] = []rulesettest{
		{"abc", ruleset, ss("abc"), "", nil},
		{"ab.c", ruleset, ss("ab"), ".c", nil},
		{"abc.", ruleset, ss("abc"), ".", nil},
		{".", ruleset, ss(), "", errs.ErrBadMatch},
	}

	ruleset = parse.NewRuleset("num", parse.Numeric)
	rstests[ruleset] = []rulesettest{
		{"123", ruleset, ss("123"), "", nil},
		{"12.3", ruleset, ss("12"), ".3", nil},
		{"123.", ruleset, ss("123"), ".", nil},
	}

	ruleset = parse.NewRuleset(
		"alphanum",
		parse.Alpha,
		parse.Numeric,
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
		parse.Alpha,
		parse.FromChar(':'),
		parse.Numeric,
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
		parse.FromChar('{'),
		parse.Alpha,
		parse.FromChar(':'),
		parse.Alpha,
		parse.FromChar('}'),
	)
	rstests[ruleset] = []rulesettest{
		{"{a:x}", ruleset, ss("{", "a", ":", "x", "}"), "", nil},
		{"{abc:xyz}", ruleset, ss("{", "abc", ":", "xyz", "}"), "", nil},
		{".", ruleset, ss(), "", errs.ErrBadMatch},
		{"{a:x}...", ruleset, ss("{", "a", ":", "x", "}"), "...", nil},
	}

	ruleset = parse.NewRuleset(
		"obj(kv(_var_ var))",
		parse.FromChar('{'),
		parse.Alpha.Wrap('_'),
		parse.FromChar(':'),
		parse.Alpha,
		parse.FromChar('}'),
	)
	rstests[ruleset] = []rulesettest{
		{"{_a_:x}", ruleset, ss("{", "_a_", ":", "x", "}"), "", nil},
		{"{_abc_:x}", ruleset, ss("{", "_abc_", ":", "x", "}"), "", nil},
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

func TestRulesetFromString(t *testing.T) {
	ss := func(s ...string) []string { return s }
	eq := func(a, b []string) bool {
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

	type rulesettest struct {
		in   string
		rs   *parse.Ruleset
		want []string
		rest string
		err  error
	}

	rstests := make(map[*parse.Ruleset][]rulesettest)

	ruleset, err := parse.NewRulesetFromStr("alpha", "ralpha")
	if err != nil {
		t.Fatalf("ruleset creation failed: %v", err)
	}

	rstests[ruleset] = []rulesettest{
		{"abc", ruleset, ss("abc"), "", nil},
		{"ab.c", ruleset, ss("ab"), ".c", nil},
		{"abc.", ruleset, ss("abc"), ".", nil},
		{".", ruleset, ss(), "", errs.ErrBadMatch},
	}

	ruleset, err = parse.NewRulesetFromStr("num", "rnum")
	if err != nil {
		t.Fatalf("ruleset creation failed: %v", err)
	}

	rstests[ruleset] = []rulesettest{
		{"123", ruleset, ss("123"), "", nil},
		{"12.3", ruleset, ss("12"), ".3", nil},
		{"123.", ruleset, ss("123"), ".", nil},
	}

	ruleset, err = parse.NewRulesetFromStr("alphanum", "ralpha | rnum")
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

	ruleset, err = parse.NewRulesetFromStr(
		"kv(var var)",
		"ralpha | c:, #1 | rnum",
	)
	if err != nil {
		t.Fatalf("ruleset creation failed: %v", err)
	}

	rstests[ruleset] = []rulesettest{
		{"a:1", ruleset, ss("a", ":", "1"), "", nil},

		{"abc:123", ruleset, ss("abc", ":", "123"), "", nil},

		{".a:1", ruleset, nil, "", errs.ErrBadMatch},
		{"a.:1", ruleset, nil, "", errs.ErrBadMatch},
		{"a:.1", ruleset, nil, "", errs.ErrBadMatch},
		{"a:1.", ruleset, ss("a", ":", "1"), ".", nil},

		{" a:1", ruleset, ss("a", ":", "1"), "", nil},
		{"a :1", ruleset, ss("a", ":", "1"), "", nil},
		{"a: 1", ruleset, ss("a", ":", "1"), "", nil},
		{"a:1 .", ruleset, ss("a", ":", "1"), ".", nil},
	}

	ruleset, err = parse.NewRulesetFromStr(
		"obj(kv(var var))",
		"c{, #1 | ralpha | c:, #1 | ralpha | c}, #1",
	)
	if err != nil {
		t.Fatalf("ruleset creation failed: %v", err)
	}

	rstests[ruleset] = []rulesettest{
		{"{a:x}", ruleset, ss("{", "a", ":", "x", "}"), "", nil},
		{" {a:x}", ruleset, ss("{", "a", ":", "x", "}"), "", nil},
		{"{ a:x}", ruleset, ss("{", "a", ":", "x", "}"), "", nil},
		{"{a :x}", ruleset, ss("{", "a", ":", "x", "}"), "", nil},
		{"{a: x}", ruleset, ss("{", "a", ":", "x", "}"), "", nil},
		{"{a:x }", ruleset, ss("{", "a", ":", "x", "}"), "", nil},
		{"{a:x} ", ruleset, ss("{", "a", ":", "x", "}"), "", nil},
		{"{abc:xyz}", ruleset, ss("{", "abc", ":", "xyz", "}"), "", nil},

		{".", ruleset, ss(), "", errs.ErrBadMatch},
		{".{a:x}", ruleset, ss(), "", errs.ErrBadMatch},
		{"{.a:x}", ruleset, ss(), "", errs.ErrBadMatch},
		{"{a.:x}", ruleset, ss(), "", errs.ErrBadMatch},
		{"{a:.x}", ruleset, ss(), "", errs.ErrBadMatch},
		{"{a:x.}", ruleset, ss(), "", errs.ErrBadMatch},
		{"{a:x}.", ruleset, ss("{", "a", ":", "x", "}"), ".", nil},
	}

	ruleset, err = parse.NewRulesetFromStr(
		"obj(kv(_var_ var))",
		"c{, #1 | ralpha, w_ | c:, #1 | ralpha | c}, #1",
	)
	if err != nil {
		t.Fatalf("ruleset creation failed: %v", err)
	}

	rstests[ruleset] = []rulesettest{
		{"{_a_:x}", ruleset, ss("{", "_a_", ":", "x", "}"), "", nil},

		{"{_abc_:xyz}", ruleset, ss("{", "_abc_", ":", "xyz", "}"), "", nil},

		{".", ruleset, nil, "", errs.ErrBadMatch},
		{".{_a_:x}", ruleset, nil, "", errs.ErrBadMatch},
		{"{._a_:x}", ruleset, nil, "", errs.ErrBadMatch},
		{"{_.a_:x}", ruleset, nil, "", errs.ErrBadMatch},
		{"{_a._:x}", ruleset, nil, "", errs.ErrBadMatch},
		{"{_a_.:x}", ruleset, nil, "", errs.ErrBadMatch},
		{"{_a_:.x}", ruleset, nil, "", errs.ErrBadMatch},
		{"{_a_:x.}", ruleset, nil, "", errs.ErrBadMatch},
		{"{_a_:x}.", ruleset, ss("{", "_a_", ":", "x", "}"), ".", nil},

		{" {_a_:x}", ruleset, ss("{", "_a_", ":", "x", "}"), "", nil},
		{"{ _a_:x}", ruleset, ss("{", "_a_", ":", "x", "}"), "", nil},
		{"{_ a_:x}", ruleset, ss("{", "_a_", ":", "x", "}"), "", nil},
		{"{_a_ :x}", ruleset, ss("{", "_a_", ":", "x", "}"), "", nil},
		{"{_a_: x}", ruleset, ss("{", "_a_", ":", "x", "}"), "", nil},
		{"{_a_:x }", ruleset, ss("{", "_a_", ":", "x", "}"), "", nil},
		{"{_a_:x} ", ruleset, ss("{", "_a_", ":", "x", "}"), "", nil},
		{"{_a_:x} .", ruleset, ss("{", "_a_", ":", "x", "}"), ".", nil},
	}

	ruleset, err = parse.NewRulesetFromStr(
		"test special vals",
		"c,, #1 | rnum, w| | c,, #2",
	)
	if err != nil {
		t.Fatalf("ruleset creation failed: %v", err)
	}

	rstests[ruleset] = []rulesettest{
		{".", ruleset, ss(), "", errs.ErrBadMatch},

		{",|123|,,", ruleset, ss(",", "|123|", ",,"), "", nil},
	}

	ruleset, err = parse.NewRulesetFromStr(
		"alpha comma",
		"ralpha, #1 | c,, #1 | ralpha, #1 | c,, #1 | ralpha, #1 | c,, #1",
	)
	if err != nil {
		t.Fatalf("ruleset creation failed: %v", err)
	}

	rstests[ruleset] = []rulesettest{
		{".", ruleset, ss(), "", errs.ErrBadMatch},
		{"a, b, c,", ruleset, ss("a", ",", "b", ",", "c", ","), "", nil},
	}

	ruleset, err = parse.NewRulesetFromStr("alpha3",
		"ralpha, #1 | ralpha, #1 | ralpha, #1",
	)
	if err != nil {
		t.Fatalf("ruleset creation failed: %v", err)
	}

	rstests[ruleset] = []rulesettest{
		{".", ruleset, ss(), "", errs.ErrBadMatch},
		{"abc", ruleset, ss("a", "b", "c"), "", nil},
	}

	ruleset, err = parse.NewRulesetFromStr("capture",
		"ralpha | c:, g0 | rnum",
	)
	if err != nil {
		t.Fatalf("ruleset creation failed: %v", err)
	}

	rstests[ruleset] = []rulesettest{
		{"a:1", ruleset, ss("a", "1"), "", nil},
		{"abc:123", ruleset, ss("abc", "123"), "", nil},
		{"a::1", ruleset, ss("a", "1"), "", nil},
		{"a: :1", ruleset, ss("a", "1"), "", nil},

		{".a:1", ruleset, ss(), "", errs.ErrBadMatch},

		{".", ruleset, ss(), "", errs.ErrBadMatch},
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
	ss := func(s ...string) []string { return s }
	eq := func(a, b []string) bool {
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
		{"xyz", rs, ss("x", "y", "z"), "", nil},

		{"a.", rs, ss("a"), ".", nil},
		{"ab.", rs, ss("a", "b"), ".", nil},

		{" ", rs, ss(), "", errs.ErrBadMatch},
		{"   a   ", rs, ss("a"), "", nil},

		{".a", rs, ss(), "", errs.ErrBadMatch},
	}

	rs, err = mk("alpha comma", "ralpha, #1 | c,, #1")
	if err != nil {
		t.Fatalf("ruleset creation failed: %v", err)
	}
	rstests[rs] = []rulesettest{
		{"a,", rs, ss("a", ","), "", nil},
		{"a,b,", rs, ss("a", ",", "b", ","), "", nil},
		{"a,b", rs, ss("a", ","), "b", nil},
		{"a ,", rs, ss("a", ","), "", nil},

		{".", rs, ss(), "", errs.ErrBadMatch},
	}

	rs, err = mk("str: alpha",
		"ralpha, w\" | c:, #1 | ralpha",
	)
	if err != nil {
		t.Fatalf("ruleset creation failed: %v", err)
	}
	rstests[rs] = []rulesettest{
		{"\"key\": value", rs, ss("\"key\"", ":", "value"), "", nil},
		{"\"a\": b \"c\": d", rs, ss("\"a\"", ":", "b", "\"c\"", ":", "d"), "", nil},

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
