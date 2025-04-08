package syntax

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/stu-k/go/parser/errors"
)

type Parsable interface {
	Parse(string) (*ParseResult, error)
	Name() string
}

type ParseFn func(string) (*ParseResult, error)

type Sequence struct {
	name    string
	list    []Parsable
	capture bool
}

func NewSequence(name string, rules ...Parsable) *Sequence {
	return &Sequence{name, rules, true}
}

func (r *Sequence) Len() int          { return len(r.list) }
func (r *Sequence) Add(v ...Parsable) { r.list = append(r.list, v...) }
func (r *Sequence) Name() string      { return r.name }
func (r *Sequence) SetCapture(v bool) { r.capture = v }
func (r *Sequence) Parse(s string) (*ParseResult, error) {
	results := NewParseResult(r.name, nil, s)
	for _, rule := range r.list {
		result, err := rule.Parse(results.Rest())
		if err != nil {
			return returnPr(r.name, s, err)
		}

		if r.capture {
			results.Append(result)
		}
		results.SetRest(result.Rest())
	}
	return results, nil
}

func (r *Sequence) UntilFail(s string) (*ParseResult, error) {
	all := NewParseResult(r.name, nil, s)
	for {
		results, err := r.Parse(all.Rest())
		if err != nil {
			if all.Len() == 0 {
				return returnPr(r.name, s, err)
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

func (r *Sequence) AnyOf(s string) (*ParseResult, error) {
	all := NewParseResult(r.name, nil, s)
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
		return returnPr(r.name, s, errors.NewBadMatchErr(r.name, s))
	}

	return all, nil
}

func returnPr(n, s string, err error) (*ParseResult, error) {
	return NewParseResult(n, nil, s), err
}

type seqStrArgs struct {
	rule        *Rule
	char        rune
	count       int
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
				seq, ok := pmap[arg[1:]]
				if !ok || seq == nil {
					return nil, errFn(arg, i, j)
				}
				sqa.rule = seq
				continue

			case '#':
				ct, err := strconv.Atoi(arg[1:])
				if err != nil {
					return nil, errFn(arg, i, j)
				}
				sqa.count = ct
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

		rule := NewRule()
		if sqa.rule != nil {
			rule = sqa.rule
		} else {
			rule = rule.Named(part)
		}
		if sqa.count != 0 {
			rule = rule.Count(sqa.count)
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
