package main

import (
	"fmt"
	"unicode"
)

type Arr struct{ val []Data }

func (a *Arr) Type() string { return "arr" }
func (a *Arr) Value() any   { return a.val }
func (a *Arr) String() string {
	sofar := "arr:["
	for _, d := range a.val {
		sofar += fmt.Sprintf(" %v", d)
	}
	sofar += " ]"
	return sofar
}

func (a *Arr) Check(r rune) bool { return r == '[' }
func (a *Arr) Parse(s string) (Data, string, error) {
	if err := checkInit(a, s); err != nil {
		return handleError(err)
	}

	toparse := s[1:]
	res := make([]Data, 0)

	lastWasComma := false
	for i := 0; i < len(toparse); i++ {
		r := rune(toparse[i])
		switch {
		case unicode.IsSpace(r):
			continue
		case r == ']':
			if lastWasComma {
				return handleError(NewExpectationErr(']', ','))
			}
			return &Arr{res}, toparse[i+1:], nil
		case r == ',':
			if len(res) == 0 {
				return handleError(NewUnexpectedCharErr("arr:comma", r))
			}
			lastWasComma = true
			continue
		default:
			if len(res) > 0 && !lastWasComma {
				return handleError(NewSingleExpectationErr(']'))
			}

			data, rest, err := parse(toparse[i:], mainOpts, false)
			if err != nil {
				return handleError(err)
			}

			res = append(res, data)
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
