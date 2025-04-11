package syntax

import (
	"unicode"

	"github.com/stu-k/go/parser/errors"
)

var ruleAny = &Rule{
	name:     "DEFAULT_ANY_RULE",
	repeat:   -1,
	capture:  true,
	check:    func(_ string) bool { return true },
	modified: false,
}

var RuleAlpha = NewRule("alpha").CheckChar(unicode.IsLetter)

var RuleNum = NewRule("num").CheckChar(unicode.IsNumber)

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

	// capture determines if the match should be returned
	capture bool

	// check is the fn used to validate if the characters
	// in a string are valid for the rule
	check func(string) bool

	// modified informs if the rule has been changed from the default
	modified bool
}

func NewRule(n string) *Rule { return ruleAny.clone().Named(n) }

func (a *Rule) clone() *Rule {
	new := &Rule{
		name:     a.name,
		capture:  a.capture,
		modified: true,
	}
	new.repeat = a.repeat
	new.check = a.check
	return new
}

func (a *Rule) IsAny() bool {
	return !a.modified && a.name != "DEFAULT_ANY_RULE"
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
	a.check = func(s string) bool {
		return len(s) > 0 && fn(rune(s[0]))
	}
	return a
}

func (a *Rule) Capture(v bool) *Rule {
	a.capture = v
	return a
}

func (a *Rule) Chars(s string) *Rule {
	m := make(map[string]struct{})

	for _, v := range s {
		m[string(v)] = struct{}{}
	}

	check := func(str string) bool {
		_, ok := m[str]
		return ok
	}

	a.check = check
	return a
}

func (a *Rule) Seq() *Sequence {
	return NewSequence(a.name, a)
}

func (a *Rule) Repeat(n int) *Rule {
	a.repeat = n
	return a
}

func (a *Rule) Parse(s string) (*Result, error) {
	shouldRepeat := a.repeat > 0

	var result string
	var count int

	for _, c := range s {
		r := rune(c)

		if shouldRepeat && count >= a.repeat {
			break
		}

		ok := a.check(string(r))
		if !ok {
			break
		}

		result += string(r)
		count++
	}

	if len(result) == 0 {
		return retErr(a.name, errors.NewBadMatchErr(a.name, s, "parsechar:nores"))
	}

	if shouldRepeat && count < a.repeat {
		return retErr(a.name, errors.NewBadMatchErr(a.name, s, "parsechar:toofew"))
	}

	if !a.capture {
		return NewResult(a.name, nil, s[count:]), nil
	}

	return NewResult(a.name, []string{result}, s[count:]), nil

}
