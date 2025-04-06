package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
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
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		input := scanner.Text()

		fmt.Printf("parsing: \"%v\"\n", input)
		result, err := mainParse(input)
		if err != nil {
			fmt.Printf("error parsing input: %v\n", err)
			return
		}

		fmt.Printf("result: \"%+v\"\n", result)
	}
}

func mainParse(input string) ([]Data, error) {
	data := make([]Data, 0)

	result, rest, err := parse(input, true)
	if err != nil {
		return nil, err
	}
	data = append(data, result)

	for rest != "" {
		result, rest, err = parse(rest, false)
		if err != nil {
			return nil, err
		}
		data = append(data, result)
	}

	return data, nil
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
	case unicode.IsSpace(r):
		return parse(input[1:], false)
	case isToken(r):
		fmt.Printf("is token: %s\n", input)
		return parseToken(input)
	case isNum(r):
		fmt.Printf("is num: %s\n", input)
		return parseNum(input)
	case isObj(r):
		fmt.Printf("is obj: %s\n", input)
		return parseObj(input)
	case isArr(r):
		fmt.Printf("is arr: %s\n", input)
		return parseArr(input)
	case isStr(r):
		fmt.Printf("is str: %s\n", input)
		return parseString(input)
	case isParen(r):
		fmt.Printf("is paren: %s\n", input)
		return parseParen(input)
	// case isOp(string(r)):
	// 	fmt.Printf("is op: %s\n", input)
	// 	return parseOp(input)
	default:
		return DataUnknown{}, "", NewUnexpectedTokenErr("initial:default", r)
	}
}
