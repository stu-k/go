package main

import "fmt"

type Parens struct{ val []Data }

func NewParens(c []Data) Parens {
	return Parens{c}
}
func (p Parens) Type() string { return "parens" }
func (p Parens) Value() any   { return p.val }
func (p Parens) String() string {
	sofar := "("
	for _, d := range p.val {
		sofar += fmt.Sprintf(" %v", d)
	}
	sofar += " )"
	return sofar
}

func parseParens(input string) (Data, string, error) {
	if input == "" {
		panic(fmt.Errorf("parens init with: \"\""))
	} else if input[0] != '(' {
		panic(fmt.Errorf("parens init with: \"%s\"", string(input[0])))
	}

	toParse := input[1:]

	parsed := make([]Data, 0)
	for i := 0; i < len(toParse); i++ {
		r := rune(toParse[i])
		switch {
		case r == ' ':
			continue
		case r == ')':
			if len(parsed) == 1 {
				return NewParens(parsed), toParse[i+1:], nil
			}
			return handleError(NewUnexpectedTokenErr("parens close paren", ')'))
		default:
			data, rest, err := parse(toParse[i:], false)
			if err != nil {
				return handleError(err)
			}
			if len(parsed) > 0 {
				return handleError(NewUnexpectedTokenErr("parens default", rune(toParse[i])))
			}
			parsed = append(parsed, data)
			if len(rest) > 0 {
				toParse = rest
				i = -1
				continue
			}
			return handleError(NewExpectationErr(')', ' '))
		}
	}

	return handleError(NewExpectationErr(')', ' '))
}
