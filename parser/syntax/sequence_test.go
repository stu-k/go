package syntax_test

import (
	"errors"
	"testing"

	errs "github.com/stu-k/go/parser/errors"
	stx "github.com/stu-k/go/parser/syntax"
)

func TestSequence(t *testing.T) {
	type seqtest struct {
		in   string
		rs   *stx.Sequence
		want []string
		rest string
		err  error
	}

	sqTests := make(map[*stx.Sequence][]seqtest)

	sequence := stx.NewSequence("alpha", stx.RuleAlpha)
	sqTests[sequence] = []seqtest{
		{"abc", sequence, ss("abc"), "", nil},
		{"ab.c", sequence, ss("ab"), ".c", nil},
		{"abc.", sequence, ss("abc"), ".", nil},
		{".", sequence, ss(), "", errs.ErrBadMatch},
	}

	sequence = stx.NewSequence("num", stx.RuleNum)
	sqTests[sequence] = []seqtest{
		{"123", sequence, ss("123"), "", nil},
		{"12.3", sequence, ss("12"), ".3", nil},
		{"123.", sequence, ss("123"), ".", nil},
	}

	sequence = stx.NewSequence(
		"alphanum",
		stx.RuleAlpha,
		stx.RuleNum,
	)
	sqTests[sequence] = []seqtest{
		{"a1", sequence, ss("a", "1"), "", nil},
		{"abc123", sequence, ss("abc", "123"), "", nil},
		{"a1.", sequence, ss("a", "1"), ".", nil},
		{"a", sequence, ss(), "", errs.ErrBadMatch},
		{"1", sequence, ss(), "", errs.ErrBadMatch},
		{"a.1", sequence, ss(), "", errs.ErrBadMatch},
	}

	sequence = stx.NewSequence(
		"kv(var var)",
		stx.RuleAlpha,
		stx.RuleAny.Chars(":"),
		stx.RuleNum,
	)
	sqTests[sequence] = []seqtest{
		{"a:1", sequence, ss("a", ":", "1"), "", nil},
		{"abc:123", sequence, ss("abc", ":", "123"), "", nil},
		{"a:1:", sequence, ss("a", ":", "1"), ":", nil},
		{".a:1", sequence, ss(), "", errs.ErrBadMatch},
		{"a.:1", sequence, ss(), "", errs.ErrBadMatch},
		{"a:.1", sequence, ss(), "", errs.ErrBadMatch},
		{"a:1.", sequence, ss("a", ":", "1"), ".", nil},
	}

	sequence = stx.NewSequence(
		"obj(kv(var var))",
		stx.RuleAny.Chars("{"),
		stx.RuleAlpha,
		stx.RuleAny.Chars(":"),
		stx.RuleAlpha,
		stx.RuleAny.Chars("}"),
	)
	sqTests[sequence] = []seqtest{
		{"{a:x}", sequence, ss("{", "a", ":", "x", "}"), "", nil},
		{"{abc:xyz}", sequence, ss("{", "abc", ":", "xyz", "}"), "", nil},
		{".", sequence, ss(), "", errs.ErrBadMatch},
		{"{a:x}...", sequence, ss("{", "a", ":", "x", "}"), "...", nil},
	}

	sequence = stx.NewSequence(
		"obj(kv(_var_ var))",
		stx.RuleAny.Chars("{"),
		stx.RuleAny.Chars("_"),
		stx.RuleAlpha,
		stx.RuleAny.Chars("_"),
		stx.RuleAny.Chars(":"),
		stx.RuleAlpha,
		stx.RuleAny.Chars("}"),
	)
	sqTests[sequence] = []seqtest{
		{"{_a_:x}", sequence, ss("{", "_", "a", "_", ":", "x", "}"), "", nil},
		{"{_abc_:x}", sequence, ss("{", "_", "abc", "_", ":", "x", "}"), "", nil},
	}

	for rs, tests := range sqTests {
		t.Run(rs.Name(), func(t *testing.T) {
			for _, test := range tests {
				got, err := test.rs.Parse(test.in)
				if !errors.Is(err, test.err) {
					t.Fatalf("for \"%v\" expected error \"%v\"; got \"%v\"", test.in, test.err, err)
				}

				if err != nil {
					return
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

func TestSeqUntilFail(t *testing.T) {
	mk := func(s string, r ...string) (stx.Parsable, error) {
		rfs, err := stx.NewSequenceFromStrs(s, r...)
		if err != nil {
			return nil, err
		}
		return rfs.UntilFail(), nil
	}

	type seqTest struct {
		in   string
		rs   stx.Parsable
		want []string
		rest string
		err  error
	}

	sqtests := make(map[stx.Parsable][]seqTest)

	sq, err := mk(
		"alpha 1",
		"ralpha, #1",
	)
	if err != nil {
		t.Fatalf("ruleset creation failed: %v", err)
	}
	sqtests[sq] = []seqTest{
		{"a", sq, ss("a"), "", nil},
		{"ab", sq, ss("a", "b"), "", nil},
		{"abc", sq, ss("a", "b", "c"), "", nil},

		{"a.", sq, ss("a"), ".", nil},
		{"ab.", sq, ss("a", "b"), ".", nil},

		{".", sq, ss(), "", errs.ErrBadMatch},
	}

	sq, err = mk(
		"alpha comma",
		"ralpha, #1", "c,, #1",
	)
	if err != nil {
		t.Fatalf("ruleset creation failed: %v", err)
	}
	sqtests[sq] = []seqTest{
		{"a,", sq, ss("a", ","), "", nil},
		{"a,b,", sq, ss("a", ",", "b", ","), "", nil},
		{"a,b", sq, ss("a", ","), "b", nil},
		{"a,", sq, ss("a", ","), "", nil},

		{".", sq, ss(), "", errs.ErrBadMatch},
	}

	sq, err = mk(
		"str: num",
		"ralpha", "c:, #1", "rnum",
	)
	if err != nil {
		t.Fatalf("ruleset creation failed: %v", err)
	}
	sqtests[sq] = []seqTest{
		{"abc:123", sq, ss("abc", ":", "123"), "", nil},
		{"a:1b:2", sq, ss("a", ":", "1", "b", ":", "2"), "", nil},
		{"a:1", sq, ss("a", ":", "1"), "", nil},
		{"a:1b:", sq, ss("a", ":", "1"), "b:", nil},

		{".", sq, ss(), "", errs.ErrBadMatch},
	}

	for rs, tests := range sqtests {
		t.Run(rs.Name(), func(t *testing.T) {
			for _, test := range tests {
				got, err := test.rs.Parse(test.in)
				if !errors.Is(err, test.err) {
					t.Fatalf("for \"%v\" expected error \"%v\"; got \"%v\"", test.in, test.err, err)
				}

				if err != nil {
					return
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

func TestSeqOneOf(t *testing.T) {
	mk := func(s string, r ...string) (stx.Parsable, error) {
		sfs, err := stx.NewSequenceFromStrs(s, r...)
		if err != nil {
			return nil, err
		}
		return sfs.OneOf(), nil
	}

	type sqTest struct {
		in   string
		want []string
		err  error
	}

	sqTests := make(map[stx.Parsable][]sqTest)

	sq1, err := mk(
		"alp",
		"ralpha",
	)
	if err != nil {
		t.Fatalf("ruleset creation failed: %v", err)
	}

	sq2, err := mk(
		"1a",
		"ca, #1",
	)
	if err != nil {
		t.Fatalf("ruleset creation failed: %v", err)
	}

	sq3, err := mk(
		"2alp",
		"ralpha, #2",
	)
	if err != nil {
		t.Fatalf("ruleset creation failed: %v", err)
	}

	sq4, err := mk(
		"2b",
		"cb, #2",
	)
	if err != nil {
		t.Fatalf("ruleset creation failed: %v", err)
	}

	sq := stx.NewSequence("a0.", sq1, sq2, sq3, sq4).OneOf()
	sqTests[sq] = []sqTest{
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

	sq1, err = mk(
		"var",
		"ralpha",
	)
	if err != nil {
		t.Fatalf("ruleset creation failed: %v", err)
	}

	sq2, err = mk(
		"str",
		"c', #1", "ralpha", "c', #1",
	)
	if err != nil {
		t.Fatalf("ruleset creation failed: %v", err)
	}

	sq = stx.NewSequence("str / var", sq1, sq2).OneOf()
	sqTests[sq] = []sqTest{
		{"'str'", ss("str"), nil},
		{"var", ss("var"), nil},

		{"", ss(), errs.ErrBadMatch},
		{".", ss(), errs.ErrBadMatch},
	}

	for sq, tests := range sqTests {
		t.Run(sq.Name(), func(t *testing.T) {
			for _, test := range tests {
				got, err := sq.Parse(test.in)
				if !errors.Is(err, test.err) {
					t.Fatalf("for \"%v\" expected error \"%v\"; got \"%v\"", test.in, test.err, err)
				}

				if err != nil {
					return
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
