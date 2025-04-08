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
	repeat:  -1,
	check:   func(_ string) bool { return true },
	capture: true,
}

var RuleAlpha = RuleAny.Named("alpha").CheckChar(unicode.IsLetter)

var RuleNum = RuleAny.Named("num").CheckChar(unicode.IsNumber)

// Rule defines a set of variables to parse a token by
type Rule struct {
	name string

	// repeat is the exact repeat of characters
	// expected in the resulting token
	// -1 denotes an infinite repeat
	//
	// ex. repeat(3) for abcd -> abc
	//     repeat(3) for ab -> error
	repeat int

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
		repeat:  a.repeat,
		check:   a.check,
		capture: a.capture,
	}
}

func (a *Rule) IsAny() bool {
	return a == RuleAny
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

func (a *Rule) Repeat(n int) *Rule {
	r := a.clone()
	r.repeat = n
	return r
}

func (a *Rule) Parse(s string) (*ParseResult, error) {
	return a.parseChar(s)
}

func (a *Rule) parseChar(s string) (*ParseResult, error) {
	name, _ := a.name, a.repeat
	shouldRepeat := a.repeat >= 0
	if len(s) == 0 {
		return returnPr(name, s, errors.NewBadMatchErr(name, s))
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
			return errors.NewBadMatchErr(name, s)
		}

		if shouldRepeat && ct < a.repeat {
			return errors.NewBadMatchErr(name, s)
		}

		return nil
	}

	for i, c := range toparse {
		r = rune(c)

		countToUse := count

		if shouldRepeat {
			if countToUse >= a.repeat {
				if a.capture {
					return NewParseResult(name, []string{result}, s[i:]), nil
				}
				return NewParseResult(name, nil, s[i:]), nil
			}
		}

		ok := a.check(string(r))
		if !ok {
			err := checkEnd(countToUse)
			if err != nil {
				return returnPr(name, s, err)
			}

			if a.capture {
				return NewParseResult(name, []string{result}, s[i:]), nil
			}
			return NewParseResult(name, nil, s[i:]), nil
		}

		result += string(r)
		count++
	}

	err := checkEnd(count)
	if err != nil {
		return returnPr(name, s, err)
	}

	if a.capture {
		return NewParseResult(name, []string{result}, s[count:]), nil
	}
	return NewParseResult(name, nil, s[count:]), nil
}
