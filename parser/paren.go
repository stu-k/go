package main

import (
	"fmt"
	"unicode"
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
	if err := checkInit(p, s); err != nil {
		return handleError(err)
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
			return handleError(NewUnexpectedCharErr("parens:close", ')'))
		default:
			data, rest, err := parse(toparse[i:], mainOpts, false)
			if err != nil {
				return handleError(err)
			}
			if len(res) > 0 {
				return handleError(NewUnexpectedCharErr("parens:default", rune(toparse[i])))
			}
			res = append(res, data)
			if len(rest) > 0 {
				toparse = rest
				i = -1
				continue
			}
			return handleError(NewExpectationErr(')', ' '))
		}
	}

	return handleError(NewExpectationErr(')', ' '))
}
