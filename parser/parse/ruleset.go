package parse

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/stu-k/go/parser/errors"
)

type Parser interface {
	Parse(string) ([]string, string, error)
}

type Ruleset struct {
	name string
	list []Parser
}

func NewRuleset(name string, rules ...Parser) *Ruleset {
	return &Ruleset{name, rules}
}

func (r *Ruleset) Len() int                     { return len(r.list) }
func (r *Ruleset) Add(rules ...Parser)          { r.list = append(r.list, rules...) }
func (r *Ruleset) Name() string                 { return r.name }
func (r *Ruleset) UntilFail() *RulesetUntilFail { return &RulesetUntilFail{r} }
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

		if result == nil {
			return nil, "", errors.NewBadMatchErr(r.name, s)
		}

		results = append(results, result...)
		toparse = rest
	}

	return results, toparse, nil
}

type RulesetUntilFail struct {
	*Ruleset
}

func (r *RulesetUntilFail) Parse(s string) ([]string, string, error) {
	var allResults []string

	toparse := s
	for {
		results, rest, err := r.Ruleset.Parse(toparse)
		if err != nil {
			if len(allResults) == 0 {
				return nil, "", err
			}
			return allResults, toparse, nil
		}
		allResults = append(allResults, results...)
		toparse = rest
		if rest == "" {
			break
		}
	}
	if len(allResults) == 0 {
		return nil, "", errors.NewBadMatchErr(r.Ruleset.name, s)
	}
	return allResults, toparse, nil
}

func NewRulesetUntilFail(name string, rules ...Parser) *RulesetUntilFail {
	return &RulesetUntilFail{NewRuleset(name, rules...)}
}

type rulesetargs struct {
	rule        *Rule
	char        rune
	count       int
	cap, usecap bool
}

func NewRulesetFromStr(name, s string) (*Ruleset, error) {
	errFn := func(arg string, i, j int) error {
		return fmt.Errorf(
			"error creating ruleset: invalid arg \"%v\" in segment %d, arg %d",
			arg, i, j,
		)
	}

	rs := NewRuleset(name)

	parts := strings.Split(s, " | ")
	if len(parts) == 0 {
		return nil, fmt.Errorf("error creating ruleset: invalid string \"%s\"", s)
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
