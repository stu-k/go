package parse

import (
	"github.com/stu-k/go/parser/errors"
)

type Parsable interface {
	Parse(string) ([]string, string, error)
}

type Ruleset struct {
	name string
	list []Parsable
}

func NewRuleset(name string, rules ...Parsable) *Ruleset {
	return &Ruleset{name, rules}
}

func (r *Ruleset) Len() int                     { return len(r.list) }
func (r *Ruleset) Add(rules ...Parsable)        { r.list = append(r.list, rules...) }
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
