package main

import (
	"fmt"
	"unicode"
)

type Var struct {
	val string
}

func (t *Var) Type() string   { return "var" }
func (t *Var) Value() any     { return t.val }
func (t *Var) String() string { return fmt.Sprintf("var:%s", t.val) }

func (t *Var) Check(r rune) bool {
	return r == '_' || unicode.IsLetter(r)
}
func (t *Var) Parse(s string) (Data, string, error) {
	if err := checkInit(t, s); err != nil {
		return handleError(err)
	}

	var res string
	for i := 0; i < len(s); i++ {
		r := rune(s[i])
		switch {
		case t.Check(r) || (i > 0 && unicode.IsDigit(r)) || r == '_':
			res += string(r)
			continue
		case unicode.IsSpace(r):
			return &Var{res}, s[i+1:], nil
		default:
			if len(res) > 0 {
				return &Var{res}, s[i:], nil
			}
			return handleError(NewUnexpectedCharErr(r))
		}
	}

	return &Var{res}, "", nil
}
