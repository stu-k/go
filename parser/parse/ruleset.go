package parse

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/stu-k/go/parser/errors"
)

type Ruleset struct {
	name string
	list []*Rule
}

func NewRuleset(name string, rules ...*Rule) *Ruleset {
	return &Ruleset{name, rules}
}

func (r *Ruleset) Len() int           { return len(r.list) }
func (r *Ruleset) Add(rules ...*Rule) { r.list = append(r.list, rules...) }
func (r *Ruleset) Name() string       { return r.name }
func (r *Ruleset) Parse(s string) ([]string, string, error) {
	if len(r.list) == 0 || s == "" {
		return nil, "", errors.NewBadMatchErr(r.name, s)
	}

	var results []string
	toparse := s
	for _, rule := range r.list {
		result, rest, err := rule.Parse(toparse)
		if err != nil {
			return nil, "", err
		}

		if result == "" {
			return nil, "", errors.NewBadMatchErr(r.name, s)
		}

		results = append(results, result)
		toparse = rest
	}

	return results, toparse, nil
}

type rulesetargs struct {
	rule          *Rule
	char, w, s, e rune
	ct            int
}

func NewRulesetFromStr(name, s string) (*Ruleset, error) {
	errFn := func(arg string, i, j int) error {
		return fmt.Errorf(
			"error creating ruleset: invalid arg \"%v\" in segment %d, arg %d",
			arg, i, j,
		)
	}

	rs := NewRuleset(name)

	parts := strings.Split(s, "|")
	if len(parts) == 0 {
		return nil, fmt.Errorf("error creating ruleset: invalid string \"%s\"", s)
	}

	for i, part := range parts {

		var rsa rulesetargs

		args := strings.Split(part, ",")
		for j, arg := range args {
			if len(arg) == 0 || len(arg) == 1 {
				return nil, errFn(arg, i, j)
			}

			r := rune(arg[1])
			if arg[1:] == "comma" {
				r = ','
			}

			switch rune(arg[0]) {

			case 'r':
				rul, ok := rulemap[arg[1:]]
				if !ok {
					return nil, errFn(arg, i, j)
				}
				rsa.rule = rul
				continue

			case '#':
				ct, err := strconv.Atoi(arg[1:])
				if err != nil {
					return nil, errFn(arg, i, j)
				}
				rsa.ct = ct
				continue

			case 's':
				rsa.s = r
				continue

			case 'e':
				rsa.e = r
				continue

			case 'w':
				rsa.w = r
				continue

			case 'c':
				rsa.char = r
				continue

			default:
				return nil, errFn(arg, i, j)

			}
		}

		rule := NewRule()
		if rsa.rule != nil {
			rule = rsa.rule
		}
		if rsa.s != 0 {
			rule = rule.Start(rsa.s)
		}
		if rsa.e != 0 {
			rule = rule.End(rsa.e)
		}
		if rsa.w != 0 {
			rule = rule.Wrap(rsa.w)
		}
		if rsa.ct != 0 {
			rule = rule.Count(rsa.ct)
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
