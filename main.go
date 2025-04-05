package main

import (
	"fmt"
	"os"
	"strings"
	"unicode"
)

type ParseError struct{ msg string }

func (p ParseError) Error() string {
	return p.msg
}
func NewUnexpedTokenErr(_ string, token rune) ParseError {
	return ParseError{fmt.Sprintf("unexpected token: %s", string(token))}
}

func NewExpectationErr(got rune, want rune) ParseError {
	return ParseError{fmt.Sprintf("expected \"%s\"; got \"%s\"", string(got), string(want))}
}

type Data interface {
	Type() string
	Value() any
	String() string
}

func main() {
	args := os.Args[1:]
	input := strings.Join(args, " ")
	fmt.Printf("parsing %v\n", input)
	result, _, err := initialParse(input)
	if err != nil {
		fmt.Printf("error parsing input: %v\n", err)
		return
	}
	fmt.Println(result)
}

func initialParse(input string) (Data, string, error) {
	fmt.Printf("\ninitialParse: %v\n", input)
	if len(input) == 0 {
		return nil, "", nil
	}

	r := rune(input[0])
	switch {
	case r == ' ':
		return initialParse(input[1:])
	case isToken(r):
		return parseToken(input)
	case unicode.IsDigit(r):
		fmt.Println("is digit")
		return parseNum(input)
	case r == '{':
		fmt.Println("is obj")
		return parseObj(input)
	case isArr(r):
		fmt.Println("is arr")
		return parseArr(input)
	case r == '"':
		fmt.Println("is string")
		return parseString(input)
	case r == '(':
		fmt.Println("is parens")
		return parseParens(input)
	default:
		return nil, "", NewUnexpedTokenErr("initial default", r)
	}
}

func parseObj(input string) (Data, string, error) {
	return nil, "", nil
}

func parseNum(input string) (Data, string, error) {
	return nil, "", nil
}

func parseParens(input string) (Data, string, error) {
	return nil, "", nil
}

func parseString(input string) (Data, string, error) {
	return nil, "", nil
}

type Array struct{ contains []Data }

func (a Array) Type() string { return "array" }
func (a Array) Value() any   { return a.contains }
func (a Array) String() string {
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
	last := input[0]
	started := false
	toParse := input

	for i := 0; i < len(toParse); i++ {
		r := rune(toParse[i])
		if !started {
			if r == '[' {
				started = true
				continue
			}
			return nil, "", NewUnexpedTokenErr("arr started", '[')
		}
		switch {
		case r == ' ':
			continue
		case r == ']':
			if last == ',' {
				return nil, "", NewExpectationErr(']', ',')
			}
			return Array{parsed}, toParse[i:], nil
		case r == ',':
			if len(parsed) == 0 {
				return nil, "", NewUnexpedTokenErr("arr \",\"", r)
			}
			last = ','
			continue
		default:
			if len(parsed) > 0 {
				return nil, "", fmt.Errorf("expected ,")
			}
			data, rest, err := initialParse(toParse[i:])
			if err != nil {
				return nil, "", err
			}
			parsed = append(parsed, data)
			toParse = rest
			i = -1
			continue
		}
	}
	return nil, "", fmt.Errorf("expected ]")
}

type Token struct {
	val string
}

func (t Token) Type() string   { return "token" }
func (t Token) Value() any     { return t.val }
func (t Token) String() string { return fmt.Sprintf("token: %s", t.val) }

func isToken(r rune) bool {
	return r == '_' || unicode.IsLetter(r)
}
func parseToken(input string) (Data, string, error) {
	var sofar string
	for i := 0; i < len(input); i++ {
		r := rune(input[i])
		switch {
		case isToken(r) || (i > 0 && unicode.IsDigit(r)):
			sofar += string(r)
			continue
		case r == ' ':
			if len(sofar) == 0 {
				panic(fmt.Sprintf("token was empty: \"%v\"", input))
			}
			return Token{sofar}, input[i+1:], nil
		default:
			if len(sofar) > 0 {
				return Token{sofar}, input[i:], nil
			}
			return nil, "", fmt.Errorf("unexpected token: %v", string(r))
		}
	}
	return Token{sofar}, "", nil
}
