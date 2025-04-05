package main

import "fmt"

type Arr struct{ contains []Data }

func (a Arr) Type() string { return "array" }
func (a Arr) Value() any   { return a.contains }
func (a Arr) String() string {
	sofar := "["
	for _, d := range a.contains {
		sofar += fmt.Sprintf(" %v", d)
	}
	sofar += " ]"
	return sofar
}

func isArr(r rune) bool { return r == '[' }
func parseArr(input string) (Data, string, error) {
	parsed := make([]Data, 0)
	lastWasComma := false
	started := false
	toParse := input

	for i := 0; i < len(toParse); i++ {
		r := rune(toParse[i])
		if !started {
			if r == '[' {
				started = true
				lastWasComma = false
				continue
			}
			return nil, "", NewUnexpedTokenErr("arr started", '[')
		}
		switch {
		case r == ' ':
			continue
		case r == ']':
			if lastWasComma {
				return nil, "", NewExpectationErr(']', ',')
			}
			return Arr{parsed}, toParse[i:], nil
		case r == ',':
			if len(parsed) == 0 {
				return nil, "", NewUnexpedTokenErr("arr comma", r)
			}
			lastWasComma = true
			continue
		default:
			data, rest, err := parse(toParse[i:], false)
			if err != nil {
				return nil, "", err
			}
			parsed = append(parsed, data)
			if rest != "" {
				toParse = rest
				i = -1
			}
			lastWasComma = false
			continue
		}
	}
	return nil, "", NewExpectationErr(']', ' ')
}
