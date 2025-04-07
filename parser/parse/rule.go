package parse

import (
	"unicode"

	"github.com/stu-k/go/parser/errors"
)

type Parser interface {
	Parse(string) (string, string, error)
}

var blank = &Rule{
	name:        "blank",
	count:       -1,
	start:       0,
	end:         0,
	check:       func(_ rune) bool { return true },
	ignoreSpace: true,
	atLeastOne:  true,
}

var Alpha = blank.Name("alpha").Check(unicode.IsLetter)

// Rule defines a set of variables to parse a token by
type Rule struct {
	name string

	// count is the exact count of characters
	// expected in the resulting token
	//
	// ex. count(3) for abcd -> abc
	//     count(3) for ab -> error
	count int

	// start, end are the characters the token is
	// expected to be wrapped in to achieve tokenization
	start, end rune

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
}

func (a *Rule) clone() *Rule {
	return &Rule{
		name:        a.name,
		count:       a.count,
		start:       a.start,
		end:         a.end,
		check:       a.check,
		ignoreSpace: a.ignoreSpace,
		atLeastOne:  a.atLeastOne,
	}
}

func (a *Rule) Count(n int) *Rule {
	new := a.clone()
	new.count = n
	return new
}

func (a *Rule) Start(r rune) *Rule {
	new := a.clone()
	new.start = r
	return new
}

func (a *Rule) End(r rune) *Rule {
	new := a.clone()
	new.end = r
	return new
}

func (a *Rule) Wrap(r rune) *Rule {
	return a.
		clone().
		Start(r).
		End(r)
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

func (a *Rule) Parse(s string) (string, string, error) {
	useStart := a.start != 0
	useEnd := a.end != 0
	useCount := a.count >= 0

	// fail if empty string or doesn't start with
	// Rule.Start value if init
	if s == "" || (useStart && rune(s[0]) != a.start) {
		return "", "", errors.NewBadMatchErr(a.name, s)
	}

	var result string

	// add Rule.Start into result and iterate
	// past the value
	toparse := s
	if a.start != 0 {
		result += string(s[0])
		toparse = s[1:]
	}

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

		if useEnd && r != a.end {
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

		if useCount && !useEnd {
			if countToUse >= a.count {
				return result, s[i:], nil
			}
		}

		if !useCount && useEnd {
			if r == a.end {
				result += string(r)
				return result, s[i+2:], nil
			}
		}

		if useCount && useEnd {
			if countToUse == a.count+1 {
				if r == a.end {
					result += string(r)
					return result, s[i+1:], nil
				}
				return "", "", errors.NewBadMatchErr(a.name, s)
			}
		}

		ok := a.check(r)
		if !ok {
			err := checkEnd(countToUse)
			if err != nil {
				return "", "", err
			}

			return result, s[i:], nil
		}

		result += string(r)
		count++
	}

	err := checkEnd(count)
	if err != nil {
		return "", "", err
	}

	return result, s[count:], nil
}
