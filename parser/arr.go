package main

import (
	"fmt"
	"unicode"
)

type Arr struct{ val []Data }

func (a Arr) Type() string { return "arr" }
func (a Arr) Value() any   { return a.val }
func (a Arr) String() string {
	sofar := "arr:["
	for _, d := range a.val {
		sofar += fmt.Sprintf(" %v", d)
	}
	sofar += " ]"
	return sofar
}

func (a Arr) Check(r rune) bool { return r == '[' }
func (a Arr) Parse(s string) (Data, string, error) {
	if err := checkInit(a, s); err != nil {
		panic(err)
	}

	parsed := make([]Data, 0)
	lastWasComma := false
	toparse := s[1:]

	for i := 0; i < len(toparse); i++ {
		r := rune(toparse[i])
		switch {
		case unicode.IsSpace(r):
			continue
		case r == ']':
			if lastWasComma {
				return handleError(NewExpectationErr(']', ','))
			}
			return Arr{parsed}, toparse[i+1:], nil
		case r == ',':
			if len(parsed) == 0 {
				return handleError(NewUnexpectedTokenErr("arr:comma", r))
			}
			lastWasComma = true
			continue
		default:
			if len(parsed) > 0 && !lastWasComma {
				return handleError(NewSingleExpectationErr(']'))
			}

			data, rest, err := parse(toparse[i:], false)
			if err != nil {
				return handleError(err)
			}

			parsed = append(parsed, data)
			if rest != "" {
				toparse = rest
				i = -1
			}
			lastWasComma = false
			continue
		}
	}
	return handleError(NewExpectationErr(']', ' '))
}
