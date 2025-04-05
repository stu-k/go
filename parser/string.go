package main

import "fmt"

type String struct{ val string }

func NewString(s string) String {
	return String{s}
}
func (s String) Type() string   { return "string" }
func (s String) Value() any     { return s.val }
func (s String) String() string { return fmt.Sprintf("string:%s", s.val) }

func parseString(input string) (Data, string, error) {
	if input == "" {
		panic(fmt.Errorf("string init with \"\""))
	} else if input[0] != '"' {
		panic(fmt.Errorf("string init with \"%s\"", string(input[0])))
	}

	toParse := input[1:]
	var sofar string
	for i := 0; i < len(toParse); i++ {
		r := rune(toParse[i])
		switch {
		case r == '"':
			return NewString(sofar), toParse[i+1:], nil
		default:
			sofar += string(r)
			continue
		}
	}

	return handleError(NewSingleExpectationErr('"'))
}
