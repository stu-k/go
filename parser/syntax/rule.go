package syntax

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/stu-k/go/parser/errors"
)

var defaultRulemap = map[string]*Rule{
	"alpha": RuleAlpha,
	"num":   RuleNum,
}

var ruleAny = &Rule{
	name:      "DEFAULT_ANY_RULE",
	repeat:    -1,
	capture:   true,
	checkChar: func(_ string) bool { return true },
	checkStr:  "",
	modified:  false,
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

	// checkChar is the fn used to validate if the characters
	// in a string are valid for the rule
	checkChar func(string) bool

	// checkStr is the exact string for a rule to parse against
	checkStr string

	// modified informs if the rule has been changed from the default
	modified bool
}

func NewRule(n string) *Rule { return ruleAny.clone().Named(n) }

func (a *Rule) clone() *Rule {
	new := &Rule{
		name:     a.name,
		capture:  a.capture,
		checkStr: a.checkStr,
		modified: true,
	}
	new.repeat = a.repeat
	new.checkChar = a.checkChar
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
	a.checkChar = func(s string) bool {
		return len(s) > 0 && fn(rune(s[0]))
	}
	a.checkStr = ""
	return a
}

func (a *Rule) CheckStr(s string) *Rule {
	a.checkStr = s
	a.checkChar = nil
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

	a.checkChar = check
	return a
}

func (a *Rule) Seq() *Sequence {
	return NewSequence(a.name, a)
}

func (a *Rule) Repeat(n int) *Rule {
	new := a.clone()
	new.repeat = n
	return new
}

func (a *Rule) Parse(s string) (*ParseResult, error) {
	if a.checkChar == nil && len(a.checkStr) == 0 {
		// defaulting to "none" rule to invalidate null pointers
		fmt.Println("[ERR] DEFAULTED_NONE_RULE")
		return a.clone().
			Named("DEFAULTED_NONE_RULE").
			CheckChar(func(_ rune) bool { return false }).
			Parse(s)
	}

	if len(a.checkStr) > 0 {
		if a.repeat > -1 {
			return a.parseStrRepeat(a.checkStr, s, a.repeat)
		}
		return a.parseStr(a.checkStr, s)
	}

	return a.parseChar(s)
}

func (a *Rule) parseStr(match, s string) (*ParseResult, error) {
	if len(s) == 0 || len(match) == 0 {
		return retErr(a.name, errors.NewBadMatchErr(a.name, s, "parsestr:emptystr"))
	}

	if !strings.HasPrefix(s, match) {
		return retErr(a.name, errors.NewBadMatchErr(a.name, s, "parsestr:noprefix"))
	}

	return NewParseResult(a.name, []string{match}, s[len(match):]), nil
}

func (a *Rule) parseStrRepeat(match, s string, n int) (*ParseResult, error) {
	if len(s) == 0 || len(match) == 0 {
		return retErr(a.name, errors.NewBadMatchErr(a.name, s, "parsestrrepeat:emptystr"))
	}

	var count int
	results := NewParseResult(a.name, nil, s)
	for i := 0; i < n; i++ {
		result, err := a.parseStr(match, results.Rest())
		if err != nil {
			if results.Len() == 0 {
				return retErr(a.name, errors.NewBadMatchErr(a.name, s, "parsestrrepeat:nofirstmatch"))
			}
			break
		}
		results.Append(result)
		results.SetRest(result.Rest())
		count++
	}

	if results.Len() == 0 {
		return retErr(a.name, errors.NewBadMatchErr(a.name, s, "parsestrrepeat:emptyresult"))
	}
	if count < n {
		return retErr(a.name, errors.NewBadMatchErr(a.name, s, "parsestrrepeat:lowcount"))
	}

	return results, nil
}

func (a *Rule) parseChar(s string) (*ParseResult, error) {
	shouldRepeat := a.repeat >= 0
	if len(s) == 0 {
		return retErr(a.name, errors.NewBadMatchErr(a.name, s, "rule:parsechar:emptystr"))
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
			return errors.NewBadMatchErr(a.name, s, "parsechar:noresult")
		}

		if shouldRepeat && ct < a.repeat {
			return errors.NewBadMatchErr(a.name, s, "parsechar:lowcount")
		}

		return nil
	}

	for i, c := range toparse {
		r = rune(c)

		countToUse := count

		if shouldRepeat {
			if countToUse >= a.repeat {
				if a.capture {
					return NewParseResult(a.name, []string{result}, s[i:]), nil
				}
				return NewParseResult(a.name, nil, s[i:]), nil
			}
		}

		ok := a.checkChar(string(r))
		if !ok {
			err := checkEnd(countToUse)
			if err != nil {
				return retErr(a.name, err)
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
		return retErr(a.name, err)
	}

	if a.capture {
		return NewParseResult(a.name, []string{result}, s[count:]), nil
	}
	return NewParseResult(a.name, nil, s[count:]), nil
}
