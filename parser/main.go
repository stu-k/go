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
			continue
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
	if unicode.IsSpace(r) {
		return parse(input[1:], false)
	}

	switch {
	case (&Var{}).Check(r):
		return (&Var{}).Parse(input)
	case (&Num{}).Check(r):
		return (&Num{}).Parse(input)
	case (&Obj{}).Check(r):
		return (&Obj{}).Parse(input)
	case (&Arr{}).Check(r):
		return (&Arr{}).Parse(input)
	case (&Str{}).Check(r):
		return (&Str{}).Parse(input)
	case (&Paren{}).Check(r):
		return (&Paren{}).Parse(input)
	case (&Op{}).Check(r):
		return (&Op{}).Parse(input)
	default:
		return nil, "", NewUnexpectedCharErr("initial:default", r)
	}
}
