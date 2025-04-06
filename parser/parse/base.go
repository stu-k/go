package parse

import (
	"fmt"
	"unicode"

	"github.com/stu-k/go/parser/errors"
)

var Alpha = &Rule{
	Name:        "alpha",
	Count:       -1,
	Start:       0,
	End:         0,
	Check:       unicode.IsLetter,
	IgnoreSpace: true,
	AtLeastOne:  true,
}

type Rule struct {
	Name        string
	Count       int
	Start, End  rune
	Check       func(rune) bool
	IgnoreSpace bool
	AtLeastOne  bool
}

func (a *Rule) Clone() *Rule {
	return &Rule{
		Name:        a.Name,
		Count:       a.Count,
		Start:       a.Start,
		End:         a.End,
		Check:       a.Check,
		IgnoreSpace: a.IgnoreSpace,
		AtLeastOne:  a.AtLeastOne,
	}
}

func (a *Rule) WithCount(n int) *Rule {
	new := a.Clone()
	new.Count = n
	return new
}

func (a *Rule) Parse(s string) (string, string, error) {
	useStart := a.Start != 0
	useEnd := a.End != 0
	useCount := a.Count >= 0

	if s == "" || (useStart && rune(s[0]) != a.Start) {
		return "", "", errors.NewBadMatchErr(a.Name, s)
	}

	toparse := s
	if a.Start != 0 {
		toparse = s[1:]
	}

	var result string
	var count int
	for i, c := range toparse {
		r := rune(c)
		if unicode.IsSpace(r) && a.IgnoreSpace {
			count++
			continue
		}

		if useCount && !useEnd {
			if count >= a.Count {
				return result, s[i:], nil
			}
		}

		if !useCount && useEnd {
			if r == a.End {
				result += string(r)
				return result, s[i+1:], nil
			}
		}

		if useCount && useEnd {
			if count == a.Count+1 {
				if r == a.End {
					result += string(r)
					return result, s[i+1:], nil
				}
				return "", "", errors.NewBadMatchErr(a.Name, s)
			}
		}

		ok := a.Check(r)
		if !ok {
			if a.AtLeastOne && result == "" {
				return "", "", errors.NewBadMatchErr(a.Name, s)
			}
			return result, s[i:], nil
		}

		result += string(r)
		count++
	}

	if useCount && count < a.Count {
		return "", "", errors.NewBadMatchErr(a.Name, s)
	}

	if a.AtLeastOne && result == "" {
		return "", "", errors.NewBadMatchErr(a.Name, s)
	}

	return result, s[count:], nil
}

type Data interface {
	Type() string
	Value() any
	String() string
}

type ParseChecker interface {
	Parse(s string) (Data, string, error)
	Check(r rune) bool
}

var mainOpts = []ParseChecker{
	&Var{},
	&Num{},
	&Obj{},
	&Arr{},
	&Str{},
	&Paren{},
	&Op{},
}

func Parse(input string) ([]Data, error) {
	data := make([]Data, 0)

	result, rest, err := parse(input, mainOpts)
	if err != nil {
		return nil, err
	}
	data = append(data, result)

	for rest != "" {
		result, rest, err = parse(rest, mainOpts)
		if err != nil {
			return nil, err
		}
		data = append(data, result)
	}

	return data, nil
}

func parse(input string, opts []ParseChecker) (Data, string, error) {
	fmt.Printf("\nparse: \"%v\"\n", input)
	if len(input) == 0 {
		return nil, "", errors.NewEndOfInputErr()
	}

	r := rune(input[0])
	if unicode.IsSpace(r) {
		return parse(input[1:], opts)
	}

	type result struct {
		data Data
		rest string
		err  error
	}

	var okResults []result
	var errResults []result
	for _, opt := range mainOpts {
		ok := opt.Check(r)
		if !ok {
			continue
		}

		res, rest, err := opt.Parse(input)
		if err != nil {
			errResults = append(okResults, result{
				res, rest, err,
			})
			continue
		}
		okResults = append(okResults, result{
			res, rest, err,
		})
	}

	fmt.Printf("ok results: %v\n", okResults)
	fmt.Printf("err results: %v\n", errResults)

	if len(okResults) > 0 {
		r := okResults[0]
		return r.data, r.rest, r.err
	}

	if len(errResults) > 0 {
		r := errResults[0]
		return r.data, r.rest, r.err
	}

	return nil, "", errors.NewUnexpectedCharErr(r)
}
