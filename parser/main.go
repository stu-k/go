package main

import (
	"fmt"
	"os"
	"strings"
)

type Data interface {
	Type() string
	Value() any
	String() string
}

type DataUnknown struct{}

func (d DataUnknown) Type() string   { return "unknown" }
func (d DataUnknown) String() string { return "unknown" }
func (d DataUnknown) Value() any     { return nil }

func main() {
	args := os.Args[1:]
	input := strings.Join(args, " ")

	fmt.Printf("parsing: \"%v\"\n", input)
	result, rest, err := parse(input, true)
	if err != nil {
		fmt.Printf("error parsing input: %v\n", err)
		return
	}

	fmt.Printf("result: \"%s\"\n", result)
	fmt.Printf("rest: \"%s\"\n", rest)
}

func parse(input string, first bool) (Data, string, error) {
	if !first {
		input = strings.Trim(input, " ")
	}

	fmt.Printf("\nparse: \"%v\"\n", input)
	if len(input) == 0 {
		return nil, "", nil
	}

	r := rune(input[0])
	switch {
	case r == ' ':
		return parse(input[1:], false)
	case isToken(r):
		fmt.Printf("is token: %s\n", input)
		return parseToken(input)
	case isNum(r):
		fmt.Printf("is num: %s\n", input)
		return parseNum(input)
	case r == '{':
		fmt.Printf("is obj: %s\n", input)
		return parseObj(input)
	case isArr(r):
		fmt.Printf("is arr: %s\n", input)
		return parseArr(input)
	case r == '"':
		fmt.Printf("is str: %s\n", input)
		return parseString(input)
	case r == '(':
		fmt.Printf("is paren: %s\n", input)
		return parseParens(input)
	default:
		return DataUnknown{}, "", NewUnexpectedTokenErr("initial default", r)
	}
}
