package main

import (
	"fmt"
	"unicode"
)

type Paren struct{ val []Data }

func NewParen(c []Data) Paren {
	return Paren{c}
}
func (p Paren) Type() string { return "paren" }
func (p Paren) Value() any   { return p.val }
func (p Paren) String() string {
	sofar := "paren:("
	for _, d := range p.val {
		sofar += fmt.Sprintf(" %v", d)
	}
	sofar += " )"
	return sofar
}

func (p Paren) Check(r rune) bool { return r == '(' }
func (p Paren) Parse(s string) (Data, string, error) {
	if err := checkInit(p, s); err != nil {
		panic(err)
	}

	toparse := s[1:]

	parsed := make([]Data, 0)
	for i := 0; i < len(toparse); i++ {
		r := rune(toparse[i])
		switch {
		case unicode.IsSpace(r):
			continue
		case r == ')':
			if len(parsed) == 1 {
				return NewParen(parsed), toparse[i+1:], nil
			}
			return handleError(NewUnexpectedTokenErr("parens:close", ')'))
		default:
			data, rest, err := parse(toparse[i:], false)
			if err != nil {
				return handleError(err)
			}
			if len(parsed) > 0 {
				return handleError(NewUnexpectedTokenErr("parens:default", rune(toparse[i])))
			}
			parsed = append(parsed, data)
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
