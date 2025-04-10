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
		stx.NewRule(":").Chars(":"),
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
		stx.NewRule("{").Chars("{"),
		stx.RuleAlpha,
		stx.NewRule(":").Chars(":"),
		stx.RuleAlpha,
		stx.NewRule("}").Chars("}"),
	)
	sqTests[sequence] = []seqtest{
		{"{a:x}", sequence, ss("{", "a", ":", "x", "}"), "", nil},
		{"{abc:xyz}", sequence, ss("{", "abc", ":", "xyz", "}"), "", nil},
		{".", sequence, ss(), "", errs.ErrBadMatch},
		{"{a:x}...", sequence, ss("{", "a", ":", "x", "}"), "...", nil},
	}

	sequence = stx.NewSequence(
		"obj(kv(_var_ var))",
		stx.NewRule("").Chars("{"),
		stx.NewRule("").Chars("_"),
		stx.RuleAlpha,
		stx.NewRule("").Chars("_"),
		stx.NewRule("").Chars(":"),
		stx.RuleAlpha,
		stx.NewRule("").Chars("}"),
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
	mk := func(s string, r ...string) (*stx.Sequence, error) {
		rfs, err := stx.NewSequenceFromStrs(s, r...)
		if err != nil {
			return nil, err
		}
		return rfs, nil
	}

	type seqTest struct {
		in   string
		want []string
		rest string
		err  error
	}

	sqtests := make(map[*stx.Sequence][]seqTest)

	sq, err := mk(
		"alpha 1",
		"ralpha, #1",
	)
	if err != nil {
		t.Fatalf("ruleset creation failed: %v", err)
	}
	sqtests[sq] = []seqTest{
		{"a", ss("a"), "", nil},
		{"ab", ss("a", "b"), "", nil},
		{"abc", ss("a", "b", "c"), "", nil},

		{"a.", ss("a"), ".", nil},
		{"ab.", ss("a", "b"), ".", nil},

		{".", ss(), "", errs.ErrBadMatch},
	}

	sq, err = mk(
		"alpha comma",
		"ralpha, #1", "c,, #1",
	)
	if err != nil {
		t.Fatalf("ruleset creation failed: %v", err)
	}
	sqtests[sq] = []seqTest{
		{"a,", ss("a", ","), "", nil},
		{"a,b,", ss("a", ",", "b", ","), "", nil},
		{"a,b", ss("a", ","), "b", nil},
		{"a,", ss("a", ","), "", nil},

		{".", ss(), "", errs.ErrBadMatch},
	}

	sq, err = mk(
		"str: num",
		"ralpha", "c:, #1", "rnum",
	)
	if err != nil {
		t.Fatalf("ruleset creation failed: %v", err)
	}
	sqtests[sq] = []seqTest{
		{"abc:123", ss("abc", ":", "123"), "", nil},
		{"a:1b:2", ss("a", ":", "1", "b", ":", "2"), "", nil},
		{"a:1", ss("a", ":", "1"), "", nil},
		{"a:1b:", ss("a", ":", "1"), "b:", nil},

		{".", ss(), "", errs.ErrBadMatch},
	}

	for rs, tests := range sqtests {
		t.Run(rs.Name(), func(t *testing.T) {
			for _, test := range tests {
				got, err := rs.UntilFail().Parse(test.in)
				if !errors.Is(err, test.err) {
					t.Errorf("for \"%v\" expected error \"%v\"; got \"%v\"", test.in, test.err, err)
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

func TestSeqAnyOf(t *testing.T) {
	mk := func(s string, r ...string) (stx.Parsable, error) {
		sfs, err := stx.NewSequenceFromStrs(s, r...)
		if err != nil {
			return nil, err
		}
		return sfs, nil
	}

	type sqTest struct {
		in   string
		want []string
		err  error
	}

	sqTests := make(map[*stx.Sequence][]sqTest)

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

	sq := stx.NewSequence("a0.", sq1, sq2, sq3, sq4)
	sqTests[sq] = []sqTest{
		{"a", ss("alp", "1a"), nil},
		{"aa", ss("alp", "1a", "2alp"), nil},
		{"aaa", ss("alp", "1a", "2alp"), nil},

		{"b", ss("alp"), nil}, // passes
		{"bb", ss("alp", "2b", "2alp"), nil},
		{"bbb", ss("alp", "2b", "2alp"), nil},

		{"x", ss("alp"), nil}, // passes
		{"xy", ss("alp", "2alp"), nil},

		{"", ss(), errs.ErrBadMatch},  // passes
		{".", ss(), errs.ErrBadMatch}, // passes
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

	sq = stx.NewSequence("str / var", sq1, sq2)
	sqTests[sq] = []sqTest{
		{"'str'", ss("str"), nil},
		{"var", ss("var"), nil},

		{"", ss(), errs.ErrBadMatch},
		{".", ss(), errs.ErrBadMatch},
	}

	for sq, tests := range sqTests {
		t.Run(sq.Name(), func(t *testing.T) {
			for _, test := range tests {
				got, err := sq.AnyOf().Parse(test.in)
				if !errors.Is(err, test.err) {
					t.Fatalf("for \"%v\" expected error \"%v\"; got \"%v\"", test.in, test.err, err)
				}

				for _, v := range test.want {
					if !got.HasResult(v) {
						t.Errorf("for \"%v\" expected key \"%v\" in nameMap %v", test.in, v, got.NameMap())
					}
				}
			}
		})
	}
}

func TestSeqFromStrs(t *testing.T) {
	type seqTest struct {
		in   string
		rs   *stx.Sequence
		want []string
		rest string
		err  error
	}

	sqTests := make(map[*stx.Sequence][]seqTest)

	sq, err := stx.NewSequenceFromStrs(
		"alpha",
		"ralpha")
	if err != nil {
		t.Fatalf("sequence creation failed: %v", err)
	}

	sqTests[sq] = []seqTest{
		{"abc", sq, ss("abc"), "", nil},
		{"ab.c", sq, ss("ab"), ".c", nil},
		{"abc.", sq, ss("abc"), ".", nil},
		{".", sq, ss(), "", errs.ErrBadMatch},
	}

	sq, err = stx.NewSequenceFromStrs(
		"num",
		"rnum",
	)
	if err != nil {
		t.Fatalf("sequence creation failed: %v", err)
	}

	sqTests[sq] = []seqTest{
		{"123", sq, ss("123"), "", nil},
		{"12.3", sq, ss("12"), ".3", nil},
		{"123.", sq, ss("123"), ".", nil},
	}

	sq, err = stx.NewSequenceFromStrs(
		"alphanum",
		"ralpha", "rnum",
	)
	if err != nil {
		t.Fatalf("sequence creation failed: %v", err)
	}

	sqTests[sq] = []seqTest{
		{"a1", sq, ss("a", "1"), "", nil},
		{"abc123", sq, ss("abc", "123"), "", nil},
		{"a1.", sq, ss("a", "1"), ".", nil},
		{"a", sq, ss(), "", errs.ErrBadMatch},
		{"1", sq, ss(), "", errs.ErrBadMatch},
		{"a.1", sq, ss(), "", errs.ErrBadMatch},
	}

	sq, err = stx.NewSequenceFromStrs(
		"kv(var var)",
		"ralpha", "c:, #1", "rnum",
	)
	if err != nil {
		t.Fatalf("sequence creation failed: %v", err)
	}

	sqTests[sq] = []seqTest{
		{"a:1", sq, ss("a", ":", "1"), "", nil},
		{"abc:123", sq, ss("abc", ":", "123"), "", nil},

		{"a::1", sq, nil, "", errs.ErrBadMatch},
		{":1", sq, nil, "", errs.ErrBadMatch},
		{"a1", sq, nil, "", errs.ErrBadMatch},
		{"a:", sq, nil, "", errs.ErrBadMatch},

		{".a:1", sq, nil, "", errs.ErrBadMatch},
		{"a.:1", sq, nil, "", errs.ErrBadMatch},
		{"a:.1", sq, nil, "", errs.ErrBadMatch},
		{"a:1.", sq, ss("a", ":", "1"), ".", nil},
	}

	sq, err = stx.NewSequenceFromStrs(
		"obj(kv(var var))",
		"c{, #1", "ralpha", "c:, #1", "ralpha", "c}, #1",
	)
	if err != nil {
		t.Fatalf("sequence creation failed: %v", err)
	}

	sqTests[sq] = []seqTest{
		{"{a:x}", sq, ss("{", "a", ":", "x", "}"), "", nil},
		{"{abc:xyz}", sq, ss("{", "abc", ":", "xyz", "}"), "", nil},

		{".", sq, ss(), "", errs.ErrBadMatch},
		{".{a:x}", sq, ss(), "", errs.ErrBadMatch},
		{"{.a:x}", sq, ss(), "", errs.ErrBadMatch},
		{"{a.:x}", sq, ss(), "", errs.ErrBadMatch},
		{"{a:.x}", sq, ss(), "", errs.ErrBadMatch},
		{"{a:x.}", sq, ss(), "", errs.ErrBadMatch},
		{"{a:x}.", sq, ss("{", "a", ":", "x", "}"), ".", nil},
	}

	sq, err = stx.NewSequenceFromStrs(
		"obj(kv(_var_ var))",
		"c{, #1", "c_, #1", "ralpha", "c_, #1", "c:, #1", "ralpha", "c}, #1",
	)
	if err != nil {
		t.Fatalf("sequence creation failed: %v", err)
	}

	sqTests[sq] = []seqTest{
		{"{_a_:x}", sq, ss("{", "_", "a", "_", ":", "x", "}"), "", nil},
		{"{_abc_:xyz}", sq, ss("{", "_", "abc", "_", ":", "xyz", "}"), "", nil},

		{".", sq, ss(), "", errs.ErrBadMatch},
	}

	sq, err = stx.NewSequenceFromStrs(
		"test special vals",
		"c,, #1", "rnum", "c|, #1",
	)
	if err != nil {
		t.Fatalf("sequence creation failed: %v", err)
	}

	sqTests[sq] = []seqTest{
		{".", sq, ss(), "", errs.ErrBadMatch},

		{",1|", sq, ss(",", "1", "|"), "", nil},
	}

	sq, err = stx.NewSequenceFromStrs(
		"alpha comma",
		"ralpha, #1", "c,, #1", "ralpha, #1", "c,, #1", "ralpha, #1", "c,, #1",
	)
	if err != nil {
		t.Fatalf("sequence creation failed: %v", err)
	}

	sqTests[sq] = []seqTest{
		{".", sq, ss(), "", errs.ErrBadMatch},
		{"a,b,c,", sq, ss("a", ",", "b", ",", "c", ","), "", nil},
	}

	sq, err = stx.NewSequenceFromStrs(
		"alpha3",
		"ralpha, #1", "ralpha, #1", "ralpha, #1",
	)
	if err != nil {
		t.Fatalf("sequence creation failed: %v", err)
	}

	sqTests[sq] = []seqTest{
		{".", sq, ss(), "", errs.ErrBadMatch},
		{"abc", sq, ss("a", "b", "c"), "", nil},
	}

	sq, err = stx.NewSequenceFromStrs(
		"capture",
		"ralpha", "c:, g0", "rnum",
	)
	if err != nil {
		t.Fatalf("sequence creation failed: %v", err)
	}

	sqTests[sq] = []seqTest{
		{"a:1", sq, ss("a", "1"), "", nil},
		{"abc:123", sq, ss("abc", "123"), "", nil},
		{"a::1", sq, ss("a", "1"), "", nil},

		{".", sq, ss(), "", errs.ErrBadMatch},
	}

	for rs, tests := range sqTests {
		t.Run(rs.Name(), func(t *testing.T) {
			for _, test := range tests {
				got, err := test.rs.Parse(test.in)
				if !errors.Is(err, test.err) {
					t.Fatalf("for \"%v\" expected error \"%v\"; got \"%v\"", test.in, test.err, err)
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

func Test_Integration(t *testing.T) {
	type testcase struct {
		in, rest string
		want     []string
		seqs     []stx.Parsable
		err      error
	}

	tests := make(map[string]testcase)

	tests["parens"] = testcase{
		"()", "",
		ss("(", ")"),
		[]stx.Parsable{
			stx.NewRule("lparen").Chars("("),
			stx.NewRule("rparen").Chars(")"),
		},
		nil,
	}

	tests["alpha in parens"] = testcase{
		"(abc)", "",
		ss("abc"),
		[]stx.Parsable{
			stx.NewRule("lparen").Chars("(").Capture(false),
			stx.RuleAlpha.Named("alpha"),
			stx.NewRule("rparen").Chars(")").Capture(false),
		},
		nil,
	}

	tests["alpha comma alpha in parens"] = testcase{
		"(abc,xyz)", "",
		ss("abc", "xyz"),
		[]stx.Parsable{
			stx.NewRule("lparen").Chars("(").Capture(false),
			stx.RuleAlpha.Named("alpha1"),
			stx.NewRule("comma").Chars(",").Capture(false),
			stx.RuleAlpha.Named("alpha2"),
			stx.NewRule("rparen").Chars(")").Capture(false),
		},
		nil,
	}

	tests["alpha comma repeat in parens"] = testcase{
		"(foo,bar,baz,)", "",
		ss("foo", "bar", "baz"),
		[]stx.Parsable{
			stx.NewRule("lparen").Chars("(").Capture(false),
			stx.NewSequence("alpha comma",
				stx.RuleAlpha.Named("alpha"),
				stx.NewRule("comma").Chars(",").Capture(false),
			).UntilFail(),
			stx.NewRule("rparen").Chars(")").Capture(false),
		},
		nil,
	}

	tests["quoted alpha comma repeat in parens"] = testcase{
		"('foo','bar','baz',)", "",
		ss("foo", "bar", "baz"),
		[]stx.Parsable{
			stx.NewRule("lparen").Chars("(").Capture(false),
			stx.NewSequence("alpha comma",
				stx.NewRule("apos").Chars("'").Capture(false),
				stx.RuleAlpha.Named("alpha"),
				stx.NewRule("apos").Chars("'").Capture(false),
				stx.NewRule("comma").Chars(",").Capture(false),
			).UntilFail(),
			stx.NewRule("rparen").Chars(")").Capture(false),
		},
		nil,
	}

	keyval := stx.NewSequence("k/'k':n",
		stx.NewSequence("'v'/v",
			stx.NewSequence("'al'",
				stx.NewRule("apos").Chars("'").Capture(false),
				stx.RuleAlpha.Named("alpha"),
				stx.NewRule("apos").Chars("'").Capture(false),
			),
			stx.NewSequence("al",
				stx.RuleAlpha.Named("alpha"),
			),
		).PickOne(),
		stx.NewRule("colon").Chars(":").Capture(false),
		stx.RuleNum.Named("num"),
	)
	keyvalcomma := keyval.Named("k/'k':n,").With(
		stx.NewRule("comma").Chars(",").Capture(false),
	)
	kvtuple := stx.NewSequence("kvtuple",
		stx.NewRule("lparen").Chars("(").Capture(false),
		keyvalcomma.UntilFail(),
		keyval,
		stx.NewRule("rparen").Chars(")").Capture(false),
	)
	tests["quoted alpha comma no trailing comma in parens"] = testcase{
		"('foo':1,bar:2,'baz':3,quux:4) ()", " ()",
		ss("foo", "1", "bar", "2", "baz", "3", "quux", "4"),
		[]stx.Parsable{kvtuple},
		nil,
	}

	for name, test := range tests {
		seq := stx.NewSequence(name, test.seqs...)
		tstFn(t, name, test.in, seq, test.want, test.err, test.rest)
		res, err := seq.Parse(test.in)
		if res == nil {
			t.Fatalf("[%v] result nil", name)
		}

		if !errors.Is(test.err, err) {
			t.Errorf("[%v] wanted error %v; got %v", name, test.err, err)
		}

		got := res.Strings()
		if !eq(got, test.want) {
			t.Errorf("[%v] wanted result %v; got %v", name, test.want, got)
		}

		rest := res.Rest()
		if rest != test.rest {
			t.Errorf("[%v] wanted rest %v; got %v", name, test.rest, rest)
		}
	}
}

var tstFn = func(t *testing.T, name string, in string, p stx.Parsable, want []string, err error, rest string) {
	t.Run(name, func(t *testing.T) {
		res, goterr := p.Parse(in)
		if res == nil {
			t.Fatalf("[%v] result nil", name)
		}

		if !errors.Is(err, goterr) {
			t.Errorf("[%v] wanted error %v; got %v", name, err, goterr)
		}

		got := res.Strings()
		if !eq(got, want) {
			t.Errorf("[%v] wanted result %v; got %v", name, want, got)
		}

		gotrest := res.Rest()
		if rest != gotrest {
			t.Errorf("[%v] wanted rest %v; got %v", name, rest, gotrest)
		}
	})
}

func Test_Integration2(t *testing.T) {
	tst := func(p stx.Parsable, in string, want []string, rest string, err error) {
		got, goterr := p.Parse(in)
		if !errors.Is(goterr, err) {
			t.Fatalf("[%v] expected error %v; got %v", p.Name(), err, goterr)
		}

		if err != nil {
			t.Fatalf("[%v] expected no error; got %v", p.Name(), err)
		}

		if got == nil {
			t.Fatalf("[%v] expected result not to be nil", p.Name())
		}

		if !eq(got.Strings(), want) {
			t.Errorf("[%v] expected result %v; got %v", p.Name(), want, got.Strings())
		}
	}

	apos := stx.NewRule("apos").Repeat(1).Chars("'").Capture(false)
	comma := stx.NewRule("comma").Repeat(1).Chars(",").Capture(false)
	lbracket := stx.NewRule("lbracket").Repeat(1).Chars("[").Capture(false)
	rbracket := stx.NewRule("rbracket").Repeat(1).Chars("]").Capture(false)

	num := stx.RuleNum
	alp := stx.RuleAlpha
	alphanum := stx.NewSequence("alphanum",
		alp, num,
	).PickOne()
	str := stx.NewSequence("quoted",
		apos, alphanum, apos,
	)

	anyVal := stx.NewSequence("anyval",
		alp, num, str,
	).PickOne()
	valComma := stx.NewSequence("anyval comma",
		anyVal, comma,
	)
	valCommaEnd := stx.NewSequence("anyval comma end optional",
		valComma.UntilFail(),
		anyVal,
	)

	arr := stx.NewSequence("arr",
		lbracket,
		valCommaEnd,
		rbracket,
	)

	t.Run("arr", func(t *testing.T) {
		tst(arr,
			"['abc',123,xyz]",
			ss("abc", "123", "xyz"),
			"", nil)
	})

	colon := stx.NewRule("colon").Repeat(1).Chars(":").Capture(false)
	kv := stx.NewSequence("kv",
		anyVal, colon, anyVal,
	)
	kvComma := kv.Named("kv comma").
		With(comma)
	kvRepeat := stx.NewSequence("kv(,)",
		kvComma.UntilFail(),
		kv,
	)

	lbrace := stx.NewRule("lbrace").Repeat(1).Chars("{").Capture(false)
	rbrace := stx.NewRule("rbrace").Repeat(1).Chars("}").Capture(false)
	object := stx.NewSequence("object",
		lbrace,
		kvRepeat,
		rbrace,
	)
	t.Run("object", func(t *testing.T) {
		tst(object,
			"{abc:123,'xyz':'321',zyx:cba}",
			ss("abc", "123", "xyz", "321", "zyx", "cba"),
			"", nil)
	})

}
