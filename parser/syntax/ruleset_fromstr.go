package syntax

import (
	"fmt"
	"strconv"
	"strings"
)

type rulesetargs struct {
	rule        *Rule
	char        rune
	count       int
	cap, usecap bool
}

func NewRulesetFromStrs(name string, parts ...string) (*Ruleset, error) {
	return newRulesetFromStrs(name, defaultRulemap, parts...)
}

func newRulesetFromStrs(name string, pmap map[string]*Rule, parts ...string) (*Ruleset, error) {
	errFn := func(arg string, i, j int) error {
		return fmt.Errorf(
			"error creating ruleset: invalid arg \"%v\" in segment %d, arg %d",
			arg, i, j,
		)
	}

	rs := NewRuleset(name)

	if len(parts) == 0 {
		return nil, fmt.Errorf("error creating ruleset: invalid string \"%s\"", parts)
	}

	for i, part := range parts {

		var rsa rulesetargs

		args := strings.Split(part, ", ")
		for j, arg := range args {
			if len(arg) == 0 || len(arg) == 1 {
				return nil, errFn(arg, i, j)
			}

			r := rune(arg[1])
			switch rune(arg[0]) {

			case 'r':
				rul, ok := pmap[arg[1:]]
				if !ok || rul == nil {
					return nil, errFn(arg, i, j)
				}
				rsa.rule = rul
				continue

			case '#':
				ct, err := strconv.Atoi(arg[1:])
				if err != nil {
					return nil, errFn(arg, i, j)
				}
				rsa.count = ct
				continue

			case 'c':
				rsa.char = r
				continue

			case 'g':
				if r == '0' {
					rsa.cap = false
					rsa.usecap = true
					continue
				}
				if r == '1' {
					rsa.cap = true
					rsa.usecap = true
					continue
				}
				return nil, errFn(arg, i, j)

			default:
				return nil, errFn(arg, i, j)

			}
		}

		rule := NewRule()
		if rsa.rule != nil {
			rule = rsa.rule
		}
		if rsa.count != 0 {
			rule = rule.Count(rsa.count)
		}
		if rsa.char != 0 {
			rule = rule.Check(func(r rune) bool { return r == rsa.char })
		}
		if rsa.usecap {
			rule = rule.Capture(rsa.cap)
		}
		rule = rule.Name(part)
		if rule.IsAny() {
			return nil, fmt.Errorf("error creating ruleset: can't add empty rule \"%s\"", part)
		}
		rs.Add(rule)
	}

	if rs.Len() == 0 {
		return nil, fmt.Errorf("error creating ruleset: can't use empty ruleset")
	}

	return rs, nil
}
