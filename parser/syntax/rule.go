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
	name:       "any",
	count:      -1,
	check:      func(_ rune) bool { return true },
	atLeastOne: true,
	capture:    true,
}

var RuleAlpha = RuleAny.Name("alpha").Check(unicode.IsLetter)

var RuleNum = RuleAny.Name("num").Check(unicode.IsNumber)

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
	check func(rune) bool

	// atLeastOne determines if a token must have
	// at least one valid character
	atLeastOne bool

	// capture determines if the match should be returned
	capture bool
}

func NewRule() *Rule { return RuleAny.clone() }

func (a *Rule) clone() *Rule {
	return &Rule{
		name:       a.name,
		count:      a.count,
		check:      a.check,
		atLeastOne: a.atLeastOne,
		capture:    a.capture,
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

func (a *Rule) Name(s string) *Rule {
	new := a.clone()
	new.name = s
	return new
}

func (a *Rule) Check(fn func(rune) bool) *Rule {
	new := a.clone()
	new.check = fn
	return new
}

func (a *Rule) Capture(v bool) *Rule {
	new := a.clone()
	new.capture = v
	return new
}

func (a *Rule) Chars(s string) *Rule {
	new := a.clone()
	m := make(map[rune]struct{})

	for _, v := range s {
		m[rune(v)] = struct{}{}
	}

	check := func(r rune) bool {
		_, ok := m[r]
		return ok
	}

	new.check = check
	return new
}

func (a *Rule) Parse(s string) ([]string, string, error) {
	useCount := a.count >= 0
	if s == "" {
		return nil, "", errors.NewBadMatchErr(a.name, s)
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
		if a.atLeastOne && result == "" {
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
					return []string{result}, s[i:], nil
				}
				return []string{}, s[i:], nil
			}
		}

		ok := a.check(r)
		if !ok {
			err := checkEnd(countToUse)
			if err != nil {
				return nil, "", err
			}

			if a.capture {
				return []string{result}, s[i:], nil
			}
			return []string{}, s[i:], nil
		}

		result += string(r)
		count++
	}

	err := checkEnd(count)
	if err != nil {
		return nil, "", err
	}

	if a.capture {
		return []string{result}, s[count:], nil
	}
	return []string{}, s[count:], nil
}
