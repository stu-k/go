package main

import "fmt"

type Paren struct{ val []Data }

func NewParen(c []Data) Paren {
	return Paren{c}
}
func (p Paren) Type() string { return "paren" }
func (p Paren) Value() any   { return p.val }
func (p Paren) String() string {
	sofar := "paren:("
	for _, d := range p.val {
		sofar += fmt.Sprintf(" %v", d)
	}
	sofar += " )"
	return sofar
}

func isParen(r rune) bool { return r == '(' }
func parseParen(input string) (Data, string, error) {
	if input == "" {
		panic(fmt.Errorf("paren init with: \"\""))
	} else if input[0] != '(' {
		panic(fmt.Errorf("paren init with: \"%s\"", string(input[0])))
	}

	toparse := input[1:]

	parsed := make([]Data, 0)
	for i := 0; i < len(toparse); i++ {
		r := rune(toparse[i])
		switch {
		case r == ' ':
			continue
		case r == ')':
			if len(parsed) == 1 {
				return NewParen(parsed), toparse[i+1:], nil
			}
			return handleError(NewUnexpectedTokenErr("parens:close", ')'))
		default:
			data, rest, err := parse(toparse[i:], false)
			if err != nil {
				return handleError(err)
			}
			if len(parsed) > 0 {
				return handleError(NewUnexpectedTokenErr("parens:default", rune(toparse[i])))
			}
			parsed = append(parsed, data)
			if len(rest) > 0 {
				toparse = rest
				i = -1
				continue
			}
			return handleError(NewExpectationErr(')', ' '))
		}
	}

	return handleError(NewExpectationErr(')', ' '))
}
