package main

import (
	"bufio"
	"fmt"
	"os"
	"unicode"

	"github.com/stu-k/go/parser/errors"
)

type Data interface {
	Type() string
	Value() any
	String() string
}

type ParseChecker interface {
	Parse(s string) (Data, string, error)
	Check(r rune) bool
}

var mainOpts = []ParseChecker{
	&Var{},
	&Num{},
	&Obj{},
	&Arr{},
	&Str{},
	&Paren{},
	&Op{},
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

		fmt.Printf("result: %+v\n", result)
	}
}

func mainParse(input string) ([]Data, error) {
	data := make([]Data, 0)

	result, rest, err := parse(input, mainOpts)
	if err != nil {
		return nil, err
	}
	data = append(data, result)

	for rest != "" {
		result, rest, err = parse(rest, mainOpts)
		if err != nil {
			return nil, err
		}
		data = append(data, result)
	}

	return data, nil
}

func parse(input string, opts []ParseChecker) (Data, string, error) {
	fmt.Printf("\nparse: \"%v\"\n", input)
	if len(input) == 0 {
		return nil, "", errors.NewEndOfInputErr()
	}

	r := rune(input[0])
	if unicode.IsSpace(r) {
		return parse(input[1:], opts)
	}

	type result struct {
		data Data
		rest string
		err  error
	}

	var okResults []result
	var errResults []result
	for _, opt := range mainOpts {
		ok := opt.Check(r)
		if !ok {
			continue
		}

		res, rest, err := opt.Parse(input)
		if err != nil {
			errResults = append(okResults, result{
				res, rest, err,
			})
			continue
		}
		okResults = append(okResults, result{
			res, rest, err,
		})
	}

	fmt.Printf("ok results: %v\n", okResults)
	fmt.Printf("err results: %v\n", errResults)

	if len(okResults) > 0 {
		r := okResults[0]
		return r.data, r.rest, r.err
	}

	if len(errResults) > 0 {
		r := errResults[0]
		return r.data, r.rest, r.err
	}

	return nil, "", errors.NewUnexpectedCharErr(r)
}
