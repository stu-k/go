package syntax

import (
	"fmt"
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
	a.repeat = n
	return a
}

func (a *Rule) Parse(s string) (*Result, error) {
	if a.checkChar == nil {
		// defaulting to "none" rule to invalidate null pointers
		fmt.Println("[ERR] DEFAULTED_NONE_RULE")
		return a.clone().
			Named("DEFAULTED_NONE_RULE").
			CheckChar(func(_ rune) bool { return false }).
			Parse(s)
	}

	return a.parseChar(s)
}

func (a *Rule) parseChar(s string) (*Result, error) {
	return charParser{
		a.name,
		a.repeat,
		a.capture,
		a.checkChar,
	}.Parse(s)
}

type charParser struct {
	name    string
	repeat  int
	capture bool
	check   func(string) bool
}

func (p charParser) Parse(s string) (*Result, error) {
	shouldRepeat := p.repeat > 0

	var result string
	var count int

	for _, c := range s {
		r := rune(c)

		if shouldRepeat && count >= p.repeat {
			break
		}

		ok := p.check(string(r))
		if !ok {
			break
		}

		result += string(r)
		count++
	}

	if len(result) == 0 {
		return retErr(p.name, errors.NewBadMatchErr(p.name, s, "parsechar:nores"))
	}

	if shouldRepeat && count < p.repeat {
		return retErr(p.name, errors.NewBadMatchErr(p.name, s, "parsechar:toofew"))
	}

	if !p.capture {
		return NewResult(p.name, nil, s[count:]), nil
	}

	return NewResult(p.name, []string{result}, s[count:]), nil

}
