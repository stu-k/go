package main

import "fmt"

type Str struct{ val string }

func NewStr(s string) Str    { return Str{s} }
func (s Str) Type() string   { return "str" }
func (s Str) Value() any     { return s.val }
func (s Str) String() string { return fmt.Sprintf("str:\"%s\"", s.val) }

func isStr(r rune) bool { return r == '"' }
func parseString(input string) (Data, string, error) {
	if input == "" {
		panic(fmt.Errorf("str init with \"\""))
	} else if input[0] != '"' {
		panic(fmt.Errorf("str init with \"%s\"", string(input[0])))
	}

	toparse := input[1:]
	var sofar string
	for i := 0; i < len(toparse); i++ {
		r := rune(toparse[i])
		switch {
		case r == '"':
			return NewStr(sofar), toparse[i+1:], nil
		default:
			sofar += string(r)
			continue
		}
	}

	return handleError(NewSingleExpectationErr('"'))
}
