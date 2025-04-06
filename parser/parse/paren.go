package parse

import (
	"fmt"
	"unicode"

	"github.com/stu-k/go/parser/errors"
)

type Paren struct{ val []Data }

func NewParen(c []Data) *Paren {
	return &Paren{c}
}
func (p *Paren) Type() string { return "paren" }
func (p *Paren) Value() any   { return p.val }
func (p *Paren) String() string {
	sofar := "paren:("
	for _, d := range p.val {
		sofar += fmt.Sprintf(" %v", d)
	}
	sofar += " )"
	return sofar
}

func (p *Paren) Check(r rune) bool { return r == '(' }
func (p *Paren) Parse(s string) (Data, string, error) {
	if err := errors.CheckInit(p, s); err != nil {
		return errors.HandeleError(err)
	}

	toparse := s[1:]

	res := make([]Data, 0)
	for i := 0; i < len(toparse); i++ {
		r := rune(toparse[i])
		switch {
		case unicode.IsSpace(r):
			continue
		case r == ')':
			if len(res) == 1 {
				return NewParen(res), toparse[i+1:], nil
			}
			return errors.HandeleError(errors.NewUnexpectedCharErr(')'))
		default:
			data, rest, err := parse(toparse[i:], mainOpts)
			if err != nil {
				return errors.HandeleError(err)
			}
			if len(res) > 0 {
				return errors.HandeleError(errors.NewUnexpectedCharErr(rune(toparse[i])))
			}
			res = append(res, data)
			if len(rest) > 0 {
				toparse = rest
				i = -1
				continue
			}
			return errors.HandeleError(errors.NewExpectedCharErr(')'))
		}
	}

	return errors.HandeleError(errors.NewExpectedCharErr(')'))
}
