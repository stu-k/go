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

func isToken(r rune) bool {
	return r == '_' || unicode.IsLetter(r)
}
func parseToken(input string) (Data, string, error) {
	if input == "" {
		panic(fmt.Errorf("token init with: \"\""))
	} else if !isToken(rune(input[0])) {
		panic(fmt.Errorf("token init with: \"%s\"", string(input[0])))
	}

	var sofar string
	for i := 0; i < len(input); i++ {
		r := rune(input[i])
		switch {
		case isToken(r) || (i > 0 && unicode.IsDigit(r)) || r == '_':
			sofar += string(r)
			continue
		case unicode.IsSpace(r):
			return Token{sofar}, input[i+1:], nil
		default:
			if len(sofar) > 0 {
				return Token{sofar}, input[i:], nil
			}
			return handleError(NewUnexpectedTokenErr("token:default", r))
		}
	}
	return Token{sofar}, "", nil
}
