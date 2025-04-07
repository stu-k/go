package parse

import (
	"strings"

	"github.com/stu-k/go/parser/errors"
)

type Ruleset struct {
	name string
	list []*Rule
}

func NewRuleset(rules ...*Rule) *Ruleset {
	var name []string
	for _, r := range rules {
		name = append(name, r.name)
	}
	return &Ruleset{strings.Join(name, "|"), rules}
}

func (r *Ruleset) Name() string { return r.name }
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
