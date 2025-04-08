package syntax

import (
	"github.com/stu-k/go/parser/errors"
)

type Parsable interface {
	Parse(string) (*ParseResult, error)
	Name() string
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
func (r *Ruleset) UntilFail() *rulesetUntilFail { return &rulesetUntilFail{r} }
func (r *Ruleset) OneOf() *rulesetOneOf         { return &rulesetOneOf{r} }
func (r *Ruleset) Parse(s string) (*ParseResult, error) {
	if len(r.list) == 0 || s == "" {
		return nil, errors.NewBadMatchErr(r.name, s)
	}

	results := NewParseResult(r.name, nil, s)
	for _, rule := range r.list {
		result, err := rule.Parse(results.Rest())
		if err != nil {
			return nil, err
		}

		if result == nil {
			return nil, errors.NewBadMatchErr(r.name, s)
		}

		results.Append(result)
		results.SetRest(result.Rest())
	}

	return results, nil
}

type rulesetUntilFail struct {
	*Ruleset
}

func (r *rulesetUntilFail) Parse(s string) (*ParseResult, error) {
	all := NewParseResult(r.name, nil, s)
	for {
		results, err := r.Ruleset.Parse(all.Rest())
		if err != nil {
			if all.Len() == 0 {
				return nil, err
			}
			return all, nil
		}
		all.Append(results)
		all.SetRest(results.Rest())
		if results.Rest() == "" {
			break
		}
	}
	if all.Len() == 0 {
		return nil, errors.NewBadMatchErr(r.Ruleset.name, s)
	}
	return all, nil
}

type rulesetOneOf struct {
	*Ruleset
}

func (r *rulesetOneOf) Parse(s string) (*ParseResult, error) {
	all := NewParseResult(r.name, nil, s)
	for _, p := range r.Ruleset.list {
		results, err := p.Parse(all.Rest())
		if err != nil {
			continue
		}
		all.Append(results)
		if results.Rest() == "" {
			break
		}
	}

	if all.Len() == 0 {
		return nil, errors.NewBadMatchErr(r.Ruleset.name, s)
	}

	return all, nil
}
