package main

import (
	"fmt"
	"unicode"
)

type Token struct {
	val string
}

func (t Token) Type() string   { return "token" }
func (t Token) Value() any     { return t.val }
func (t Token) String() string { return fmt.Sprintf("token:%s", t.val) }
func (t Token) Check(r rune) bool {
	return r == '_' || unicode.IsLetter(r)
}

func (t Token) Parse(s string) (Data, string, error) {
	if err := checkInit(t, s); err != nil {
		panic(err)
	}

	var sofar string
	for i := 0; i < len(s); i++ {
		r := rune(s[i])
		switch {
		case t.Check(r) || (i > 0 && unicode.IsDigit(r)) || r == '_':
			sofar += string(r)
			continue
		case unicode.IsSpace(r):
			return Token{sofar}, s[i+1:], nil
		default:
			if len(sofar) > 0 {
				return Token{sofar}, s[i:], nil
			}
			return handleError(NewUnexpectedTokenErr("token:default", r))
		}
	}

	return Token{sofar}, "", nil
}
