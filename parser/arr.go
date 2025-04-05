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

func isArr(r rune) bool { return r == '[' }
func parseArr(input string) (Data, string, error) {
	if input == "" {
		panic(fmt.Errorf("arr init with: \"\""))
	} else if input[0] != '[' {
		panic(fmt.Errorf("arr init with: \"%s\"", string(input[0])))
	}

	parsed := make([]Data, 0)
	lastWasComma := false
	toparse := input[1:]

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
