package parse

import (
	"fmt"
	"unicode"

	"github.com/stu-k/go/parser/errors"
)

var blank = &Rule{
	name:        "blank",
	count:       -1,
	check:       func(_ rune) bool { return true },
	ignoreSpace: true,
	atLeastOne:  true,
	capture:     true,
}

var Alpha = blank.Name("alpha").Check(unicode.IsLetter)

var Numeric = blank.Name("numeric").Check(unicode.IsNumber)

var rulemap = map[string]*Rule{
	"alpha": Alpha,
	"num":   Numeric,
}

func FromChar(c rune) *Rule {
	return blank.
		Name(fmt.Sprintf("char %s", string(c))).
		Check(func(r rune) bool { return r == c }).
		Count(1)
}

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

	// ignoreSpace determines if spaces will be
	// omitted in tokenization
	//
	// TODO: remove in favor of Check fn
	ignoreSpace bool

	// atLeastOne determines if a token must have
	// at least one valid character
	atLeastOne bool

	// capture determines if the match should be returned
	capture bool
}

func NewRule() *Rule { return blank.clone() }

func (a *Rule) clone() *Rule {
	return &Rule{
		name:        a.name,
		count:       a.count,
		check:       a.check,
		ignoreSpace: a.ignoreSpace,
		atLeastOne:  a.atLeastOne,
		capture:     a.capture,
	}
}

func (a *Rule) IsAny() bool {
	return a == blank
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

func (a *Rule) IgnoreSpace(v bool) *Rule {
	new := a.clone()
	new.ignoreSpace = v
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
	var ignoredSpaces int
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

		if a.ignoreSpace && unicode.IsSpace(r) {
			count++
			ignoredSpaces++
			continue
		}

		countToUse := count

		// ignored spaces shouldn't be counted
		// as chars from Rule.Count
		if a.ignoreSpace {
			countToUse -= ignoredSpaces
		}

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
