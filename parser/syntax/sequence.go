package syntax

import (
	"github.com/stu-k/go/parser/errors"
)

type Parsable interface {
	Parse(string) (*ParseResult, error)
	Name() string
}

type Sequence struct {
	name string
	list []Parsable
}

func NewSequence(name string, rules ...Parsable) *Sequence {
	return &Sequence{name, rules}
}

func (r *Sequence) Len() int                 { return len(r.list) }
func (r *Sequence) Add(rules ...Parsable)    { r.list = append(r.list, rules...) }
func (r *Sequence) Name() string             { return r.name }
func (r *Sequence) UntilFail() *seqUntilFail { return &seqUntilFail{r} }
func (r *Sequence) OneOf() *seqOneOf         { return &seqOneOf{r} }
func (r *Sequence) Parse(s string) (*ParseResult, error) {
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

type seqUntilFail struct {
	*Sequence
}

func (r *seqUntilFail) Parse(s string) (*ParseResult, error) {
	all := NewParseResult(r.name, nil, s)
	for {
		results, err := r.Sequence.Parse(all.Rest())
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
		return nil, errors.NewBadMatchErr(r.Sequence.name, s)
	}
	return all, nil
}

type seqOneOf struct {
	*Sequence
}

func (r *seqOneOf) Parse(s string) (*ParseResult, error) {
	all := NewParseResult(r.name, nil, s)
	for _, p := range r.Sequence.list {
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
		return nil, errors.NewBadMatchErr(r.Sequence.name, s)
	}

	return all, nil
}
