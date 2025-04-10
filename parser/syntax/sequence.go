package syntax

import (
	"fmt"
	"strconv"
	"strings"

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
		results, err := r.Parse(all.Rest())
		if err != nil {
			if all.Len() == 0 {
				return retErr(r.name, err)
			}
			return all, nil
		}

		if r.capture {
			all.Append(results)
		}

		all.SetRest(results.Rest())
		if len(results.Rest()) == 0 {
			break
		}
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
	all := NewResult(r.name, nil, s)
	for _, p := range r.list {
		results, err := p.Parse(s)
		if err != nil {
			continue
		}

		if r.capture {
			all.Append(results)
		}
	}

	if all.Len() == 0 {
		return retErr(r.name, errors.NewBadMatchErr(r.name, s, "anyof:emptyres"))
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
	res, err := r.anyOf(s)
	if err != nil {
		return retErr(r.name, err)
	}
	if res.IsEmpy() {
		return retErr(r.name, errors.NewBadMatchErr(r.name, s, "pickone:isempty"))
	}

	for _, result := range res.ResultMap() {
		return result, nil
	}

	return retErr(r.name, errors.NewBadMatchErr(r.name, s, "pickone:nores"))
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

type seqStrArgs struct {
	rule        *Rule
	char        rune
	repeat      int
	cap, usecap bool
}

func NewSequenceFromStrs(name string, parts ...string) (*Sequence, error) {
	return newSequenceFromStrs(name, defaultRulemap, parts...)
}

func newSequenceFromStrs(name string, pmap map[string]*Rule, parts ...string) (*Sequence, error) {
	errFn := func(arg string, i, j int) error {
		return fmt.Errorf(
			"error creating sequence from strs: invalid arg \"%v\" in segment %d, arg %d",
			arg, i, j,
		)
	}

	if len(parts) == 0 {
		return nil, fmt.Errorf("error creating sequence from strs: invalid string \"%s\"", parts)
	}

	sq := NewSequence(name)
	for i, part := range parts {

		var sqa seqStrArgs

		args := strings.Split(part, ", ")
		for j, arg := range args {
			if len(arg) == 0 || len(arg) == 1 {
				return nil, errFn(arg, i, j)
			}

			r := rune(arg[1])
			switch rune(arg[0]) {

			case 'r':
				// dangerous to pull from pmap
				// can mutate origin rules
				seq, ok := pmap[arg[1:]]
				if !ok || seq == nil {
					return nil, errFn(arg, i, j)
				}
				sqa.rule = seq.Named(arg[1:])
				continue

			case '#':
				ct, err := strconv.Atoi(arg[1:])
				if err != nil {
					return nil, errFn(arg, i, j)
				}
				sqa.repeat = ct
				continue

			case 'c':
				sqa.char = r
				continue

			case 'g':
				if r == '0' {
					sqa.cap = false
					sqa.usecap = true
					continue
				}
				if r == '1' {
					sqa.cap = true
					sqa.usecap = true
					continue
				}
				return nil, errFn(arg, i, j)

			default:
				return nil, errFn(arg, i, j)

			}
		}

		rule := NewRule("any")
		if sqa.rule != nil {
			rule = sqa.rule
		} else {
			rule = rule.Named(part)
		}
		if sqa.repeat > 0 {
			rule = rule.Repeat(sqa.repeat)
		}
		if sqa.char != 0 {
			rule = rule.CheckChar(func(r rune) bool { return r == sqa.char })
		}
		if sqa.usecap {
			rule = rule.Capture(sqa.cap)
		}
		if rule.IsAny() {
			return nil, fmt.Errorf("error creating sequence from strs: can't add empty rule \"%s\"", part)
		}
		sq.Add(rule)
	}

	if sq.Len() == 0 {
		return nil, fmt.Errorf("error creating sequence from strs: can't use empty sequence")
	}

	return sq, nil
}
