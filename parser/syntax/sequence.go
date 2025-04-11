package syntax

import (
	"github.com/stu-k/go/parser/errors"
)

type Sequencer interface {
	Seq() *Sequence
}

type Sequence struct {
	name    string
	list    []Parsable
	capture bool
}

func NewSequence(name string, rules ...Parsable) *Sequence {
	return &Sequence{name, rules, true}
}

func (r *Sequence) Seq() *Sequence {
	return r
}
func (r *Sequence) clone() *Sequence {
	return &Sequence{
		name:    r.name,
		list:    r.list[:],
		capture: r.capture,
	}
}

func (r *Sequence) With(p ...Parsable) *Sequence {
	new := r.clone()
	new.list = append(new.list, p...)
	return new
}
func (r *Sequence) Named(n string) *Sequence {
	new := r.clone()
	new.name = n
	return new
}
func (r *Sequence) Len() int          { return len(r.list) }
func (r *Sequence) Add(v ...Parsable) { r.list = append(r.list, v...) }
func (r *Sequence) Name() string      { return r.name }
func (r *Sequence) SetCapture(v bool) *Sequence {
	new := r.clone()
	new.capture = v
	return new
}
func (r *Sequence) Parse(s string) (*Result, error) {
	return Parse(r.name, s, r.list...)
}

func Parse(name, s string, pars ...Parsable) (*Result, error) {
	results := NewResult(name, nil, s)
	for _, p := range pars {
		result, err := p.Parse(results.Rest())
		if err != nil {
			return retErr(name, err)
		}

		results.Append(result)
		results.SetRest(result.Rest())
	}
	return results, nil
}

type untilfailparser struct {
	s *Sequence
}

func (p *untilfailparser) Parse(s string) (*Result, error) {
	return p.s.untilFail(s)
}

func (p *untilfailparser) Name() string {
	return p.s.Name()
}

func (r *Sequence) UntilFail() Parsable {
	return &untilfailparser{r}
}

func (r *Sequence) untilFail(s string) (*Result, error) {
	all := NewResult(r.name, nil, s)
	for {
		results, _ := r.Parse(all.Rest())
		if results.IsEmpy() {
			if len(all.Strings()) == 0 {
				return NewResult(r.name, nil, ""), nil
			}
			return all, nil
		}

		all.Append(results)
		all.SetRest(results.Rest())
		if len(results.Rest()) == 0 {
			break
		}
	}

	if !r.capture {
		return NewResult(r.name, nil, ""), nil
	}

	return all, nil
}

type anyofparser struct{ s *Sequence }

func (p *anyofparser) Parse(s string) (*Result, error) { return p.s.anyOf(s) }
func (p *anyofparser) Name() string                    { return p.s.Name() }

func (r *Sequence) AnyOf() Parsable {
	return &anyofparser{r}
}

func (r *Sequence) anyOf(s string) (*Result, error) {
	all := NewResult(r.name, nil, "")
	for _, p := range r.list {
		results, _ := p.Parse(s)
		if results.IsEmpy() {
			continue
		}
		all.Append(results)
	}

	if all.IsEmpy() {
		return retErr(r.name, errors.NewBadMatchErr(r.name, s, "anyof:emptyres"))
	}

	if !r.capture {
		return NewResult(r.name, nil, ""), nil
	}

	return all, nil
}

type pickoneparser struct{ s *Sequence }

func (p *pickoneparser) Parse(s string) (*Result, error) { return p.s.pickOne(s) }
func (p *pickoneparser) Name() string                    { return p.s.Name() }

func (r *Sequence) PickOne() Parsable {
	return &pickoneparser{r}
}

func (r *Sequence) pickOne(s string) (*Result, error) {
	res, _ := r.anyOf(s)
	if res.IsEmpy() {
		return NewResult(r.name, nil, ""), nil
	}

	for _, result := range res.ResultMap() {
		return result, nil
	}

	return NewResult(r.name, nil, ""), nil
}

func retErr(n string, err error) (*Result, error) {
	return NewResult(n, nil, ""), err
}

type optionalparser struct {
	s *Sequence
}

func (p *optionalparser) Parse(s string) (*Result, error) {
	return p.s.optional(s)
}

func (p *optionalparser) Name() string {
	return p.s.Name()
}

func (r *Sequence) Optional() Parsable {
	return &optionalparser{r}
}

func (r *Sequence) optional(s string) (*Result, error) {
	res, err := r.Parse(s)
	if err != nil {
		return NewResult(r.name, nil, s), nil
	}
	return res, nil
}
