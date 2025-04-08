package syntax

import (
	"unicode"

	"github.com/stu-k/go/parser/errors"
)

var defaultRulemap = map[string]*Rule{
	"alpha": RuleAlpha,
	"num":   RuleNum,
}

var RuleAny = &Rule{
	name:    "any",
	count:   -1,
	check:   func(_ string) bool { return true },
	capture: true,
}

var RuleAlpha = RuleAny.Named("alpha").CheckChar(unicode.IsLetter)

var RuleNum = RuleAny.Named("num").CheckChar(unicode.IsNumber)

// Rule defines a set of variables to parse a token by
type Rule struct {
	name string

	// count is the exact count of characters
	// expected in the resulting token
	//
	// ex. count(3) for abcd -> abc
	//     count(3) for ab -> error
	count int

	// check is the fn used to validate if the characters
	// in a string are valid for the rule
	check func(string) bool

	// capture determines if the match should be returned
	capture bool
}

func NewRule() *Rule { return RuleAny.clone() }

func (a *Rule) clone() *Rule {
	return &Rule{
		name:    a.name,
		count:   a.count,
		check:   a.check,
		capture: a.capture,
	}
}

func (a *Rule) IsAny() bool {
	return a == RuleAny
}

func (a *Rule) Count(n int) *Rule {
	new := a.clone()
	new.count = n
	return new
}

func (a *Rule) Named(s string) *Rule {
	new := a.clone()
	new.name = s
	return new
}

func (a *Rule) Name() string {
	return a.name
}

func (a *Rule) CheckChar(fn func(rune) bool) *Rule {
	new := a.clone()
	new.check = func(s string) bool {
		return len(s) > 0 && fn(rune(s[0]))
	}
	return new
}

func (a *Rule) Capture(v bool) *Rule {
	new := a.clone()
	new.capture = v
	return new
}

func (a *Rule) Chars(s string) *Rule {
	new := a.clone()
	m := make(map[string]struct{})

	for _, v := range s {
		m[string(v)] = struct{}{}
	}

	check := func(str string) bool {
		_, ok := m[str]
		return ok
	}

	new.check = check
	return new
}

func (a *Rule) ToSeq() *Sequence {
	s := NewSequence(a.name, a)
	return s
}

func (a *Rule) Parse(s string) (*ParseResult, error) {
	return a.parseChar(s)
}

func (a *Rule) parseChar(s string) (*ParseResult, error) {
	useCount := a.count >= 0
	if len(s) == 0 {
		return returnPr(a.name, s, errors.NewBadMatchErr(a.name, s))
	}

	var result string

	// add Rule.Start into result and iterate
	// past the value
	toparse := s

	var count int
	var r rune

	// handle checking results on end of string
	// or invalid character
	checkEnd := func(ct int) error {
		if len(result) == 0 {
			return errors.NewBadMatchErr(a.name, s)
		}

		if useCount && ct < a.count {
			return errors.NewBadMatchErr(a.name, s)
		}

		return nil
	}

	for i, c := range toparse {
		r = rune(c)

		countToUse := count

		if useCount {
			if countToUse >= a.count {
				if a.capture {
					return NewParseResult(a.name, []string{result}, s[i:]), nil
				}
				return NewParseResult(a.name, nil, s[i:]), nil
			}
		}

		ok := a.check(string(r))
		if !ok {
			err := checkEnd(countToUse)
			if err != nil {
				return returnPr(a.name, s, err)
			}

			if a.capture {
				return NewParseResult(a.name, []string{result}, s[i:]), nil
			}
			return NewParseResult(a.name, nil, s[i:]), nil
		}

		result += string(r)
		count++
	}

	err := checkEnd(count)
	if err != nil {
		return returnPr(a.name, s, err)
	}

	if a.capture {
		return NewParseResult(a.name, []string{result}, s[count:]), nil
	}
	return NewParseResult(a.name, nil, s[count:]), nil
}
